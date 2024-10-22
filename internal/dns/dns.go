package dns

import (
	"fmt"
	"os/exec"
	"runtime"
)

func FlushDNS() error {
	switch runtime.GOOS {
	case "windows":
		return flushDNSWindows()
	case "darwin":
		return flushDNSDarwin()
	case "linux":
		return flushDNSLinux()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func flushDNSWindows() error {
	cmd := exec.Command("ipconfig", "/flushdns")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error flushing DNS: %v, output: %s", err, string(output))
	}
	return nil
}

func flushDNSDarwin() error {
	// First flush the DNS cache
	cmd := exec.Command("dscacheutil", "-flushcache")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error flushing DNS cache: %v, output: %s", err, string(output))
	}

	// Then restart the mDNSResponder
	cmd = exec.Command("killall", "-HUP", "mDNSResponder")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error restarting mDNSResponder: %v, output: %s", err, string(output))
	}

	return nil
}

func flushDNSLinux() error {
	// Try different methods in order of preference
	methods := []struct {
		name string
		cmd  []string
	}{
		{"systemd-resolve", []string{"systemd-resolve", "--flush-caches"}},
		{"resolvectl", []string{"resolvectl", "flush-caches"}},
		{"service", []string{"service", "nscd", "restart"}},
		{"nscd", []string{"nscd", "-K"}},
		{"dnsmasq", []string{"systemctl", "restart", "dnsmasq"}},
		{"network-manager", []string{"systemctl", "restart", "NetworkManager"}},
		{"pdnsd", []string{"systemctl", "restart", "pdnsd"}},
	}

	var lastErr error
	for _, method := range methods {
		// Check if the command exists
		if _, err := exec.LookPath(method.cmd[0]); err == nil {
			cmd := exec.Command(method.cmd[0], method.cmd[1:]...)
			if output, err := cmd.CombinedOutput(); err == nil {
				return nil
			} else {
				lastErr = fmt.Errorf("error with %s: %v, output: %s", method.name, err, string(output))
			}
		}
	}

	// If none of the above methods worked, try a simple hosts file reload
	// This is a fallback that should work on most systems
	if err := reloadHostsFile(); err != nil {
		return fmt.Errorf("all DNS flush methods failed. Last error: %v", lastErr)
	}

	return nil
}

func reloadHostsFile() error {
	// Simple touch of hosts file to trigger reload
	cmd := exec.Command("touch", "/etc/hosts")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error reloading hosts file: %v, output: %s", err, string(output))
	}
	return nil
}
