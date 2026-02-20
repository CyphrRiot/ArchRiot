#!/usr/bin/env bash
# ArchRiot — hyprlock-crypto.sh (single-source config with optional blank separators + P/L)
# Modes:
#   <SYM> (e.g., BTC|ETH|LTC|XMR|ZEC) → "SYM $Price  [$±Gain $ / ±Pct%]" when entry provided
#   ROW                                 → single line: "SYM $.. • SYM $.. • ..." (separators and P/L skipped for compactness)
#   ROWML                               → multi-line block with arrows and staleness; if entry provided, appends P/L amount and %
# Implementation:
# - Exactly ONE ordered list defines everything. You can:
#   - Insert a blank line using ("", "")
#   - Add holdings and entry price by using a 4-tuple: ("SYM", "coingecko_id", held, entry_price)
#     Example: ("ZEC", "zcash", 100, 230.10)

set -euo pipefail

MODE="${1:-BTC}"

CACHE_DIR="$HOME/.cache"
CUR_FILE="$CACHE_DIR/hyprlock-crypto.json"
PREV_FILE="$CACHE_DIR/hyprlock-crypto-prev.json"
CACHE_TTL=1800 # 30 minutes
mkdir -p "$CACHE_DIR"

python3 - "$MODE" "$CUR_FILE" "$PREV_FILE" "$CACHE_TTL" <<'PY'
import json, sys, os, time, math, urllib.request

# SINGLE SOURCE OF TRUTH
# Order + ids; optional holdings and entry enable P/L display.
PAIRS = [
    ("ZEC", "zcash",   0,   0.0),   # edit held/entry to enable P/L; keep 0/0.0 to hide
    ("XMR", "monero",  0,   0.0),
    ("LTC", "litecoin", 0,   0.0),
    ("",   ""),                  # blank line between LTC and BTC (ROWML only)
    ("BTC", "bitcoin", 0,   0.0),
    ("ETH", "ethereum",0,   0.0),
]

mode, cur_path, prev_path, ttl_s = sys.argv[1], sys.argv[2], sys.argv[3], int(sys.argv[4])

# Normalize to list of dicts
items = []
for entry in PAIRS:
    if len(entry) == 2:
        sym, cid = entry
        items.append({"sym": sym, "cid": cid, "held": None, "entry": None})
    elif len(entry) == 4:
        sym, cid, held, entry_px = entry
        items.append({"sym": sym, "cid": cid, "held": held, "entry": entry_px})
    else:
        raise SystemExit("Invalid PAIRS entry; use (sym,cid) or (sym,cid,held,entry)")

# Build URL from non-blank ids only
ids = [it["cid"] for it in items if it["sym"] and it["cid"]]
url_ids = ",".join(ids)
URL = f"https://api.coingecko.com/api/v3/simple/price?ids={url_ids}&vs_currencies=usd"
UA  = "ArchRiot/hyprlock-crypto"

# Try to refresh cache non-fatally
try:
    req = urllib.request.Request(URL, headers={"User-Agent": UA})
    with urllib.request.urlopen(req, timeout=6) as resp:
        data = resp.read()
    tmp = cur_path + ".tmp"
    with open(tmp, "wb") as f:
        f.write(data)
    os.replace(tmp, cur_path)
except Exception:
    pass

# Load current and previous snapshots
try:
    with open(cur_path, 'r', encoding='utf-8') as f:
        cur = json.load(f)
except Exception:
    cur = {}
try:
    with open(prev_path, 'r', encoding='utf-8') as f:
        prev = json.load(f)
except Exception:
    prev = {}

# Current prices by symbol (non-blank only)
prices = {}
for it in items:
    sym, cid = it["sym"], it["cid"]
    if not sym or not cid:
        continue
    v = None
    try:
        v = cur.get(cid, {}).get('usd')
    except Exception:
        v = None
    prices[sym] = v

# Formatting helpers
UP_EMO = "▲"; DOWN_EMO = "▼"; SAME_EMO = "•"

def fmt_amount(v: float, width: int = 10) -> str:
    if isinstance(v, (int, float)) and not (v is None or math.isnan(v)):
        return f"$ {v:>{width},.2f}"
    return "$ " + "--".rjust(width)

def fmt_signed_amount(v: float, width: int = 8) -> str:
    if v is None or (isinstance(v, float) and (math.isnan(v) or math.isinf(v))):
        return "         --"
    sign = "+" if v >= 0 else "-"
    return f"{sign}$ {abs(v):>{width-3},.2f}"  # width minus sign+"$ "

def fmt_percent(p: float) -> str:
    if p is None or (isinstance(p, float) and (math.isnan(p) or math.isinf(p))):
        return "--%"
    sign = "+" if p >= 0 else "-"
    return f"{sign}{abs(p):.2f}%"

def fmt_one(sym: str, include_pl: bool = False) -> str:
    v = prices.get(sym)
    base = f"{sym} {fmt_amount(v)}"
    if include_pl:
        # Find this item's held/entry
        for it in items:
            if it["sym"] == sym:
                held = it["held"]; entry = it["entry"]; break
        else:
            held = entry = None
        if v is not None and entry and entry > 0:
            h = float(held) if held is not None else 1.0
            gl_amt = (float(v) - float(entry)) * h
            gl_pct = ((float(v) - float(entry)) / float(entry)) * 100.0
            return f"{base}  {fmt_signed_amount(gl_amt)} {fmt_percent(gl_pct)}"
    return base

# Helper symbol list (non-blank)
SYMS = [it["sym"] for it in items if it["sym"] and it["cid"]]

m = mode.upper()
if m in prices:
    print(fmt_one(m, include_pl=True))
elif m == "ROW":
    print(" • ".join(fmt_one(s, include_pl=False) for s in SYMS))
elif m == "ROWML":
    try:
        age_s = max(0, int(time.time() - os.path.getmtime(cur_path)))
    except Exception:
        age_s = 0
    STALE = " ⌛" if age_s > (ttl_s * 3 // 2) else ""
    lines = []
    for it in items:
        sym, cid = it["sym"], it["cid"]
        if not sym or not cid:
            lines.append("")
            continue
        v = prices.get(sym)
        pv = prev.get(sym)
        arrow = " •"
        if isinstance(v, (int, float)) and isinstance(pv, (int, float)):
            arrow = f" {UP_EMO}" if v > pv else (f" {DOWN_EMO}" if v < pv else f" {SAME_EMO}")
        line = f"{fmt_one(sym, include_pl=True)}{arrow}{STALE}"
        lines.append(line)
    print("\n".join(lines))
    # Save snapshot for next comparison (non-blank only)
    try:
        tmp = prev_path + ".tmp"
        with open(tmp, 'w', encoding='utf-8') as f:
            json.dump({k: v for k, v in prices.items() if isinstance(v, (int, float))}, f)
        os.replace(tmp, prev_path)
    except Exception:
        pass
else:
    print("--")
PY
