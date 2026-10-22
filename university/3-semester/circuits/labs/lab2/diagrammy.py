# Векторные диаграммы к ЛР-2 (полярные координаты), простой вид.
# Запуск: python diagrammy.py
import math
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
import numpy as np

INK = "#000000"
plt.rcParams.update({
    "font.family": "DejaVu Sans",
    "font.size": 10,
    "text.color": INK,
    "xtick.color": INK,
    "ytick.color": INK,
    "figure.facecolor": "white",
    "axes.facecolor": "white",
})


def num(fmt, value):
    return fmt.format(value).replace(".", ",")


def arrow(ax, angle_deg, r, label, off):
    th = math.radians(angle_deg)
    ax.annotate("", xy=(th, r), xytext=(0, 0),
                arrowprops=dict(arrowstyle="-|>", color=INK, linewidth=1.6,
                                shrinkA=0, shrinkB=0, mutation_scale=14))
    ax.annotate(label, xy=(th, r), xytext=off, textcoords="offset points",
                fontsize=10, ha="center", va="center")


def angle_arc(ax, a1, a2, r, text):
    th = np.radians(np.linspace(a1, a2, 60))
    ax.plot(th, [r] * 60, color=INK, linewidth=0.9)
    ax.annotate(text, xy=(math.radians((a1 + a2) / 2), r * 1.25),
                fontsize=10, ha="center", va="center")


def setup(ax, tmin, tmax, title):
    ax.set_thetamin(tmin)
    ax.set_thetamax(tmax)
    ax.set_rlim(0, 2.2)
    ax.set_rticks([1, 2])
    ax.set_rlabel_position(tmax - 8)
    ax.set_thetagrids(range(tmin, tmax + 1, 30), [""] * len(range(tmin, tmax + 1, 30)))
    ax.grid(color="#bbbbbb", linewidth=0.7)
    ax.set_title(title, fontsize=10, pad=16)


def rc_diagram():
    U, R, C, f = 2.0, 2000.0, 62e-9, 1000.0
    Xc = 1 / (2 * math.pi * f * C)
    I = U / math.hypot(R, Xc)
    UR, UC = I * R, I * Xc
    phi = math.degrees(math.atan(Xc / R))

    fig = plt.figure(figsize=(5.4, 4.8))
    ax = fig.add_subplot(111, projection="polar")
    setup(ax, -90, 30, "RC-цепь: R = 2 кОм, C = 62 нФ, f = 1 кГц")
    arrow(ax, 0, UR, num("$U_R$ = {:.2f} В", UR), (6, 14))
    arrow(ax, -90, UC, num("$U_C$ = {:.2f} В", UC), (-42, 4))
    arrow(ax, -phi, U, num("U = {:.2f} В", U), (16, -18))
    angle_arc(ax, -phi, 0, 0.55, num("φ = {:.0f}°", phi))
    fig.tight_layout()
    fig.savefig("diagramma_rc.png", dpi=200)


def rl_diagram():
    U, R, L, f = 2.0, 2000.0, 10e-3, 5000.0
    XL = 2 * math.pi * f * L
    I = U / math.hypot(R, XL)
    UR, UL = I * R, I * XL
    phi = math.degrees(math.atan(XL / R))

    fig = plt.figure(figsize=(5.4, 4.8))
    ax = fig.add_subplot(111, projection="polar")
    setup(ax, -30, 90, "RL-цепь: R = 2 кОм, L = 10 мГн, f = 5 кГц")
    arrow(ax, 0, UR, num("$U_R$ = {:.2f} В", UR), (4, -16))
    arrow(ax, 90, UL, num("$U_L$ = {:.2f} В", UL), (0, 16))
    arrow(ax, phi, U, num("U = {:.2f} В", U), (26, 10))
    angle_arc(ax, 0, phi, 1.3, num("φ = {:.1f}°", phi))
    fig.tight_layout()
    fig.savefig("diagramma_rl.png", dpi=200)


if __name__ == "__main__":
    rc_diagram()
    rl_diagram()
    print("готово")
