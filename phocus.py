import os
import subprocess
import sys
import time
import socket

from rich.console import Console
from rich.prompt import Prompt

HOSTS_FILE = "/etc/hosts"
REDIRECT_IP = "127.0.0.1"

console = Console()

def get_ip_addresses(domain):
    try:
        return socket.gethostbyname_ex(domain)[2]
    except socket.gaierror:
        return []

def block_websites():
    if not os.geteuid() == 0:
        console.print("[yellow]Elevating privileges to root...[/yellow]")
        venv_python = sys.executable
        subprocess.run(['sudo', venv_python] + sys.argv)
        return

    websites = []
    while True:
        website = Prompt.ask("Enter website to block (or press Enter to finish)", default=None)
        if not website:
            break
        websites.append(website)

    if not websites:
        console.print("[bold red]No websites provided. Exiting...[/bold red]")
        return

    with open(HOSTS_FILE, "r") as file:
        original_content = file.readlines()

    try:
        with open(HOSTS_FILE, "a") as file:
            for website in websites:
                # Block main domain
                file.write(f"{REDIRECT_IP} {website}\n")
                # Block www subdomain
                file.write(f"{REDIRECT_IP} www.{website}\n")
                # Block IP addresses associated with the domain
                ip_addresses = get_ip_addresses(website)
                for ip in ip_addresses:
                    file.write(f"{REDIRECT_IP} {ip}\n")
                console.print(f"[green]Blocked:[/green] {website} and its IP addresses")
        
        # Flush DNS cache
        if sys.platform == "darwin":  # macOS
            subprocess.run(["sudo", "killall", "-HUP", "mDNSResponder"])
        elif sys.platform == "linux":
            subprocess.run(["sudo", "systemd-resolve", "--flush-caches"])
        else:
            console.print("[bold red]Unsupported platform. Exiting...[/bold red]")
            return
        
        console.print("[bold cyan]The websites have been blocked. Press Ctrl+C to unblock and exit.[/bold cyan]")
        
        while True:
            time.sleep(1)

    except KeyboardInterrupt:
        console.print("\n[bold yellow]Restoring the original hosts file...[/bold yellow]")
        
        with open(HOSTS_FILE, "w") as file:
            file.writelines(original_content)
        
        console.print("[bold green]Websites have been unblocked.[/bold green]")

if __name__ == "__main__":
    block_websites()
