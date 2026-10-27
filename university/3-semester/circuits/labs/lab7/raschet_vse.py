"""ЛР-7: расчёт усилителя на 2N2222A по SPICE-модели из базы Multisim.

Модель (MSCOMP_S.PRD, Motorola / Symmetry Design Systems) — Гуммеля–Пуна, n-p-n,
с ёмкостями переходов и диффузионными ёмкостями (TF, TR), как в SPICE3.
Свой решатель методом узловых потенциалов (MNA) на numpy:
  1) рабочая точка — Ньютон по постоянному току (C — разрыв, L — перемычка);
  2) АЧХ — линеаризация в рабочей точке (G = dI/dU, C = dQ/dU) и (G + jωC)·U = I.
Multisim делает то же самое (AC-анализ / плоттер Боде = малосигнальный режим).

Запуск: python raschet_vse.py  → печатает таблицы и пишет zamery.py
"""
import math
from pathlib import Path

import numpy as np

LAB = Path(__file__).resolve().parent

# .MODEL 2N2222A npn
# Q2N2222A__BJT_NPN__1 — та, что стоит на схеме пользователя (Q1 → Edit model, 2026-10-26)
M = dict(IS=3.0611e-14, BF=220, NF=1.00124, VAF=104, IKF=0.52,
         ISE=7.5e-15, NE=1.41, BR=4, NR=1.005, VAR=28,
         IKR=0.24, ISC=10.6525e-12, NC=1.3728, RB=0.13, RE=0.22,
         RC=0.12, CJE=27.01e-12, VJE=0.75, MJE=0.33, TF=0.325e-9,
         XTF=0, VTF=math.inf, ITF=0, CJC=9.12e-12, VJC=0.4089,
         MJC=0.3508, XCJC=1, FC=0.5, TR=100e-9)
VT = 1.380649e-23 * 300.15 / 1.602176634e-19   # 27 °C, как в Multisim
GMIN = 1e-12


def ex(v, n):
    return math.exp(min(v / (n * VT), 80.0)) - 1.0


def qdep(v, cj, vj, mj):
    """заряд обеднённого слоя перехода (SPICE3)"""
    fc = M["FC"]
    if v < fc * vj:
        return vj * cj * (1 - (1 - v / vj) ** (1 - mj)) / (1 - mj)
    f1 = vj * (1 - (1 - fc) ** (1 - mj)) / (1 - mj)
    f2 = (1 - fc) ** (1 + mj)
    f3 = 1 - fc * (1 + mj)
    return cj * (f1 + (f3 * (v - fc * vj) + mj / (2 * vj) * (v * v - (fc * vj) ** 2)) / f2)


def bjt(vbe, vbc):
    """внутренний транзистор: (iб, iк, qбэ, qбк) по SPICE3 bjtload"""
    cbe = M["IS"] * ex(vbe, M["NF"])
    cbc = M["IS"] * ex(vbc, M["NR"])
    cben = M["ISE"] * ex(vbe, M["NE"])
    cbcn = M["ISC"] * ex(vbc, M["NC"])
    q1 = 1.0 / (1.0 - vbc / M["VAF"] - vbe / M["VAR"])
    q2 = cbe / M["IKF"] + cbc / M["IKR"]
    qb = q1 * (1.0 + math.sqrt(max(0.0, 1.0 + 4.0 * q2))) / 2.0
    ib = cbe / M["BF"] + cben + cbc / M["BR"] + cbcn + GMIN * (vbe + vbc)
    ic = (cbe - cbc) / qb - cbc / M["BR"] - cbcn - GMIN * vbc
    argtf = 0.0
    if M["XTF"]:
        argtf = M["XTF"] * (cbe / (cbe + M["ITF"])) ** 2 * math.exp(vbc / (1.44 * M["VTF"]))
    qbe = M["TF"] * (1 + argtf) * cbe / qb + qdep(vbe, M["CJE"], M["VJE"], M["MJE"])
    qbc = M["TR"] * cbc + qdep(vbc, M["CJC"] * M["XCJC"], M["VJC"], M["MJC"])
    return ib, ic, qbe, qbc


class Shema:
    """элементы: ("R"|"C"|"L", a, b, номинал), ("V", a, b, U пост., U перем.), ("Q", к, б, э)"""

    def __init__(self, elems):
        self.el = []
        uzly = []
        for e in elems:
            if e[0] == "Q":
                k, b, em = e[1:]
                vn = [f"{k}'", f"{b}'", f"{em}'"]
                self.el += [("R", k, vn[0], M["RC"]), ("R", b, vn[1], M["RB"]),
                            ("R", em, vn[2], M["RE"]), ("Q", *vn, b)]
            else:
                self.el.append(e)
        for e in self.el:
            for u in e[1:4] if e[0] == "Q" else e[1:3]:
                if u != "0" and u not in uzly:
                    uzly.append(u)
        self.uzly = uzly
        self.idx = {u: i for i, u in enumerate(uzly)}
        self.vetvi = [i for i, e in enumerate(self.el) if e[0] in ("V", "L")]
        self.n = len(uzly) + len(self.vetvi)

    def u(self, x, uz):
        return 0.0 if uz == "0" else x[self.idx[uz]]

    def tok_q(self, x):
        """нелинейная часть: токи и заряды, втекающие в узлы из транзисторов"""
        i = np.zeros(self.n)
        q = np.zeros(self.n)
        for e in self.el:
            if e[0] != "Q":
                continue
            k, b, em, bx = e[1:]
            vb, vk, ve = self.u(x, b), self.u(x, k), self.u(x, em)
            ib, ic, qbe, qbc = bjt(vb - ve, vb - vk)
            qbx = qdep(self.u(x, bx) - vk, M["CJC"] * (1 - M["XCJC"]), M["VJC"], M["MJC"])
            for uz, val in ((b, -ib), (k, -ic), (em, ib + ic)):
                i[self.idx[uz]] += val
            for uz, val in ((b, -qbe - qbc), (k, qbc), (em, qbe)):
                q[self.idx[uz]] += val
            if bx != "0":
                q[self.idx[bx]] -= qbx
            q[self.idx[k]] += qbx
        return i, q

    def lin(self, w):
        """линейная часть: матрица Y (при w = 0 — постоянный ток) и вектор источников"""
        Y = np.zeros((self.n, self.n), complex)
        src_dc = np.zeros(self.n)
        src_ac = np.zeros(self.n, complex)
        nu = len(self.uzly)
        for k in range(nu):             # GMIN на землю, как в SPICE: узел за конденсатором не висит
            Y[k, k] += GMIN

        def stamp(a, b, y):
            ia = self.idx.get(a)
            ib = self.idx.get(b)
            if ia is not None:
                Y[ia, ia] += y
            if ib is not None:
                Y[ib, ib] += y
            if ia is not None and ib is not None:
                Y[ia, ib] -= y
                Y[ib, ia] -= y

        for j, e in enumerate(self.el):
            kind, a, b = e[0], e[1], e[2]
            if kind == "R":
                stamp(a, b, 1 / e[3])
            elif kind == "C":
                stamp(a, b, 1j * w * e[3])
            elif kind in ("V", "L"):
                r = nu + self.vetvi.index(j)
                for uz, s in ((a, 1), (b, -1)):
                    if uz in self.idx:
                        Y[self.idx[uz], r] += s
                        Y[r, self.idx[uz]] += s
                if kind == "L":
                    Y[r, r] -= 1j * w * e[3]
                else:
                    src_dc[r] = e[3]
                    src_ac[r] = e[4]
        return Y, src_dc, src_ac

    def rabochaya_tochka(self):
        Y0, sdc, _ = self.lin(0.0)
        G = Y0.real
        x = np.zeros(self.n)
        for shag in np.linspace(0.1, 1.0, 10):          # плавно поднимаем питание
            for _ in range(300):
                inl, _ = self.tok_q(x)
                F = G @ x - inl - shag * sdc
                J = np.empty((self.n, self.n))
                for k in range(self.n):
                    d = np.zeros(self.n)
                    d[k] = 1e-7
                    J[:, k] = G[:, k] - (self.tok_q(x + d)[0] - inl) / 1e-7
                dx = np.linalg.solve(J, -F)
                m = np.max(np.abs(dx[:len(self.uzly)]))
                if m > 0.05:
                    dx *= 0.05 / m
                x = x + dx
                if m < 1e-10:
                    break
        self.x0 = x
        # малосигнальные G и C в рабочей точке
        i0, q0 = self.tok_q(x)
        self.Gn = np.empty((self.n, self.n))
        self.Cn = np.empty((self.n, self.n))
        for k in range(self.n):
            d = np.zeros(self.n)
            d[k] = 1e-6
            ip, qp = self.tok_q(x + d)
            im, qm = self.tok_q(x - d)
            self.Gn[:, k] = (ip - im) / 2e-6
            self.Cn[:, k] = (qp - qm) / 2e-6
        return x

    def perehod(self, f, amp, periodov=40, na_period=200):
        """форма сигнала во времени (как осциллограф): неявный Эйлер от рабочей точки.
        Возвращает (t, u_вых) за последние 4 периода."""
        G = self.lin(0.0)[0].real
        Cl = self.lin(1.0)[0].imag            # матрица линейных ёмкостей
        for j in self.vetvi:                   # у ветвей L/V мнимая часть — не ёмкость
            Cl[len(self.uzly) + self.vetvi.index(j), :] = 0
            Cl[:, len(self.uzly) + self.vetvi.index(j)] = 0
        _, sdc, sac = self.lin(0.0)
        sac = sac.real
        dt = 1 / f / na_period
        x = self.x0.copy()
        qp = Cl @ x - self.tok_q(x)[1]
        ts, us = [], []
        n_all = periodov * na_period
        for s in range(1, n_all + 1):
            t = s * dt
            src = sdc + sac * amp * math.sin(2 * math.pi * f * t)
            for _ in range(50):
                inl, qnl = self.tok_q(x)
                F = G @ x - inl - src + (Cl @ x - qnl - qp) / dt
                J = np.empty((self.n, self.n))
                for k in range(self.n):
                    d = np.zeros(self.n)
                    d[k] = 1e-7
                    i2, q2 = self.tok_q(x + d)
                    J[:, k] = G[:, k] + Cl[:, k] / dt - (i2 - inl) / 1e-7 - (q2 - qnl) / 1e-7 / dt
                dx = np.linalg.solve(J, -F)
                x = x + dx
                if np.max(np.abs(dx)) < 1e-9:
                    break
            qp = Cl @ x - self.tok_q(x)[1]
            if s > n_all - 4 * na_period:
                ts.append(t)
                us.append(self.u(x, "out"))
        return np.array(ts), np.array(us)

    def napr(self, uz):
        return self.u(self.x0, uz)

    def ac(self, f, uz):
        """комплексное напряжение узла при единичном переменном источнике"""
        w = 2 * math.pi * f
        Y, _, sac = self.lin(w)
        A = Y - self.Gn - 1j * w * self.Cn
        v = np.linalg.solve(A, sac)
        return 0.0 if uz == "0" else v[self.idx[uz]]


# ---------------------------------------------------------------------------
# схемы методички
# ---------------------------------------------------------------------------
def shema_ris3(C1=1e-6, C3=1e-6, C4=100e-9, R4=2e3):
    """рис. 3: резистивный каскад ОЭ. Узлы: vcc, gen, b, c, e, out"""
    el = [("V", "vcc", "0", 10.0, 0.0),      # V1
          ("V", "gen", "0", 0.0, 1.0),       # XFG1 (идеальный)
          ("C", "vcc", "0", 10e-9),          # C2
          ("R", "vcc", "b", 20e3),           # R1
          ("R", "b", "0", 20e3),             # R2
          ("C", "gen", "b", C1),             # C1
          ("R", "vcc", "c", 2.7e3),          # R3
          ("Q", "c", "b", "e"),              # Q1
          ("R", "e", "0", R4),               # R4
          ("C", "c", "out", C3),             # C3
          ("C", "out", "0", 10e-9)]          # C5
    if C4:
        el.append(("C", "e", "0", C4))       # C4
    return Shema(el)


def shema_ris4():
    """рис. 4: резонансный каскад, в коллекторе контур L2 ‖ R3 ‖ C6"""
    el = [("V", "vcc", "0", 10.0, 0.0),
          ("V", "gen", "0", 0.0, 1.0),
          ("C", "vcc", "0", 10e-9),          # C2
          ("R", "vcc", "b", 20e3),           # R1
          ("R", "b", "0", 20e3),             # R2
          ("C", "gen", "b", 10e-9),          # C1 = 10 нФ
          ("L", "vcc", "c", 700e-6),         # L2
          ("R", "vcc", "c", 1e3),            # R3
          ("C", "vcc", "c", 100e-9),         # C6
          ("Q", "c", "b", "e"),
          ("R", "e", "0", 1e3),              # R4
          ("C", "e", "0", 100e-9),           # C4
          ("C", "c", "out", 100e-9),         # C3
          ("C", "out", "0", 10e-9)]          # C5
    return Shema(el)


def K(sh, fr):
    v = sh.ac(fr, "out")
    return abs(v), math.degrees(math.atan2(v.imag, v.real))


def polosa(sh, f_lo=10, f_hi=1e7):
    """Kmax, f на максимуме и граничные частоты по уровню Kmax / √2"""
    fs = np.logspace(math.log10(f_lo), math.log10(f_hi), 400)
    ks = np.array([K(sh, x)[0] for x in fs])
    i = int(np.argmax(ks))
    a, b = fs[max(i - 1, 0)], fs[min(i + 1, len(fs) - 1)]
    for _ in range(80):                 # уточнить максимум
        m1, m2 = a * (b / a) ** 0.382, a * (b / a) ** 0.618
        if K(sh, m1)[0] < K(sh, m2)[0]:
            a = m1
        else:
            b = m2
    fmax = math.sqrt(a * b)
    kmax = K(sh, fmax)[0]
    cel = kmax / math.sqrt(2)

    def bis(a, b):
        ka = K(sh, a)[0] - cel
        for _ in range(60):
            m = math.sqrt(a * b)
            km = K(sh, m)[0] - cel
            if (km > 0) == (ka > 0):
                a, ka = m, km
            else:
                b = m
        return math.sqrt(a * b)

    fn = bis(fs[0], fmax) if ks[0] < cel else None
    fv = bis(fmax, fs[-1]) if ks[-1] < cel else None
    return kmax, fmax, fn, fv


def tochka(sh, R4):
    sh.rabochaya_tochka()
    ub, uc, ue = sh.napr("b"), sh.napr("c"), sh.napr("e")
    return dict(Ub=ub, Uc=uc, Ue=ue, Ube=ub - ue, Uce=uc - ue, Ie=ue / R4)


def thd(sh, fr, amp):
    """размах выхода и коэффициент гармоник по форме сигнала (переходный расчёт)"""
    _, u = sh.perehod(fr, amp, periodov=40, na_period=100)
    u = u - u.mean()
    Y = np.abs(np.fft.rfft(u))
    k = math.sqrt(sum(Y[4 * h] ** 2 for h in range(2, 8))) / Y[4]
    return float(u.max()), float(-u.min()), float(k)


# --- точки ------------------------------------------------------------------
UVH = 0.010                                                     # амплитуда генератора, В
F1 = [100, 200, 300, 500, 700, 1e3, 2e3, 3e3, 5e3, 10e3, 20e3, 30e3, 50e3, 70e3, 100e3]
C1_VAR = [5e-9, 1e-6]
C3_VAR = [1e-9, 0.8e-6]
C4_VAR = [0, 10e-9, 100e-9, 1e-6, 10e-6]
F_C4 = [100, 300, 1e3, 3e3, 10e3, 30e3, 100e3]
R4_VAR = [500, 1e3, 2e3, 3e3, 5e3, 10e3]
F4 = [5e3, 10e3, 14e3, 16e3, 17e3, 17.5e3, 18e3, 18.25e3, 18.5e3, 19e3, 19.5e3, 20e3,
      22e3, 25e3, 30e3, 50e3]


def vse():
    Z = {}
    sh = shema_ris3()
    Z["TOCHKA"] = tochka(sh, 2e3)
    Z["P1"] = [K(sh, x) for x in F1]
    Z["P1_POLOSA"] = polosa(sh)
    Z["P1_FORMA"] = thd(sh, 10e3, UVH)
    Z["P2_C1"] = {}
    Z["P2_C1_POLOSA"] = {}
    for c in C1_VAR:
        s = shema_ris3(C1=c)
        s.rabochaya_tochka()
        Z["P2_C1"][c] = [K(s, x)[0] for x in F1]
        Z["P2_C1_POLOSA"][c] = polosa(s)
    Z["P2_C3"] = {}
    Z["P2_C3_POLOSA"] = {}
    for c in C3_VAR:
        s = shema_ris3(C3=c)
        s.rabochaya_tochka()
        Z["P2_C3"][c] = [K(s, x)[0] for x in F1]
        Z["P2_C3_POLOSA"][c] = polosa(s)
    Z["P3_C4"] = {}
    Z["P3_C4_FORMA"] = {}
    for c in C4_VAR:
        s = shema_ris3(C4=c)
        s.rabochaya_tochka()
        Z["P3_C4"][c] = [K(s, x)[0] for x in F_C4]
        Z["P3_C4_FORMA"][c] = thd(s, 10e3, UVH)
    Z["P3_R4"] = {}
    for r in R4_VAR:
        s = shema_ris3(R4=r)
        tc = tochka(s, r)
        tc["K1k"] = K(s, 1e3)[0]
        tc["K10k"] = K(s, 10e3)[0]
        Z["P3_R4"][r] = tc
    s4 = shema_ris4()
    s4.rabochaya_tochka()
    Z["TOCHKA4"] = {"Ub": s4.napr("b"), "Ue": s4.napr("e"), "Uc": s4.napr("c"),
                    "Ie": s4.napr("e") / 1e3}
    Z["P4"] = [K(s4, x) for x in F4]
    Z["P4_POLOSA"] = polosa(s4, 1e3, 1e6)
    return Z


NL = chr(10)


def main():
    Z = vse()
    for k, v in Z.items():
        print(k, v)
    with open(LAB / "zamery.py", "w", encoding="utf-8") as fh:
        fh.write('"""ЛР-7: результаты расчёта (В, А, Гц). Пишет raschet_vse.py — руками не править."""' + NL)
        fh.write("import numpy as np" + NL)
        for name, val in [("UVH", UVH), ("F1", F1), ("C1_VAR", C1_VAR), ("C3_VAR", C3_VAR),
                          ("C4_VAR", C4_VAR), ("F_C4", F_C4), ("R4_VAR", R4_VAR), ("F4", F4)]:
            fh.write(f"{name} = {val!r}" + NL)
        for k, v in Z.items():
            fh.write(f"{k} = {v!r}" + NL)


if __name__ == "__main__":
    main()
