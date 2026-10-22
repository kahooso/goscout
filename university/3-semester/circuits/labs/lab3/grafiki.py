# -*- coding: utf-8 -*-
"""ЛР-3. Построение всех графиков АЧХ."""
import math, io, sys, os
sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt

def ZC(f, C): return 1/(1j*2*math.pi*f*C)
def ZL(f, L): return 1j*2*math.pi*f*L
def par(a, b): return (a*b)/(a+b)
R = 1000.0

def rc_hp(f, C, U=10): return abs(U*R/(R+ZC(f, C)))
def rc_lp(f, C, U=10): z = ZC(f, C); return abs(U*z/(R+z))
def rl_hp(f, L, U=10): z = ZL(f, L); return abs(U*z/(R+z))
def rl_lp(f, L, U=10): return abs(U*R/(R+ZL(f, L)))
def t_hp(f, C, U=10):
    z1 = z2 = ZC(f, C); zb = par(R, z2+R)
    return abs(U*zb/(R+z1+zb)*R/(z2+R))
def p_lp(f, C, U=10):
    z1 = z2 = ZC(f, C); zb = par(z2, R); za = par(z1, R+zb)
    return abs(U*za/(R+za)*zb/(R+zb))
def t_lp(f, L, U=10):
    z1 = z2 = ZL(f, L); zb = par(R, z2+R)
    return abs(U*zb/(R+z1+zb)*R/(z2+R))
def p_hp(f, L, U=10):
    z1 = z2 = ZL(f, L); zb = par(z2, R); za = par(z1, R+zb)
    return abs(U*za/(R+za)*zb/(R+zb))
def notch(f, C, L=405e-6, U=10): return abs(U*R/(R+par(ZL(f, L), ZC(f, C))+R))
def bp(f, C, L=9e-3, U=0.5):
    zs = par(par(ZL(f, L), ZC(f, C)), R); return abs(U*zs/(R+zs))

G1  = [200, 500, 1000, 2000, 5000, 8000, 16000, 30000, 60000, 150000]
N10 = [20000, 50000, 65000, 74000, 79000, 84000, 95000, 120000, 150000, 200000]
N20 = [15000, 35000, 46000, 52000, 56000, 60000, 68000, 85000, 120000, 200000]
B   = [2000, 5000, 8000, 12000, 17000, 22000, 30000, 50000, 80000, 150000]

def norm(fn, p, grid, U):
    k = [fn(f, p)/U for f in grid]
    m = max(k)
    return [v/m for v in k]

PLOTS = [
 ("grafik_1_rc_fvch",   "АЧХ RC-фильтра высоких частот",
  (rc_hp, 10e-9, G1, 10, "C = 10 нФ"), (rc_hp, 20e-9, G1, 10, "C = 20 нФ"),
  [(15915, "15,9 кГц", 'r'), (7958, "7,96 кГц", 'l')]),
 ("grafik_2_rl_fvch",   "АЧХ RL-фильтра высоких частот",
  (rl_hp, 10e-3, G1, 10, "L = 10 мГн"), (rl_hp, 20e-3, G1, 10, "L = 20 мГн"),
  [(15915, "15,9 кГц", 'r'), (7958, "7,96 кГц", 'l')]),
 ("grafik_3_rc_fnch",   "АЧХ RC-фильтра низких частот",
  (rc_lp, 10e-9, G1, 10, "C = 10 нФ"), (rc_lp, 20e-9, G1, 10, "C = 20 нФ"),
  [(15915, "15,9 кГц", 'r'), (7958, "7,96 кГц", 'l')]),
 ("grafik_4_rl_fnch",   "АЧХ RL-фильтра низких частот",
  (rl_lp, 10e-3, G1, 10, "L = 10 мГн"), (rl_lp, 20e-3, G1, 10, "L = 20 мГн"),
  [(15915, "15,9 кГц", 'r'), (7958, "7,96 кГц", 'l')]),
 ("grafik_5_t_fvch",    "АЧХ Т-образного фильтра высоких частот",
  (t_hp, 10e-9, G1, 10, "C = 10 нФ"), (t_hp, 20e-9, G1, 10, "C = 20 нФ"),
  [(17255, "17,3 кГц", 'r'), (8702, "8,7 кГц", 'l')]),
 ("grafik_6_t_fnch",    "АЧХ Т-образного фильтра низких частот",
  (t_lp, 10e-3, G1, 10, "L = 10 мГн"), (t_lp, 20e-3, G1, 10, "L = 20 мГн"),
  [(14514, "14,5 кГц", 'r'), (7261, "7,3 кГц", 'l')]),
 ("grafik_7_p_fvch",    "АЧХ П-образного фильтра высоких частот",
  (p_hp, 10e-3, G1, 10, "L = 10 мГн"), (p_hp, 20e-3, G1, 10, "L = 20 мГн"),
  [(17255, "17,3 кГц", 'r'), (8702, "8,7 кГц", 'l')]),
 ("grafik_8_p_fnch",    "АЧХ П-образного фильтра низких частот",
  (p_lp, 5e-9, G1, 10, "C = 5 нФ"), (p_lp, 10e-9, G1, 10, "C = 10 нФ"),
  [(29025, "29,0 кГц", 'r'), (14514, "14,5 кГц", 'l')]),
 ("grafik_9_rezhektornyy", "АЧХ режекторного фильтра",
  (notch, 10e-9, N10, 10, "C = 10 нФ"), (notch, 20e-9, N20, 10, "C = 20 нФ"),
  [(79085, "79,1 кГц", 'r'), (55921, "55,9 кГц", 'l')]),
 ("grafik_10_polosovoy", "АЧХ полосового фильтра",
  (bp, 10e-9, B, 0.5, "C = 10 нФ"), (bp, 20e-9, B, 0.5, "C = 20 нФ"),
  [(16776, "16,8 кГц", 'r'), (11862, "11,9 кГц", 'l')]),
]

base = os.path.dirname(os.path.abspath(__file__))
for name, title, a, b, marks in PLOTS:
    fig, ax = plt.subplots(figsize=(8, 5))
    for (fn, p, grid, U, lab), style in ((a, dict(ls='-', marker='o', mfc='k')),
                                         (b, dict(ls='--', marker='s', mfc='w'))):
        ax.semilogx(grid, norm(fn, p, grid, U), color='k', lw=1.4, ms=5,
                    label=lab, **style)
    ax.axhline(0.707, color='gray', lw=0.9, ls=':')
    ax.text(ax.get_xlim()[0]*1.08, 0.725, '0,707', fontsize=9, color='gray')
    for fv, lab, side in marks:
        ax.axvline(fv, color='gray', lw=0.9, ls='-.')
        ax.text(fv*(1.07 if side == 'r' else 0.93), 0.04, lab, fontsize=9,
                color='gray', ha='left' if side == 'r' else 'right')
    ax.set_xlabel('Частота F, Гц')
    ax.set_ylabel(r'$\hat{K} = K\,/\,K_{max}$')
    ax.set_title(title)
    ax.set_ylim(0, 1.06)
    ax.grid(True, which='both', lw=0.4, color='0.85')
    ax.legend(loc='best', frameon=False)
    fig.tight_layout()
    fig.savefig(os.path.join(base, name + '.png'), dpi=200)
    plt.close(fig)
    print("готов:", name + '.png')
