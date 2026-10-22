# Графики к ЛР-2. Запуск: python grafiki.py
# Данные — table1_rc.csv (таблица 1) и table2_rl.csv (таблица 2).
import csv
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from matplotlib.ticker import FixedLocator, FixedFormatter

BLUE, ORANGE, AQUA = "#2a78d6", "#eb6834", "#1baf7a"
INK, MUTED = "#0b0b0b", "#52514e"
STYLES = [(BLUE, "o", "-"), (ORANGE, "s", "--"), (AQUA, "^", ":")]

plt.rcParams.update({
    "font.family": "DejaVu Sans",
    "font.size": 9,
    "axes.edgecolor": MUTED,
    "axes.labelcolor": INK,
    "text.color": INK,
    "xtick.color": MUTED,
    "ytick.color": MUTED,
    "figure.facecolor": "white",
    "axes.facecolor": "white",
})


def read(path):
    with open(path, encoding="utf-8-sig") as f:
        rows = list(csv.reader(f))
    head = rows[0]
    fcol = [i for i, h in enumerate(head) if h.replace(".", "").isdigit()]
    freqs = [float(head[i]) for i in fcol]
    data = {r[0].replace('"', "").replace(",", ""): [float(r[i]) for i in fcol] for r in rows[1:]}
    return freqs, data


def panel(ax, freqs, data, keys, labels, title, ylabel):
    for (color, marker, ls), key, label in zip(STYLES, keys, labels):
        ax.plot(freqs, data[key], color=color, marker=marker, linestyle=ls,
                linewidth=2, markersize=5, label=label, clip_on=False, zorder=3)
    ax.set_title(title, fontsize=10, pad=8)
    ax.set_xlabel("Частота f, кГц")
    if ylabel:
        ax.set_ylabel(ylabel)
    ax.tick_params(axis="x", labelsize=8, rotation=45)
    ax.grid(True, color="#e3e2de", linewidth=0.8, zorder=0)
    ax.set_axisbelow(True)
    for side in ("top", "right"):
        ax.spines[side].set_visible(False)


def legend(fig, ax):
    handles, labels = ax.get_legend_handles_labels()
    fig.legend(handles, labels, loc="upper center", bbox_to_anchor=(0.5, 0.93),
               ncol=3, frameon=False, fontsize=9, columnspacing=2.5)


def fig_rc():
    freqs, data = read("table1_rc.csv")
    fig, axes = plt.subplots(1, 3, figsize=(10, 3.6), sharey=True)
    res = ["2,0 кОм", "2,5 кОм", "2,8 кОм"]
    labs = ["C1 = 3 нФ", "C2 = 62 нФ", "C3 = 1 мкФ"]
    for n, ax in enumerate(axes, start=1):
        keys = [f"R{n}C{i}" for i in (1, 2, 3)]
        panel(ax, freqs, data, keys, labs,
              f"R{n}: R0 + R = {res[n - 1]}", "$U_C$, В" if n == 1 else None)
        ax.set_xscale("log")
        ax.xaxis.set_major_locator(FixedLocator(freqs))
        ax.xaxis.set_major_formatter(FixedFormatter(
            [str(f).replace(".", ",").rstrip(",0") or "0" for f in freqs]))
        ax.xaxis.set_minor_locator(FixedLocator([]))
        ax.set_ylim(0, 2.15)
    legend(fig, axes[0])
    fig.suptitle("Напряжение на конденсаторе $U_C$ в зависимости от частоты  "
                 "(U = 2 В действующее)", fontsize=11, y=0.99)
    fig.tight_layout(rect=(0, 0, 1, 0.86))
    fig.savefig("grafik1_uc.png", dpi=200)


def fig_rl():
    freqs, data = read("table2_rl.csv")
    fig, axes = plt.subplots(1, 3, figsize=(10, 3.6), sharey=True)
    res = ["2,0 кОм", "2,5 кОм", "2,8 кОм"]
    labs = ["L1 = 3 мГн", "L2 = 7 мГн", "L3 = 10 мГн"]
    for n, ax in enumerate(axes, start=1):
        keys = [f"R{n}L{i}" for i in (1, 2, 3)]
        panel(ax, freqs, data, keys, labs,
              f"R{n}: R0 + R = {res[n - 1]}", "$U_L$, мВ" if n == 1 else None)
        ax.set_xlim(0, 5.3)
        ax.set_ylim(0, 330)
        ax.xaxis.set_major_locator(FixedLocator(freqs))
        ax.xaxis.set_major_formatter(FixedFormatter(
            [str(f).replace(".", ",").rstrip(",0") or "0" for f in freqs]))
    legend(fig, axes[0])
    fig.suptitle("Напряжение на катушке $U_L$ в зависимости от частоты  "
                 "(U = 2 В действующее)", fontsize=11, y=0.99)
    fig.tight_layout(rect=(0, 0, 1, 0.86))
    fig.savefig("grafik2_ul.png", dpi=200)


if __name__ == "__main__":
    fig_rc()
    fig_rl()
    print("готово: grafik1_uc.png, grafik2_ul.png")
