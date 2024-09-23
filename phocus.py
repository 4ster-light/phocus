import os
import subprocess
import sys
import time
import click

HOSTS_FILE = "/etc/hosts"
REDIRECT_IP = "127.0.0.1"

def is_root():
    return os.geteuid() == 0

@click.command()
@click.option("--websites", multiple=True, help="Websites to block (space-separated)")

def block_websites(websites):
    if not is_root():
        click.echo("Elevating privileges to root...")
        venv_python = sys.executable
        subprocess.run(['sudo', venv_python] + sys.argv)
        return

    with open(HOSTS_FILE, "r") as file:
        original_content = file.readlines()

    try:
        with open(HOSTS_FILE, "a") as file:
            for website in websites:
                file.write(f"{REDIRECT_IP}  {website}\n")
                click.echo(f"Blocked: {website}")
        
        click.echo("The websites have been blocked. Press Ctrl+C to unblock and exit.")
        
        while True:
            time.sleep(1)

    except KeyboardInterrupt:
        click.echo("\nRestoring the original hosts file...")
        
        with open(HOSTS_FILE, "w") as file:
            file.writelines(original_content)
        
        click.echo("Websites have been unblocked.")

if __name__ == "__main__":
    block_websites()

