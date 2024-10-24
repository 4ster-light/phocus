# Check if running as administrator
if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Warning "Please run this script as Administrator"
    exit 1
}

# Check if Go is installed
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Go is not installed. Please install Go before continuing."
    Write-Host "Visit https://golang.org/doc/install for installation instructions."
    exit 1
}

# Install phocus
Write-Host "Installing phocus..."
go install github.com/4ster-light/phocus@latest

if ($LASTEXITCODE -ne 0) {
    Write-Host "Installation failed"
    exit 1
}

# Get the binary path from GOPATH
$GOPATH = go env GOPATH
$SourceBinary = Join-Path $GOPATH "bin\phocus.exe"

# Ensure the binary exists
if (!(Test-Path $SourceBinary)) {
    Write-Host "Binary not found after installation"
    exit 1
}

# Move to Windows system directory
$TargetDir = "C:\Windows\System32"
$TargetPath = Join-Path $TargetDir "phocus.exe"

try {
    Move-Item -Path $SourceBinary -Destination $TargetPath -Force
    Write-Host "Phocus has been successfully installed to $TargetPath"
    Write-Host "Now you can run the program with (from an administrator shell): phocus"
    Write-Host "If you wish to remove the program run just delete the binary from $TargetPath"
} catch {
    Write-Host "Failed to move binary. Please ensure you're running as Administrator"
    exit 1
}
