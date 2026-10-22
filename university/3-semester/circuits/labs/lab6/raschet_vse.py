"""ЛР-6: расчёт ВАХ транзистора 2N1132A по SPICE-модели из базы Multisim.

Модель (MSCOMP_S.PRD, Symmetry Design Systems) — Гуммеля–Пуна, p-n-p.
Считаем как n-p-n с модулями напряжений: |UБЭ|, |UКЭ| и т. д. — на приборах
Multisim у p-n-p будут те же числа (знак зависит от того, как подключены щупы).

Запуск: python raschet_vse.py  → печатает таблицы и пишет zamery.py
"""
import math
from pathlib import Path

import numpy as np

LAB = Path(__file__).resolve().parent

# .MODEL 2N1132A pnp (только то, что влияет на постоянный ток)
M = dict(IS=8.69081e-12, BF=354.256, NF=1.23906, VAF=117.414, IKF=1.40313,
         ISE=9.13878e-13, NE=1.52672, BR=16.818, NR=1.2902, VAR=1.01024,
         IKR=0.0999733, ISC=9.42699e-14, NC=2.26776, RB=3.75863, RE=0.0844064,
         RC=0.461029)
VT = 1.380649e-23 * 300.15 / 1.602176634e-19   # 27 °C, как в Multisim
GMIN = 1e-12


def ex(v, n):
    a = v / (n * VT)
    return math.exp(min(a, 80.0)) - 1.0


def jn(vbe, vbc):
    """токи внутреннего транзистора (ib, ic) по SPICE3 bjtload, без RB/RE/RC"""
    cbe = M["IS"] * ex(vbe, M["NF"])
    cbc = M["IS"] * ex(vbc, M["NR"])
    cben = M["ISE"] * ex(vbe, M["NE"])
    cbcn = M["ISC"] * ex(vbc, M["NC"])
    q1 = 1.0 / (1.0 - vbc / M["VAF"] - vbe / M["VAR"])
    q2 = cbe / M["IKF"] + cbc / M["IKR"]
    qb = q1 * (1.0 + math.sqrt(max(0.0, 1.0 + 4.0 * q2))) / 2.0
    ib = cbe / M["BF"] + cben + cbc / M["BR"] + cbcn + GMIN * (vbe + vbc)
    ic = (cbe - cbc) / qb - cbc / M["BR"] - cbcn - GMIN * vbc
    return ib, ic


def tranzistor(ube, uce):
    """внешние UБЭ, UКЭ (эмиттер = 0) → (IБ, IК, IЭ) с учётом RB, RC, RE"""
    x = np.array([ube, uce, 0.0])      # внутренние узлы b', c', e'

    def F(x):
        b, c, e = x
        ib, ic = jn(b - e, b - c)
        return np.array([(ube - b) / M["RB"] - ib,
                         (uce - c) / M["RC"] - ic,
                         e / M["RE"] - ib - ic])

    for _ in range(200):
        f0 = F(x)
        J = np.empty((3, 3))
        for k in range(3):
            d = np.zeros(3)
            d[k] = 1e-7
            J[:, k] = (F(x + d) - f0) / 1e-7
        dx = np.linalg.solve(J, -f0)
        m = np.max(np.abs(dx))
        if m > 0.02:
            dx *= 0.02 / m
        x = x + dx
        if m < 1e-12:
            break
    b, c, e = x
    ib = (ube - b) / M["RB"]
    ic = (uce - c) / M["RC"]
    return float(ib), float(ic), float(ib + ic)


def ube_dlya(cel, uce_ili_ucb, ob=False, kakoy=0):
    """подобрать UБЭ (UЭБ в ОБ), при котором ток номер kakoy (0 — IБ, 2 — IЭ) = cel"""
    lo, hi = 0.0, 1.0
    for _ in range(80):
        mid = (lo + hi) / 2
        uce = mid + uce_ili_ucb if ob else uce_ili_ucb
        if tranzistor(mid, uce)[kakoy] < cel:
            lo = mid
        else:
            hi = mid
    return (lo + hi) / 2


# --- точки таблиц ----------------------------------------------------------
UBE = [0.50, 0.55, 0.60, 0.62, 0.64, 0.66, 0.68, 0.70, 0.72]      # табл. 1
UCE_VH = [15, 10, 5]
IB_VYH = [0.10e-3, 0.15e-3, 0.20e-3]                                 # табл. 2
UCE = [0, 0.05, 0.1, 0.2, 0.3, 0.5, 1, 2, 5, 10, 15]
UEB = [0.50, 0.55, 0.60, 0.62, 0.64, 0.66, 0.68, 0.70, 0.72]  # табл. 3
UCB_VH = [15, 10, 5]
IE_VYH = [2e-3, 4e-3, 8e-3]                                          # табл. 4
UCB = [0, 0.5, 1, 2, 5, 10, 15]


def vse():
    oe_vh = {u: [tranzistor(v, u)[0] for v in UBE] for u in UCE_VH}
    oe_vyh = {}
    for ib in IB_VYH:
        col = []
        for u in UCE:
            ube = ube_dlya(ib, u)
            col.append((ube, tranzistor(ube, u)[1]))
        oe_vyh[ib] = col
    ob_vh = {u: [tranzistor(v, v + u)[2] for v in UEB] for u in UCB_VH}
    ob_vyh = {}
    for ie in IE_VYH:
        col = []
        for u in UCB:
            ueb = ube_dlya(ie, u, ob=True, kakoy=2)
            col.append((ueb, tranzistor(ueb, ueb + u)[1]))
        ob_vyh[ie] = col
    return oe_vh, oe_vyh, ob_vh, ob_vyh


def main():
    oe_vh, oe_vyh, ob_vh, ob_vyh = vse()
    print("Табл.1 ОЭ вход: UБЭ | IБ мА при UКЭ", UCE_VH)
    for i, v in enumerate(UBE):
        print(f"  {v:.2f}", *(f"{oe_vh[u][i] * 1e3:10.5f}" for u in UCE_VH))
    print("Табл.2 ОЭ выход: UКЭ | IК мА при IБ", IB_VYH, "(в скобках UБЭ)")
    for i, u in enumerate(UCE):
        print(f"  {u:5}", *(f"{oe_vyh[ib][i][1] * 1e3:9.3f} ({oe_vyh[ib][i][0]:.3f})"
                            for ib in IB_VYH))
    print("Табл.3 ОБ вход: UЭБ | IЭ мА при UКБ", UCB_VH)
    for i, v in enumerate(UEB):
        print(f"  {v:.2f}", *(f"{ob_vh[u][i] * 1e3:10.4f}" for u in UCB_VH))
    print("Табл.4 ОБ выход: UКБ | IК мА при IЭ", IE_VYH)
    for i, u in enumerate(UCB):
        print(f"  {u:5}", *(f"{ob_vyh[ie][i][1] * 1e3:9.4f} ({ob_vyh[ie][i][0]:.3f})"
                            for ie in IE_VYH))

    with open(LAB / "zamery.py", "w", encoding="utf-8") as fh:
        fh.write('"""ЛР-6: точки таблиц (А, В). Пишет raschet_vse.py — руками не править."""\n')
        for name, val in [("UBE", UBE), ("UCE_VH", UCE_VH), ("IB_VYH", IB_VYH),
                          ("UCE", UCE), ("UEB", UEB), ("UCB_VH", UCB_VH),
                          ("IE_VYH", IE_VYH), ("UCB", UCB),
                          ("OE_VH", oe_vh), ("OE_VYH", oe_vyh),
                          ("OB_VH", ob_vh), ("OB_VYH", ob_vyh)]:
            fh.write(f"{name} = {val!r}\n")


if __name__ == "__main__":
    main()
