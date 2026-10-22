"""ЛР-4: графики для расчётной части отчёта.

Сигнал для расчёта: импульсы 0…1 В, f = 2 кГц, коэффициент заполнения 50 % (q = 2),
импульс начинается в момент t = 0.
Экспериментальные значения (EXP) вписываются после замеров курсором; пока их нет,
график сравнения не строится.
"""
import numpy as np
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt

HERE = r"C:/Users/ryche/Desktop/goscout/university/3-semester/circuits/labs/lab4"
plt.rcParams["font.size"] = 12

A, F, Q = 1.0, 2.0, 2          # В, кГц, скважность
N = np.arange(0, 6)            # гармоники 0…5
CALC = np.where(N == 0, A / Q, 2 * A / (np.maximum(N, 1) * np.pi) * np.abs(np.sin(N * np.pi / Q)))
EXP = [0.488, 0.620, 0, 0.167, 0, 0.085]  # курсор анализатора: 093423, 091318, 091358, 093215


def vosstanovlennyy():
    t = np.linspace(0, 1.0, 4000)          # мс, 2 периода
    w = 2 * np.pi * F                      # рад/мс
    pulses = np.where((t % (1 / F)) < 0.5 / F, A, 0.0)
    s = A / 2 + 2 * A / np.pi * (np.sin(w * t) + np.sin(3 * w * t) / 3 + np.sin(5 * w * t) / 5)
    fig, ax = plt.subplots(figsize=(9, 4))
    ax.plot(t, pulses, color="0.6", lw=3, label="импульсы")
    ax.plot(t, s, color="black", lw=1.8, label="A0 + 1, 3, 5 гармоники")
    ax.set_xlabel("t, мс")
    ax.set_ylabel("u, В")
    ax.set_ylim(-0.3, 1.55)
    ax.grid(alpha=0.3)
    ax.legend(loc="upper center", ncol=2)
    fig.tight_layout()
    fig.savefig(HERE + "/grafik_vosstanovlennyy.png", dpi=150)
    plt.close(fig)


def fazovyy():
    n = np.array([1, 3, 5])
    fig, ax = plt.subplots(figsize=(9, 3.4))
    ax.stem(n * F, np.full(3, -90), linefmt="black", markerfmt="ko", basefmt="black")
    ax.set_xticks(N * F)
    ax.set_yticks([-90, 0])
    ax.set_ylim(-110, 20)
    ax.set_xlim(-0.5, 11)
    ax.set_xlabel("f, кГц")
    ax.set_ylabel("φₙ, град")
    ax.grid(alpha=0.3)
    fig.tight_layout()
    fig.savefig(HERE + "/grafik_fazovyy.png", dpi=150)
    plt.close(fig)


def sravnenie():
    if EXP is None:
        return
    fig, ax = plt.subplots(figsize=(9, 4))
    fr = N * F
    ax.bar(fr - 0.25, CALC, width=0.5, color="0.75", edgecolor="black", label="расчёт")
    ax.bar(fr + 0.25, EXP, width=0.5, color="black", label="эксперимент")
    ax.set_xticks(fr)
    ax.set_xlabel("f, кГц")
    ax.set_ylabel("Aₙ, В")
    ax.grid(alpha=0.3, axis="y")
    ax.legend()
    fig.tight_layout()
    fig.savefig(HERE + "/grafik_sravnenie.png", dpi=150)
    plt.close(fig)


if __name__ == "__main__":
    print("расчёт:", np.round(CALC, 3))
    vosstanovlennyy()
    fazovyy()
    sravnenie()
    print("ok")
