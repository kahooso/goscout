"""ЛР-8: расчёт RC-автогенератора на ОУ (мост Вина) и схем интегратора / дифференциатора.

Модель:
  - диоды D1, D2 — IDEALD из базы Multisim (DIODE_VIRTUAL): IS = 1e-14, N = 1, RS = 0;
    встречно-параллельно: i = 2·IS·sh(u / Vt);
  - ОУ в автогенераторе и режимах 2, 3 — идеальный (V+ = V−): частоты 1…5 кГц много меньше
    полосы LM741 (≈ 1 МГц / K ≈ 300 кГц), ошибка частоты < 1 %;
  - ОУ в интеграторе / дифференциаторе — LM741 «типовой»: A0 = 2·10⁵, GBW = 1 МГц,
    скорость нарастания SR = 0,5 В/мкс, насыщение ±14 В (питание ±15 В).
  - R4 (потенциометр 40 кОм): R_ос = доля·40 кОм (какой конец подключён — сверить по скрину).

Запуск: python raschet_vse.py  → печатает таблицы и пишет zamery.py
"""
import math
from pathlib import Path

import numpy as np

LAB = Path(__file__).resolve().parent

VT = 1.380649e-23 * 300.15 / 1.602176634e-19   # 27 °C, как в Multisim
IS = 1e-14
R3, R4 = 10e3, 40e3
R, C = 5e3, 15e-9                               # R1 = R2, C1 = C2
A0, GBW, SR, VSAT = 2e5, 1e6, 0.5e6, 14.0


# ---------------------------------------------------------------------------
# обратная связь: ОУ идеальный, V− = V+, ток через R3 идёт через R_ос ‖ D1, D2
# ---------------------------------------------------------------------------
def u_os(i, roc):
    """напряжение на цепи ООС (R_ос ‖ диоды) при токе i; roc = None — R4 отключён"""
    if roc is None:
        return VT * math.asinh(i / (2 * IS))
    lo, hi = -50.0, 50.0                    # монотонно — делим пополам
    for _ in range(80):
        m = (lo + hi) / 2
        if m / roc + 2 * IS * math.sinh(max(min(m / VT, 700), -700)) < i:
            lo = m
        else:
            hi = m
    return (lo + hi) / 2


def tabl_vyhod(roc, vmax=3.0, n=1201):
    """таблица U+ → Uвых = U+ + u_ос(U+/R3) для быстрой интерполяции"""
    vp = np.linspace(-vmax, vmax, n)
    vo = np.array([v + u_os(v / R3, roc) for v in vp])
    return vp, np.clip(vo, -VSAT, VSAT)


# ---------------------------------------------------------------------------
# мост Вина: β(f) = 1 / (3 + j(ωRC − 1/(ωRC)))
# ---------------------------------------------------------------------------
def beta(f, r=R, c=C):
    x = 2 * math.pi * f * r * c
    b = 1 / complex(3, x - 1 / x)
    return abs(b), math.degrees(math.atan2(b.imag, b.real))


def f0(r=R, c=C):
    return 1 / (2 * math.pi * r * c)


# ---------------------------------------------------------------------------
# режим 1: автоколебания (во времени, RK4)
#   C2 — последовательно с R2 (выход → «+»), C1 ‖ R1 — «+» на землю
# ---------------------------------------------------------------------------
def generator(dolya, r=R, c=C, t_end=0.06, dt=0.5e-6):
    roc = dolya * R4
    vp_t, vo_t = tabl_vyhod(roc)

    def deriv(s):
        vc2, vp = s
        vo = np.interp(vp, vp_t, vo_t)
        i = (vo - vc2 - vp) / r
        return np.array([i / c, (i - vp / r) / c])

    s = np.array([0.0, 1e-3])               # толчок 1 мВ — «шум при включении»
    nst = int(t_end / dt)
    tt = np.empty(nst)
    vo_s = np.empty(nst)
    vp_s = np.empty(nst)
    for k in range(nst):
        k1 = deriv(s)
        k2 = deriv(s + dt / 2 * k1)
        k3 = deriv(s + dt / 2 * k2)
        k4 = deriv(s + dt * k3)
        s = s + dt / 6 * (k1 + 2 * k2 + 2 * k3 + k4)
        tt[k] = (k + 1) * dt
        vp_s[k] = s[1]
        vo_s[k] = np.interp(s[1], vp_t, vo_t)
    return tt, vo_s, vp_s


def analiz(tt, vo, vp, t_okno=0.01):
    """последние t_okno секунд: частота, амплитуда, действующие, Kг"""
    m = tt > tt[-1] - t_okno
    t, y, z = tt[m], vo[m], vp[m]
    up = float(y.max())
    if up < 0.05:
        return dict(gen=False, Um=up, U2=float(np.sqrt(np.mean(y ** 2))), U3=0.0, f=0.0, kg=0.0)
    zc = [t[i] - y[i] * (t[i + 1] - t[i]) / (y[i + 1] - y[i])
          for i in range(len(y) - 1) if y[i] < 0 <= y[i + 1]]
    fr = (len(zc) - 1) / (zc[-1] - zc[0])
    a, b = np.searchsorted(t, zc[0]), np.searchsorted(t, zc[-1])
    seg = y[a:b]
    sp = np.abs(np.fft.rfft(seg))
    per = len(zc) - 1                       # номер гармоники основной частоты
    h1 = sp[per]
    hv = math.sqrt(sum(sp[per * k] ** 2 for k in range(2, 10)))
    return dict(gen=True, Um=float(up), Umin=float(y.min()),
                U2=float(np.sqrt(np.mean(seg ** 2))),
                U3=float(np.sqrt(np.mean(z[a:b] ** 2))), f=float(fr), kg=float(hv / h1))


def rost(dolya, r=R, c=C):
    """время нарастания: от толчка до 0,9 установившейся амплитуды"""
    tt, vo, vp = generator(dolya, r, c)
    a = analiz(tt, vo, vp)
    if not a["gen"]:
        return None
    t90 = tt[np.argmax(np.abs(vo) > 0.9 * a["Um"])]
    return float(t90)


# ---------------------------------------------------------------------------
# режимы 2, 3: усилитель, вход — синус от генератора на «+» (C3 = 15 мкФ — перемычка)
# ---------------------------------------------------------------------------
def usilitel(roc, um, f=1e3, n=2000):
    t = np.arange(n) / n / f
    vin = um * np.sin(2 * math.pi * f * t)
    vo = np.array([v + u_os(v / R3, roc) for v in vin])
    vo = np.clip(vo, -VSAT, VSAT)
    sp = np.abs(np.fft.rfft(vo))
    kg = math.sqrt(sum(sp[k] ** 2 for k in range(2, 20))) / sp[1]
    return dict(Um=float(vo.max()), K=float(vo.max() / um), kg=float(kg),
                U2=float(np.sqrt(np.mean(vo ** 2))), t=t, vin=vin, vo=vo)


def k_usil_f(f, roc):
    """АЧХ усилителя режима 2 (малый сигнал, диоды закрыты), ОУ с одним полюсом"""
    k0 = 1 + roc / R3
    a = A0 / complex(1, f / (GBW / A0))
    h = a / (1 + a / k0)
    return abs(h), math.degrees(math.atan2(h.imag, h.real))


# ---------------------------------------------------------------------------
# интегратор и дифференциатор: вход — меандр, ОУ LM741 (полюс, SR, насыщение)
# ---------------------------------------------------------------------------
def meandr(t, um, f):
    return um if (t * f) % 1 < 0.5 else -um


def opamp_step(vo, vd, dt):
    wp = 2 * math.pi * GBW / A0
    dv = wp * (A0 * vd - vo)
    dv = max(-SR, min(SR, dv))
    return max(-VSAT, min(VSAT, vo + dv * dt))


INT = dict(R=10e3, C=10e-9, Rsh=1e6, Um=1.0, f=1e3)
DIF = dict(Rin=1e3, C=100e-9, Rf=10e3, Um=0.5, f=1e3)


def integrator(p=INT, perioda=40, dt=20e-9):
    """«−» : (Uвх − U−)/R = C·d(U− − Uвых)/dt + (U− − Uвых)/Rш"""
    nst = int(perioda / p["f"] / dt)
    vc, vo = 0.0, 0.0                        # vc = U− − Uвых
    out = []
    for k in range(nst):
        t = k * dt
        vin = meandr(t, p["Um"], p["f"])
        vm = vo + vc
        vc += ((vin - vm) / p["R"] - vc / p["Rsh"]) / p["C"] * dt
        vo = opamp_step(vo, -(vo + vc), dt)
        if t >= (perioda - 2) / p["f"] and k % 50 == 0:
            out.append((t, vin, vo))
    return np.array(out)


def differenciator(p=DIF, perioda=6, dt=5e-9):
    """Uвх → Rin → C → «−»; Rf с «−» на выход. Состояние: напряжение на C"""
    nst = int(perioda / p["f"] / dt)
    vc, vo, vm = 0.0, 0.0, 0.0
    out = []
    for k in range(nst):
        t = k * dt
        vin = meandr(t, p["Um"], p["f"])
        # ток через Rin–C равен току через Rf: (vin − vc − vm)/Rin = (vm − vo)/Rf
        vm = ((vin - vc) / p["Rin"] + vo / p["Rf"]) / (1 / p["Rin"] + 1 / p["Rf"])
        i = (vin - vc - vm) / p["Rin"]
        vc += i / p["C"] * dt
        vo = opamp_step(vo, -vm, dt)
        if t >= (perioda - 2) / p["f"] and k % 20 == 0:
            out.append((t, vin, vo))
    return np.array(out)


# ---------------------------------------------------------------------------
def main():
    z = {}
    print(f"f0 = {f0():.1f} Гц")

    # мост Вина
    F_B = [100, 200, 500, 1000, 1500, 2000, 2122, 2500, 3000, 5000, 10000, 20000]
    z["F_B"] = F_B
    z["BETA"] = [beta(x) for x in F_B]
    z["F0"] = f0()
    print("β(f):", [(x, round(b, 3), round(ph)) for x, (b, ph) in zip(F_B, z["BETA"])])

    # режим 1, положения R4
    DOLI = [0.45, 0.48, 0.52, 0.57, 0.65, 0.80, 1.00]
    z["DOLI"] = DOLI
    z["GEN"] = {}
    for d in DOLI:
        a = analiz(*generator(d))
        a["t90"] = rost(d) if a["gen"] else None
        z["GEN"][d] = a
        print(f"R4 {d:.0%}: ", {k: (round(v, 4) if isinstance(v, float) else v) for k, v in a.items()})

    # осциллограммы для графиков: пуск и установившийся режим при 57 %
    tt, vo, vp = generator(0.57, t_end=0.012, dt=0.5e-6)
    z["PUSK"] = [(float(t), float(a), float(b)) for t, a, b in zip(tt[::20], vo[::20], vp[::20])]

    # другая частота
    VAR = [("C1 = C2 = 10 нФ", R, 10e-9), ("C1 = C2 = 30 нФ", R, 30e-9),
           ("R1 = R2 = 10 кОм", 10e3, C), ("R1 = R2 = 2,5 кОм", 2.5e3, C)]
    z["FVAR"] = []
    for name, r, c in VAR:
        a = analiz(*generator(0.57, r, c, t_end=0.12 if c > 20e-9 or r > 6e3 else 0.06))
        z["FVAR"].append((name, r, c, f0(r, c), a["f"], a["Um"], a["U2"]))
        print(name, round(f0(r, c), 1), round(a["f"], 1), round(a["Um"], 3))

    # режим 2: линейная ООС
    U_VH = [0.1, 0.5, 1.0]
    DOLI2 = [0.25, 0.57, 1.00]
    z["U_VH"], z["DOLI2"] = U_VH, DOLI2
    z["US2"] = {}
    for d in DOLI2:
        for u in U_VH:
            a = usilitel(d * R4, u)
            z["US2"][(d, u)] = (a["Um"], a["K"], a["kg"], a["U2"])
            print(f"реж.2 {d:.0%} Uвх {u}: Um={a['Um']:.3f} K={a['K']:.2f} Kг={a['kg']:.3f}")
    F_K = [10, 100, 1e3, 10e3, 100e3, 300e3, 1e6]
    z["F_K"] = F_K
    z["K_F"] = [k_usil_f(x, 0.57 * R4) for x in F_K]

    # режим 3: без R4
    U_VH3 = [0.05, 0.1, 0.5, 1.0, 2.0]
    z["U_VH3"] = U_VH3
    z["US3"] = {}
    for u in U_VH3:
        a = usilitel(None, u)
        z["US3"][u] = (a["Um"], a["K"], a["kg"], a["U2"])
        print(f"реж.3 Uвх {u}: Um={a['Um']:.3f} K={a['K']:.2f} Kг={a['kg']:.3f}")
    for name, roc, u in [("OSC2", 0.57 * R4, 1.0), ("OSC3", None, 0.5)]:
        a = usilitel(roc, u)
        z[name] = [(float(t), float(x), float(y)) for t, x, y in zip(a["t"][::10], a["vin"][::10], a["vo"][::10])]

    # интегратор, дифференциатор
    z["INT"], z["DIF"] = INT, DIF
    oi = integrator()
    od = differenciator()
    z["OSC_INT"] = [tuple(map(float, r)) for r in oi]
    z["OSC_DIF"] = [tuple(map(float, r)) for r in od]
    z["INT_RAZMAH"] = (float(oi[:, 2].max()), float(oi[:, 2].min()))
    z["DIF_PIK"] = (float(od[:, 2].max()), float(od[:, 2].min()))
    print("интегратор: Uвых", z["INT_RAZMAH"], " идеал ±", INT["Um"] / (4 * INT["f"] * INT["R"] * INT["C"]))
    print("дифференциатор: пики", z["DIF_PIK"], " идеал ±", 2 * DIF["Um"] * DIF["Rf"] / DIF["Rin"])

    with open(LAB / "zamery.py", "w", encoding="utf-8") as fh:
        fh.write('"""ЛР-8: результаты расчёта (В, А, Гц, с). Пишет raschet_vse.py — руками не править."""\n')
        for k, v in z.items():
            fh.write(f"{k} = {v!r}\n")


if __name__ == "__main__":
    main()
