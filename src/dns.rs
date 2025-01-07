use anyhow::Result;
use std::process::Command;

pub async fn flush_dns() -> Result<()> {
    #[cfg(target_os = "windows")]
    {
        Command::new("ipconfig")
            .arg("/flushdns")
            .output()
            .context("Failed to flush DNS cache")?;
    }

    #[cfg(target_os = "macos")]
    {
        Command::new("dscacheutil")
            .arg("-flushcache")
            .output()
            .context("Failed to flush DNS cache")?;

        Command::new("killall")
            .args(["-HUP", "mDNSResponder"])
            .output()
            .context("Failed to restart mDNSResponder")?;
    }

    #[cfg(target_os = "linux")]
    {
        // Try to detect the init system
        let is_systemd = std::path::Path::new("/run/systemd/system").exists();

        if is_systemd {
            // Try systemd-specific methods first
            let systemd_methods = [
                ("systemd-resolve", &["--flush-caches"][..]),
                ("resolvectl", &["flush-caches"][..]),
            ];

            for (cmd, args) in systemd_methods {
                if let Ok(_) = Command::new(cmd).args(args).output() {
                    return Ok(());
                }
            }
        }

        // Try common DNS cache services
        let service_methods = [
            ("service", vec!["nscd", "restart"]),
            ("service", vec!["dnsmasq", "restart"]),
            ("service", vec!["named", "restart"]),
            ("systemctl", vec!["restart", "NetworkManager"]),
        ];

        for (cmd, args) in service_methods {
            if let Ok(output) = Command::new(cmd).args(&args).output() {
                if output.status.success() {
                    return Ok(());
                }
            }
        }

        // Try direct service control
        let direct_methods = [
            ("nscd", vec!["-K"]),
            ("killall", vec!["-HUP", "dnsmasq"]),
            ("rndc", vec!["flush"]),
        ];

        for (cmd, args) in direct_methods {
            if let Ok(output) = Command::new(cmd).args(&args).output() {
                if output.status.success() {
                    return Ok(());
                }
            }
        }

        // Method 1: Touch the hosts file
        let _ = Command::new("touch").arg("/etc/hosts").output();

        // Method 2: Notify NetworkManager if it exists
        if let Ok(_) = Command::new("nmcli").arg("general").arg("reload").output() {
            return Ok(());
        }
    }

    Ok(())
}
