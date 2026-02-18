#!/usr/bin/env bash
# ArchRiot — hyprlock-crypto.sh (refactored)
# Single source of truth for crypto list and ids. No duplication.
# Modes:
#   BTC|ETH|LTC|XMR|ZEC  -> single line, e.g. "BTC $91,350.50"
#   ROW                  -> one-line row: "BTC $.. • ZEC $.. • XMR $.. • ETH $.. • LTC $.."
#   ROWML                -> multi-line block with alignment and basic arrows
# Notes:
# - Uses CoinGecko simple/price (no API key). Caches to ~/.cache.
# - Entire fetch/cache/format pipeline is implemented in Python stdlib (no curl/jq deps).

set -euo pipefail

MODE="${1:-BTC}"

CACHE_DIR="$HOME/.cache"
CUR_FILE="$CACHE_DIR/hyprlock-crypto.json"
PREV_FILE="$CACHE_DIR/hyprlock-crypto-prev.json"
CACHE_TTL=1800 # 30 minutes
mkdir -p "$CACHE_DIR"

python3 - "$MODE" "$CUR_FILE" "$PREV_FILE" "$CACHE_TTL" <<'PY'
import json, sys, os, time, urllib.request, urllib.error

# Single source of truth
ORDER = ["BTC", "ZEC", "XMR", "ETH", "LTC"]
ID_MAP = {
    "BTC": "bitcoin",
    "ZEC": "zcash",
    "XMR": "monero",
    "ETH": "ethereum",
    "LTC": "litecoin",
}

mode, cur_path, prev_path, ttl_s = sys.argv[1], sys.argv[2], sys.argv[3], int(sys.argv[4])

# Build URL dynamically from ORDER/ID_MAP
ids = ",".join(ID_MAP[s] for s in ORDER)
URL = f"https://api.coingecko.com/api/v3/simple/price?ids={ids}&vs_currencies=usd"
UA  = "ArchRiot/hyprlock-crypto"

# Decide if a refresh is needed based on cache mtime; we will try to refresh anyway, but this prevents over-frequent writes
now = int(time.time())
mtime = 0
try:
    mtime = int(os.path.getmtime(cur_path))
except Exception:
    mtime = 0

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
    # keep existing cache
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

# Build prices dict from ORDER/ID_MAP
prices = {}
for sym in ORDER:
    cid = ID_MAP[sym]
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
if m in ORDER:
    print(fmt_one(m))
elif m == "ROW":
    print(" • ".join(fmt_one(s) for s in ORDER))
elif m == "ROWML":
    try:
        age_s = max(0, int(time.time() - os.path.getmtime(cur_path)))
    except Exception:
        age_s = 0
    STALE = " ⌛" if age_s > (ttl_s * 3 // 2) else ""  # stale if > 1.5x TTL
    lines = []
    for s in ORDER:
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
