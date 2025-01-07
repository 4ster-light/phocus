use std::process;
use anyhow::Result;

mod app;
mod dns;
mod hosts;
mod ui;

#[tokio::main]
async fn main() -> Result<()> {
    // Check for root/admin privileges
    #[cfg(not(windows))]
    if !nix::unistd::Uid::effective().is_root() {
        eprintln!("This program requires root privileges to modify the hosts file and flush DNS.");
        eprintln!("Please run it with sudo.");
        process::exit(1);
    }

    let mut app = app::App::new();
    if let Err(err) = ui::run_app(&mut app).await {
        eprintln!("Error running program: {}", err);
        process::exit(1);
    }

    // Cleanup on exit
    if let Err(err) = app.cleanup().await {
        eprintln!("Error during cleanup: {}", err);
        process::exit(1);
    }

    println!("✨ All domains have been successfully unblocked! ✨");
    Ok(())
}
