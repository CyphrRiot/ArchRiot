#!/bin/bash

# Get the status
STATUS=$(mullvad status)
LOCATION=""

if echo "$STATUS" | grep -q "Connected"; then
    LOCATION=$(echo "$STATUS" | grep "Relay:" | awk '{print $2}' | cut -d'-' -f2 | tr '[:lower:]' '[:upper:]')
fi

# Display all icon options
echo "Status: $STATUS"
echo ""
echo "Choose which icon displays best on your system:"
echo ""
if [ -n "$LOCATION" ]; then
    echo "1. 󰌆 $LOCATION  (Shield Lock)"
    echo "2. 󰒄 $LOCATION  (Network Lock)"
    echo "3. 󰌾 $LOCATION  (Lock)"
    echo "4. 󰿆 $LOCATION  (VPN)"
    echo "5.  $LOCATION  (Unicode Lock)"
else
    echo "1. 󰌉  (Shield Unlock)"
    echo "2. 󰒃  (Network Unlock)"
    echo "3. 󰌿  (Unlock)"
    echo "4. 󰀉  (Disconnect)"
    echo "5.   (Unicode Unlock)"
fi