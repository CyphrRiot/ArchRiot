#!/usr/bin/env bash
# ArchRiot — hyprlock-crypto.sh (single-source config)
# Modes:
#   <SYM> (e.g., BTC|ETH|LTC|XMR|ZEC|SOL) → prints e.g. "BTC $12,345.67"
#   ROW                                   → single line: "SYM $.. • SYM $.. • ..."
#   ROWML                                 → multi-line block with arrows and staleness
# Implementation:
# - Exactly ONE ordered list defines both symbols and CoinGecko ids.
# - Everything (URL, ordering, parsing, output) derives from this single list.

set -euo pipefail

MODE="${1:-BTC}"

CACHE_DIR="$HOME/.cache"
CUR_FILE="$CACHE_DIR/hyprlock-crypto.json"
PREV_FILE="$CACHE_DIR/hyprlock-crypto-prev.json"
CACHE_TTL=1800 # 30 minutes
mkdir -p "$CACHE_DIR"

python3 - "$MODE" "$CUR_FILE" "$PREV_FILE" "$CACHE_TTL" <<'PY'
import json, sys, os, time, urllib.request

# SINGLE SOURCE OF TRUTH: ordered (SYMBOL, coingecko_id) pairs
PAIRS = [
    ("ZEC", "zcash"),
    ("XMR", "monero"),
    ("LTC", "litecoin"),
    ("BTC", "bitcoin"),
    ("ETH", "ethereum"),
    ("SOL", "solana"),
]

mode, cur_path, prev_path, ttl_s = sys.argv[1], sys.argv[2], sys.argv[3], int(sys.argv[4])
SYMS = [s for s, _ in PAIRS]
IDS  = [i for _, i in PAIRS]

# Build URL from single list
ids_param = ",".join(IDS)
URL = f"https://api.coingecko.com/api/v3/simple/price?ids={ids_param}&vs_currencies=usd"
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

# Build prices dict from single list
prices = {}
for sym, cid in PAIRS:
    v = None
    try:
        v = cur.get(cid, {}).get('usd')
    except Exception:
        v = None
    prices[sym] = v

# Formatting helpers
UP_EMO = "▲"; DOWN_EMO = "▼"; SAME_EMO = "•"

def fmt_amount(v: float, width: int = 10) -> str:
    if isinstance(v, (int, float)):
        return f"$ {v:>{width},.2f}"
    return "$ " + "--".rjust(width)

def fmt_one(sym: str) -> str:
    v = prices.get(sym)
    return f"{sym} {fmt_amount(v)}" if isinstance(v, (int, float)) else f"{sym} $ --"

m = mode.upper()
if m in SYMS:
    print(fmt_one(m))
elif m == "ROW":
    print(" • ".join(fmt_one(s) for s in SYMS))
elif m == "ROWML":
    try:
        age_s = max(0, int(time.time() - os.path.getmtime(cur_path)))
    except Exception:
        age_s = 0
    STALE = " ⌛" if age_s > (ttl_s * 3 // 2) else ""  # stale if > 1.5x TTL
    lines = []
    for s in SYMS:
        v = prices.get(s)
        pv = prev.get(s)
        arrow = " •"
        if isinstance(v, (int, float)) and isinstance(pv, (int, float)):
            arrow = f" {UP_EMO}" if v > pv else (f" {DOWN_EMO}" if v < pv else f" {SAME_EMO}")
        lines.append(f"{s} {fmt_amount(v)}{arrow}{STALE}")
    print("\n".join(lines))
    # Save snapshot for next comparison
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
