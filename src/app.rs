use crate::hosts;
use anyhow::Result;
use std::vec::Vec;

pub struct App {
    pub input: String,
    pub messages: Vec<String>,
    pub blocked_domains: Vec<String>,
}

impl App {
    pub fn new() -> Self {
        Self {
            input: String::new(),
            messages: vec!["Waiting for domains to block...".to_string()],
            blocked_domains: Vec::new(),
        }
    }

    pub async fn block_domain(&mut self, domain: String) -> Result<()> {
        let domain = domain.trim().to_string();
        if domain.is_empty() {
            return Ok(());
        }

        match hosts::block_domain(&domain).await {
            Ok(_) => {
                self.blocked_domains.push(domain.clone());
                self.add_message(format!("✓ Blocked: {}", domain));
            }
            Err(err) => {
                self.add_message(format!("✗ Error blocking {}: {}", domain, err));
            }
        }

        Ok(())
    }

    pub async fn cleanup(&mut self) -> Result<()> {
        hosts::unblock_domains(&self.blocked_domains).await
    }

    fn add_message(&mut self, message: String) {
        if self.messages.len() == 1 && self.messages[0] == "Waiting for domains to block..." {
            self.messages[0] = message;
        } else {
            self.messages.push(message);
        }
    }
}
