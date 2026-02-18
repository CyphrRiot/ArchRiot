#!/usr/bin/env bash
# ArchRiot — hyprlock-crypto.sh
# Show crypto price(s) for Hyprlock labels.
# Modes:
#   BTC|ETH|LTC|XMR|ZEC  -> single line, e.g. "BTC $91,350.50"
#   ROW                  -> one-line row: "BTC $.. • ZEC $.. • XMR $.. • LTC $.. • ETH $.."
#   ROWML                -> multi-line block with alignment and basic arrows
# Notes:
# - Uses CoinGecko simple/price (no API key). Caches to ~/.cache.
# - Adds litecoin to the default set so LTC is never "--".

set -euo pipefail

MODE="${1:-BTC}"

CACHE_DIR="$HOME/.cache"
CUR_FILE="$CACHE_DIR/hyprlock-crypto.json"
PREV_FILE="$CACHE_DIR/hyprlock-crypto-prev.json"
CACHE_TTL=1800 # 30 minutes
mkdir -p "$CACHE_DIR"

now=$(date +%s)
fetch_needed=true
if [[ -f "$CUR_FILE" ]]; then
  mtime=$(stat -c %Y "$CUR_FILE" 2>/dev/null || echo 0)
  age=$((now - mtime))
  if ((age < CACHE_TTL)); then
    fetch_needed=false
  fi
fi

# Always attempt to refresh cache (non-fatal); fallback to previous on failure
URL='https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum,litecoin,monero,zcash&vs_currencies=usd'
CURL_OPTS=(--fail --silent --show-error --max-time 6 --connect-timeout 2 --retry 2 --retry-delay 1 -H 'User-Agent: ArchRiot/hyprlock-crypto')
if curl "${CURL_OPTS[@]}" "$URL" -o "$CUR_FILE.tmp"; then
  mv -f "$CUR_FILE.tmp" "$CUR_FILE"
fi

python3 - "$MODE" "$CUR_FILE" "$PREV_FILE" <<'PY'
import json, sys, os, time
mode, cur_path, prev_path = sys.argv[1:4]
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

prices = {
    'BTC': cur.get('bitcoin', {}).get('usd'),
    'ZEC': cur.get('zcash', {}).get('usd'),
    'XMR': cur.get('monero', {}).get('usd'),
    'LTC': cur.get('litecoin', {}).get('usd'),
    'ETH': cur.get('ethereum', {}).get('usd'),
}

order_row = ['BTC','ZEC', 'XMR', 'LTC', 'ETH']

UP_EMO = "▲"
DOWN_EMO = "▼"
SAME_EMO = "•"

def fmt_amount(v: float, width: int = 10) -> str:
    if isinstance(v, (int, float)):
        return f"$ {v:>{width},.2f}"
    return "$ " + "--".rjust(width)

def fmt_one(sym: str) -> str:
    v = prices.get(sym)
    return f"{sym} {fmt_amount(v)}" if isinstance(v, (int, float)) else f"{sym} $ --"

if mode.upper() in ("BTC","ETH","LTC","XMR","ZEC"):
    s = mode.upper()
    print(fmt_one(s))
elif mode.upper() == "ROW":
    parts = [fmt_one(s) for s in order_row]
    print(" • ".join(parts))
elif mode.upper() == "ROWML":
    try:
        mtime = os.path.getmtime(cur_path)
        age_s = max(0, int(time.time() - mtime))
    except Exception:
        age_s = 0
    STALE = " ⌛" if age_s > 5400 else ""
    lines = []
    for s in order_row:
        v = prices.get(s)
        pv = prev.get(s)
        arrow = " •"
        if isinstance(v, (int, float)) and isinstance(pv, (int, float)):
            if v > pv: arrow = f" {UP_EMO}"
            elif v < pv: arrow = f" {DOWN_EMO}"
            else: arrow = f" {SAME_EMO}"
        lines.append(f"{s} {fmt_amount(v)}{arrow}{STALE}")
    print("\n".join(lines))
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
