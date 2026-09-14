#!/bin/bash
# ==============================================================================
# Script Name:  harden_network.sh
# Description:  Automates host-based firewall hardening (Zero Trust profile).
# Author:       Your Name
# ==============================================================================

# Root enforcement verification
if [ "$EUID" -ne 0 ]; then
  echo "[-] Critical Error: Please execute this automation script as root (sudo)."
  exit 1
fi

echo "=== Initiating Automated Network Firewall Hardening ==="

# 1. Reset firewall configurations to factory defaults safely
echo "[*] Resetting existing network tables..."
ufw --force reset

# 2. Enforce strict baseline Zero-Trust rules (Deny everything incoming)
echo "[*] Setting Global Intercept Default Policies (Deny Incoming)..."
ufw default deny incoming
ufw default allow outgoing

# 3. Selectively authorize secure administrative ingress ports
echo "[*] Provisioning explicit administrative access criteria..."
ufw allow 22/tcp comment 'Secure SSH Administrative Entry'
ufw allow 80/tcp comment 'Standard Web Payload Delivery'
ufw allow 443/tcp comment 'Secure TLS Web Ingress'

# 4. Firewalls Ignition
echo "[*] Activating host-based network protection rules..."
ufw --force enable

# 5. Diagnostic validation confirmation report
echo ""
echo "=== Final Operational Network Guard Security Profile ==="
ufw status verbose
