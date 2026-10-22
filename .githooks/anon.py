"""git filter: replaces personal data with placeholders on the way into git and back.

clean   (worktree -> index): real value  -> placeholder
smudge  (index -> worktree): placeholder -> real value

Mapping lives in .anon-map at the repo root (gitignored), one pair per line:
    Real Value = {{PLACEHOLDER}}
Lines starting with # are comments. Without the map, clean fails so nothing leaks.
"""
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parent.parent
MAP = ROOT / ".anon-map"


def load_pairs():
    pairs = []
    for line in MAP.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        real, placeholder = (part.strip() for part in line.rsplit("=", 1))
        if real and placeholder:
            pairs.append((real.encode("utf-8"), placeholder.encode("utf-8")))
    return pairs


def main():
    mode = sys.argv[1] if len(sys.argv) > 1 else ""
    data = sys.stdin.buffer.read()

    if not MAP.exists():
        if mode == "clean":
            sys.stderr.write("anon: .anon-map not found, refusing to stage university files\n")
            sys.exit(1)
        sys.stdout.buffer.write(data)
        return

    pairs = load_pairs()
    if mode == "clean":
        # longest first, so "Full Name" wins over "Surname"
        for real, placeholder in sorted(pairs, key=lambda p: -len(p[0])):
            data = data.replace(real, placeholder)
    elif mode == "smudge":
        for real, placeholder in sorted(pairs, key=lambda p: -len(p[1])):
            data = data.replace(placeholder, real)
    else:
        sys.stderr.write("usage: anon.py clean|smudge\n")
        sys.exit(2)
    sys.stdout.buffer.write(data)


if __name__ == "__main__":
    main()
