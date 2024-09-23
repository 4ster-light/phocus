import os
import time

HOSTS_FILE = "/etc/hosts"
REDIRECT_IP = "127.0.0.1"

def is_root():
    return os.geteuid() == 0

def get_websites():
    websites = input("Enter websites to block, separated by spaces: ").split()
    return websites

def block_websites(websites_to_block):
    if not is_root():
        print("This script must be run as root!")
        return

    with open(HOSTS_FILE, "r") as file:
        original_content = file.readlines()

    try:
        with open(HOSTS_FILE, "a") as file:
            for website in websites_to_block:
                file.write(f"{REDIRECT_IP}  {website}\n")
                print(f"Blocked: {website}")
    
        print("The websites have been blocked. Press Ctrl+C to unblock and exit.")

        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\nRestoring the original hosts file...")

        with open(HOSTS_FILE, "w") as file:
            file.writelines(original_content)

        print("Websites have been unblocked.")

if __name__ == "__main__":
    websites = get_websites()
    block_websites(websites)

