# -*- coding: utf-8 -*-
"""ЛР-3. Расчёт таблиц АЧХ и построение графиков."""
import math, io, sys, os
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

U = 10.0
R = 1000.0
F = [200, 500, 1000, 2000, 5000, 8000, 16000, 30000, 60000, 150000]

def hp_rc(f, C):   # ФВЧ на RC: выход с резистора
    Xc = 1/(2*math.pi*f*C);  return U*R/math.hypot(R, Xc)

def rows(fn, par):
    out = []
    for f in F:
        u = fn(f, par)
        out.append([f, u, u/U])
    kmax = max(r[2] for r in out)
    for r in out:
        r.append(r[2]/kmax)
    return out

def fmtU(v):
    return f"{v*1000:.1f} мВ" if v < 1 else f"{v:.3f} В"

def show(title, data):
    print(f"\n{title}")
    print("| F, Гц | Uвх, В | Uвых | K = Uвых/Uвх | K̂ = K/Kmax |")
    print("|---|---|---|---|---|")
    for f, u, k, kn in data:
        fs = f"{f//1000} 000" if f >= 1000 else str(f)
        print(f"| {fs} | 10.000 | {fmtU(u)} | {k:.4f} | {kn:.4f} |")

d10 = rows(hp_rc, 10e-9)
d20 = rows(hp_rc, 20e-9)
show("Таблица 1 — RC-ФВЧ, C1 = 10 нФ, R1 = 1 кОм", d10)
show("Таблица 2 — RC-ФВЧ, C2 = 20 нФ, R1 = 1 кОм", d20)
print(f"\nf_гр(10 нФ) = {1/(2*math.pi*R*10e-9):.0f} Гц")
print(f"f_гр(20 нФ) = {1/(2*math.pi*R*20e-9):.0f} Гц")

# ---- график ----
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt

fig, ax = plt.subplots(figsize=(8, 5))
f10 = [r[0] for r in d10]
ax.semilogx(f10, [r[3] for r in d10], 'ko-', ms=5, lw=1.4, label='C = 10 нФ')
ax.semilogx(f10, [r[3] for r in d20], 'ks--', ms=5, lw=1.4, mfc='w', label='C = 20 нФ')
ax.axhline(0.707, color='gray', lw=0.9, ls=':')
ax.text(215, 0.73, '0,707', fontsize=9, color='gray')
for fg, st in ((15915, '-.'), (7958, '-.')):
    ax.axvline(fg, color='gray', lw=0.9, ls=st)
ax.text(16800, 0.05, '15,9 кГц', fontsize=9, color='gray')
ax.text(3100, 0.05, '7,96 кГц', fontsize=9, color='gray', ha='right')
ax.set_xlabel('Частота F, Гц')
ax.set_ylabel(r'$\hat{K} = K / K_{max}$')
ax.set_title('АЧХ RC-фильтра высоких частот')
ax.set_ylim(0, 1.05)
ax.grid(True, which='both', lw=0.4, color='0.85')
ax.legend(loc='upper left', frameon=False)
fig.tight_layout()
out = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'grafik_1_1_rc_fvch.png')
fig.savefig(out, dpi=200)
print("\nграфик:", out)
