import datetime
from netmiko import ConnectHandler

# Define the network device parameters
switch_ios = {
    'device_type': 'cisco_ios',
    'host': '192.168.1.10',
    'username': 'admin',
    'password': 'YourSecurePassword123!',
    'secret': 'EnableSecretPassword',
}

def backup_config(device):
    # Establish secure SSH connection
    print(f"Connecting to {device['host']}...")
    net_connect = ConnectHandler(**device)

    # Enter enable mode
    net_connect.enable()

    # Run the backup command
    print("Fetching running-configuration...")
    output = net_connect.send_command('show running-config')

    # Generate a timestamped filename
    date_str = datetime.datetime.now().strftime("%Y-%m-%d")
    filename = f"backup_{device['host']}_{date_str}.txt"

    # Write the output to a file
    with open(filename, 'w') as f:
        f.write(output)

    print(f"Success! Configuration saved to {filename}")
    net_connect.disconnect()

if __name__ == "__main__":
    backup_config(switch_ios)
