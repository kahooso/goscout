"""ЛР-4: структурные схемы параллельного и последовательного анализаторов спектра."""
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from matplotlib.patches import FancyBboxPatch

HERE = r"C:/Users/ryche/Desktop/goscout/university/3-semester/circuits/labs/lab4"
plt.rcParams["font.size"] = 11


def box(ax, x, y, w, h, text):
    ax.add_patch(FancyBboxPatch((x, y), w, h, boxstyle="square,pad=0", fill=False, lw=1.5))
    ax.text(x + w / 2, y + h / 2, text, ha="center", va="center")


def arrow(ax, x1, y1, x2, y2):
    ax.annotate("", xy=(x2, y2), xytext=(x1, y1), arrowprops=dict(arrowstyle="->", lw=1.3))


# ---------- параллельный ----------
fig, ax = plt.subplots(figsize=(10, 5))
ax.set_xlim(0, 10)
ax.set_ylim(0, 5)
ax.axis("off")
ax.text(0.2, 2.5, "Вход", va="center")
arrow(ax, 0.8, 2.5, 1.5, 2.5)
ax.plot([1.5, 1.5], [0.7, 4.3], color="black", lw=1.3)
rows = [(4.3, "ПФ f₁"), (3.3, "ПФ f₂"), (2.3, "ПФ f₃"), (0.7, "ПФ fₙ")]
for y, name in rows:
    arrow(ax, 1.5, y, 2.3, y)
    box(ax, 2.3, y - 0.3, 1.6, 0.6, name)
    arrow(ax, 3.9, y, 4.6, y)
    box(ax, 4.6, y - 0.3, 1.5, 0.6, "Детектор")
    arrow(ax, 6.1, y, 7.2, y)
ax.text(3.1, 1.55, "⋮", ha="center", fontsize=16)
ax.text(5.35, 1.55, "⋮", ha="center", fontsize=16)
box(ax, 7.2, 0.4, 2.2, 4.2, "Индикатор\n(экран)")
fig.tight_layout()
fig.savefig(HERE + "/ris_18_parallelnyy.png", dpi=150)
plt.close(fig)

# ---------- последовательный ----------
fig, ax = plt.subplots(figsize=(11, 4.6))
ax.set_xlim(0, 11)
ax.set_ylim(0, 4.6)
ax.axis("off")
ax.text(0.1, 3.4, "Вход", va="center")
arrow(ax, 0.7, 3.4, 1.2, 3.4)
box(ax, 1.2, 3.0, 1.6, 0.8, "Смеситель")
arrow(ax, 2.8, 3.4, 3.4, 3.4)
box(ax, 3.4, 3.0, 2.0, 0.8, "Узкополосный\nфильтр (УПЧ)")
arrow(ax, 5.4, 3.4, 6.0, 3.4)
box(ax, 6.0, 3.0, 1.5, 0.8, "Детектор")
arrow(ax, 7.5, 3.4, 8.4, 3.4)
box(ax, 8.4, 1.0, 2.3, 3.2, "ЭЛТ (экран)\n\nверт. откл. —\nамплитуда\n\nгориз. откл. —\nчастота")
box(ax, 1.2, 1.2, 1.6, 0.8, "Гетеродин")
arrow(ax, 2.0, 2.0, 2.0, 3.0)
box(ax, 1.2, 0.1, 1.6, 0.7, "Генератор\nразвёртки")
arrow(ax, 2.0, 0.8, 2.0, 1.2)
ax.plot([2.8, 7.9], [0.45, 0.45], color="black", lw=1.3)
ax.plot([7.9, 7.9], [0.45, 1.6], color="black", lw=1.3)
arrow(ax, 7.9, 1.6, 8.4, 1.6)
ax.text(5.3, 0.55, "пилообразное напряжение", ha="center", fontsize=9)
fig.tight_layout()
fig.savefig(HERE + "/ris_19_posledovatelnyy.png", dpi=150)
print("ok")
