#!/bin/bash
# ==============================================================================
# Script Name:  provision_users.sh
# Description:  Automates provisioning of IAM users with secure group mapping.
# Author:       Your Name
# ==============================================================================

# Root enforcement verification
if [ "$EUID" -ne 0 ]; then
  echo "[-] Critical Error: Please execute this automation script as root (sudo)."
  exit 1
fi

# Define local target identities array (Format: username:group)
USERS=("cloud_admin:sysadmins" "cloud_auditor:compliance" "cloud_dev:developers")

echo "=== Starting Enterprise Identity Provisioning ==="

for entry in "${USERS[@]}"; do
    IFS=":" read -r username group <<< "$entry"

    # 1. Create regional operational group if non-existent
    if ! getent group "$group" > /dev/null 2>&1; then
        groupadd "$group"
        echo "[+] Group Created: $group"
    fi

    # 2. Safely provision system identity account
    if ! id "$username" > /dev/null 2>&1; then
        useradd -m -g "$group" -s /bin/bash "$username"
        echo "[+] Identity Provisioned: User '$username' mapped to Primary Group '$group'"

        # 3. Restrict home directory to Least-Privilege Standard (chmod 700)
        chmod 700 "/home/$username"
        chown "$username:$group" "/home/$username"
        echo "    [↳] Least-Privilege Enforced: /home/$username permission hardened to 700."
    else
        echo "[-] Conflict Warning: User account '$username' already exists. Skipping."
    fi
done

echo "=== Identity Provisioning Completed Successfully ==="
