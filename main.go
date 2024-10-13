package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#04B575",
		Dark:  "#04B575",
	}).Render
)

type model struct {
	textInput      textinput.Model
	viewport       viewport.Model
	blockedDomains []string
	err            error
	content        string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter domain to block (e.g., example.com)"
	ti.Focus()

	vp := viewport.New(80, 20)
	vp.SetContent("Blocked domains will appear here.\n")

	return model{
		textInput: ti,
		viewport:  vp,
		content:   "Blocked domains will appear here.\n",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			domain := strings.TrimSpace(m.textInput.Value())
			if domain != "" {
				m.blockedDomains = append(m.blockedDomains, domain)
				if err := blockDomain(domain); err != nil {
					m.err = err
					m.content += fmt.Sprintf("Error blocking %s: %v\n", domain, err)
				} else {
					m.content += fmt.Sprintf("Blocked: %s\n", domain)
				}
				m.textInput.SetValue("")
				m.viewport.SetContent(m.content)
				m.viewport.GotoBottom()
			}
		}

	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return appStyle.Render(fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s",
		titleStyle.Render("DOMAIN BLOCKER"),
		m.viewport.View(),
		m.textInput.View(),
		"(esc to quit)",
	) + "\n")
}

func blockDomain(domain string) error {
	hostsFile := "/etc/hosts"
	if runtime.GOOS == "windows" {
		hostsFile = `C:\Windows\System32\drivers\etc\hosts`
	}

	f, err := os.OpenFile(hostsFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening hosts file: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("127.0.0.1 %s\n", domain)); err != nil {
		return fmt.Errorf("error writing to hosts file: %v", err)
	}

	// Add a small delay before flushing DNS
	time.Sleep(100 * time.Millisecond)

	return flushDNS()
}

func unblockDomains(domains []string) error {
	hostsFile := "/etc/hosts"
	if runtime.GOOS == "windows" {
		hostsFile = `C:\Windows\System32\drivers\etc\hosts`
	}

	input, err := os.ReadFile(hostsFile)
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
	if err := os.WriteFile(hostsFile, []byte(output), 0644); err != nil {
		return fmt.Errorf("error writing hosts file: %v", err)
	}

	// Add a small delay before flushing DNS
	time.Sleep(100 * time.Millisecond)

	return flushDNS()
}

func containsBlockedDomain(line string, domains []string) bool {
	for _, domain := range domains {
		if strings.Contains(line, domain) {
			return true
		}
	}
	return false
}

func flushDNS() error {
	var cmd *exec.Cmd
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ipconfig", "/flushdns")
	case "darwin":
		cmd = exec.Command("dscacheutil", "-flushcache")
		args = []string{"-c", "killall -HUP mDNSResponder"}
	case "linux":
		cmd = exec.Command("resolvectl", "flush-caches")
		if _, err := exec.LookPath("resolvectl"); err != nil {
			// Fallback for systems without systemd-resolved
			cmd = exec.Command("systemd-resolve", "--flush-caches")
			if _, err := exec.LookPath("systemd-resolve"); err != nil {
				// Final fallback
				cmd = exec.Command("sudo", "service", "network-manager", "restart")
			}
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error flushing DNS: %v, output: %s", err, string(output))
	}

	if runtime.GOOS == "darwin" {
		// Run additional command for macOS
		cmd = exec.Command("bash", args...)
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error flushing DNS (mDNSResponder): %v, output: %s", err, string(output))
		}
	}

	return nil
}

func main() {
	if runtime.GOOS != "windows" && os.Geteuid() != 0 {
		fmt.Println("This program requires root privileges to modify the hosts file and flush DNS.")
		fmt.Println("Please run it with sudo.")
		os.Exit(1)
	}

	m := initialModel()
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(model); ok {
		if err := unblockDomains(m.blockedDomains); err != nil {
			fmt.Println("Error unblocking domains:", err)
		} else {
			fmt.Println("All domains have been unblocked.")
		}
	}
}
