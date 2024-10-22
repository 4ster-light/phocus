package hosts

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/4ster-light/phocus/internal/dns"
)

// getHostsFilePath returns the appropriate hosts file path for the current OS
func getHostsFilePath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}

func BlockDomain(domain string) error {
	f, err := os.OpenFile(getHostsFilePath(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening hosts file: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("127.0.0.1 %s\n", domain)); err != nil {
		return fmt.Errorf("error writing to hosts file: %v", err)
	}

	// Add a small delay before flushing DNS
	time.Sleep(100 * time.Millisecond)

	return dns.FlushDNS()
}

func UnblockDomains(domains []string) error {
	input, err := os.ReadFile(getHostsFilePath())
	if err != nil {
		return fmt.Errorf("error reading hosts file: %v", err)
	}

	lines := strings.Split(string(input), "\n")
	var newLines []string

	for _, line := range lines {
		if !containsBlockedDomain(line, domains) {
			newLines = append(newLines, line)
		}
	}

	output := strings.Join(newLines, "\n")
	if err := os.WriteFile(getHostsFilePath(), []byte(output), 0644); err != nil {
		return fmt.Errorf("error writing hosts file: %v", err)
	}

	// Add a small delay before flushing DNS
	time.Sleep(100 * time.Millisecond)

	return dns.FlushDNS()
}

func containsBlockedDomain(line string, domains []string) bool {
	for _, domain := range domains {
		if strings.Contains(line, domain) {
			return true
		}
	}
	return false
}
