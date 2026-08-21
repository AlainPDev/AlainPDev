import logging
from netmiko import ConnectHandler
from netmiko.exceptions import NetmikoTimeoutException, NetmikoAuthenticationException

# Configure logging for audit trails
logging.basicConfig(filename='network_changes.log', level=logging.INFO,
                    format='%(asctime)s - %(levelname)s - %(message)s')

# Target device definitions
device = {
    'device_type': 'cisco_ios',
    'host': '192.168.1.20',
    'username': 'admin',
    'password': 'YourSecurePassword123!',
    'secret': 'EnableSecretPassword',
}

# Target data to deploy (Can be easily moved to an external JSON or CSV file)
vlan_payload = [
    {"id": "10", "name": "Data_HR"},
    {"id": "20", "name": "Voice_VoIP"},
    {"id": "30", "name": "Guest_WiFI"}
]

def deploy_vlans(target_device, vlans):
    try:
        connection = ConnectHandler(**target_device)
        connection.enable()

        config_commands = []
        for vlan in vlans:
            # Simple input sanitization to prevent injection
            if not str(vlan['id']).isdigit():
                raise ValueError(f"Invalid VLAN ID: {vlan['id']}")

            config_commands.extend([
                f"vlan {vlan['id']}",
                f"name {vlan['name']}"
            ])

        # Deploy all configuration changes in a batch execution
        print(f"Deploying {len(vlans)} VLANs to {target_device['host']}...")
        output = connection.send_config_set(config_commands)

        # Save modifications permanently to NVRAM
        connection.send_command('write memory')

        logging.info(f"Successfully configured VLANs on {target_device['host']}. Details: \n{output}")
        print("Deployment successful. Logs updated.")

    except (NetmikoTimeoutException, NetmikoAuthenticationException) as net_err:
        logging.error(f"Network Connection Failed to {target_device['host']}: {net_err}")
        print(f"Network error encountered. Check logs for details.")
    except Exception as e:
        logging.error(f"Unexpected processing error: {e}")
        print(f"An unexpected failure occurred: {e}")
    finally:
        if 'connection' in locals():
            connection.disconnect()

if __name__ == "__main__":
    deploy_vlans(device, vlan_payload)
