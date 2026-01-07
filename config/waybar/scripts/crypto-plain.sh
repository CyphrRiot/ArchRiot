#!/usr/bin/env bash
# ArchRiot — crypto-plain.sh
# Plain-text crypto prices for Hyprlock/Waybar (no JSON output)
#
# - Reads ~/.config/archriot/crypto.cfg for a list of up to 5 assets (space- or comma-separated)
#   Use CoinGecko IDs (e.g., bitcoin, ethereum, solana). Defaults: "bitcoin ethereum".
# - Fetches USD prices via CoinGecko (no API key).
# - Requires: curl, jq. If missing or any error occurs, prints an empty line (no error spam).
# - Output: a single line, e.g., "BTC 43.2k • ETH 2.34k"
#
# Safe, idempotent, suitable for Waybar/Hyprlock `exec =` and Hyprlock `text = cmd[...]` usage.

set -euo pipefail

CFG="$HOME/.config/archriot/crypto.cfg"
DEFAULT_IDS="bitcoin ethereum zcash monero"

have() { command -v "$1" >/dev/null 2>&1; }

# Read config (up to 5 ids), allow commas or spaces, normalize to lowercase
read_ids() {
  local raw ids
  if [[ -f "$CFG" ]]; then
    raw="$(tr ',' ' ' <"$CFG" | tr -s '[:space:]' ' ' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')"
    # shellcheck disable=SC2206
    arr=($raw)
    ids=""
    local count=0
    for id in "${arr[@]}"; do
      [[ -n "$id" ]] || continue
      ids+="${id,,} "
      count=$((count+1))
      [[ $count -ge 5 ]] && break
    done
    ids="$(echo -n "$ids" | sed 's/[[:space:]]\+$//')"
    echo "${ids:-$DEFAULT_IDS}"
  else
    echo "$DEFAULT_IDS"
  fi
}

# Map CoinGecko id -> ticker symbol (fallback: first 4 chars uppercased)
ticker() {
  case "$1" in
    bitcoin) echo "BTC" ;;
    ethereum) echo "ETH" ;;
    solana) echo "SOL" ;;
    cardano) echo "ADA" ;;
    ripple|xrp) echo "XRP" ;;
    dogecoin) echo "DOGE" ;;
    polygon|matic) echo "MATIC" ;;
    litecoin) echo "LTC" ;;
    zcash) echo "ZEC" ;;
    monero) echo "XMR" ;;
    polkadot) echo "DOT" ;;
    chainlink) echo "LINK" ;;
    avalanche-2|avalanche) echo "AVAX" ;;
    tron) echo "TRX" ;;
    stellar) echo "XLM" ;;
    binancecoin) echo "BNB" ;;
    *) printf "%s" "$1" | cut -c1-4 | tr '[:lower:]' '[:upper:]' ;;
  esac
}

# Format a numeric USD value with commas and 2 decimals (no $ sign)
fmt_usd() {
  local v="$1"
  # numeric?
  if ! printf "%s" "$v" | grep -qE '^[0-9]+(\.[0-9]+)?$'; then
    printf "%s" "$v"
    return
  fi
  awk -v n="$v" '
  function commify(x,    s,neg,int,frac,ret) {
    s = sprintf("%.2f", x)
    neg = (s ~ /^-/)
    if (neg) s = substr(s, 2)
    split(s, a, ".")
    int = a[1]; frac = a[2]
    ret = ""
    while (length(int) > 3) {
      ret = "," substr(int, length(int)-2) ret
      int = substr(int, 1, length(int)-3)
    }
    ret = int ret
    if (neg) ret = "-" ret
    return ret "." frac
  }
  BEGIN { print commify(n) }'
}

# Check tool availability; defer placeholder handling until after IDS/ROW are computed
TOOLS_OK=1
if ! have curl || ! have jq; then
  TOOLS_OK=0
fi

# Optional row selector: --row 1 (first pair) | --row 2 (second pair)
ROW=""
if [[ "${1:-}" == "--row" && "${2:-}" =~ ^[12]$ ]]; then
  ROW="${2}"
fi

IDS="$(read_ids)"

# If a row is requested, select two coins from the configured list:
# Row 1 -> first two (e.g., BTC/ETH), Row 2 -> next two (e.g., ZEC/XMR)
if [[ -n "$ROW" ]]; then
  # shellcheck disable=SC2206
  tmp_arr=($IDS)
  if [[ "$ROW" == "1" ]]; then
    IDS="${tmp_arr[0]:-} ${tmp_arr[1]:-}"
  else
    IDS="${tmp_arr[2]:-} ${tmp_arr[3]:-}"
  fi
  IDS="$(echo "$IDS" | sed 's/[[:space:]]\+/ /g' | sed 's/^ //;s/ $//')"
fi

# Build API URL
API="https://api.coingecko.com/api/v3/simple/price?vs_currencies=usd&ids=$(echo "$IDS" | tr ' ' ',' | tr '[:upper:]' '[:lower:]')"

# Fetch with short timeout; if tools missing or API empty, print aligned placeholders instead of empty
JSON=""
if [[ "$TOOLS_OK" -eq 1 ]]; then
  JSON="$(curl -fsSL --max-time 6 -H 'User-Agent: ArchRiot/crypto-plain' "$API" 2>/dev/null || true)"
fi
if [[ -z "$JSON" ]]; then
  OUT=""
  for id in ${IDS}; do
    sym="$(ticker "$id")"
    placeholder="--.--"
    if [[ -z "$OUT" ]]; then
      OUT="$sym $placeholder"
    else
      OUT="$OUT • $sym $placeholder"
    fi
  done
  echo "$OUT"
  exit 0
fi

# Assemble output
OUT=""
# shellcheck disable=SC2206
for id in ${IDS}; do
  # Pull price; jq returns null or empty if id missing
  usd="$(echo "$JSON" | jq -r --arg id "$id" '.[$id].usd // empty')"
  [[ -n "$usd" ]] || continue
  price="$(fmt_usd "$usd")"
  sym="$(ticker "$id")"
  if [[ -z "$OUT" ]]; then
    OUT="$sym $price"
  else
    OUT="$OUT • $sym $price"
  fi
done

echo "${OUT}"
