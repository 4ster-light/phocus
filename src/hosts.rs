use anyhow::{Context, Result};
use std::fs;
use std::path::PathBuf;
use std::time::Duration;
use tokio::time::sleep;
use crate::dns::flush_dns;
use std::io::Write;

fn get_hosts_path() -> PathBuf {
    if cfg!(windows) {
        PathBuf::from(r"C:\Windows\System32\drivers\etc\hosts")
    } else {
        PathBuf::from("/etc/hosts")
    }
}

pub async fn block_domain(domain: &str) -> Result<()> {
    let hosts_path = get_hosts_path();
    let entry = format!("127.0.0.1 {}\n", domain);

    fs::OpenOptions::new()
        .append(true)
        .open(&hosts_path)
        .context("Failed to open hosts file")?
        .write_all(entry.as_bytes())
        .context("Failed to write to hosts file")?;

    // Add a small delay before flushing DNS
    sleep(Duration::from_millis(100)).await;
    flush_dns().await
}

pub async fn unblock_domains(domains: &[String]) -> Result<()> {
    let hosts_path = get_hosts_path();
    let content = fs::read_to_string(&hosts_path).context("Failed to read hosts file")?;

    let new_content: String = content
        .lines()
        .filter(|line| !domains.iter().any(|domain| line.contains(domain)))
        .collect::<Vec<_>>()
        .join("\n");

    fs::write(&hosts_path, new_content).context("Failed to write hosts file")?;

    // Add a small delay before flushing DNS
    sleep(Duration::from_millis(100)).await;
    flush_dns().await
}
