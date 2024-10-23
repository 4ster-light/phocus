package app

import (
	"fmt"
	"strings"

	"github.com/4ster-light/phocus/internal/hosts"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
)

// App state
type Model struct {
	textInput      textinput.Model
	viewport       viewport.Model
	blockedDomains []string
	err            error
	content        string
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter domain to block (e.g., example.com)"
	ti.Focus()
	ti.Prompt = "> "

	vp := viewport.New(78, 15)
	vp.Style = viewportStyle
	initialContent := "Waiting for domains to block..."
	vp.SetContent(initialContent)

	return Model{
		textInput: ti,
		viewport:  vp,
		content:   initialContent,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m.handleDomainInput()
		}

	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m Model) View() string {
	var builder strings.Builder

	// Add title
	builder.WriteString(titleStyle.Render("P H O C U S"))
	builder.WriteString("\n\n")

	// Add viewport with messages
	builder.WriteString(m.viewport.View())
	builder.WriteString("\n")

	// Add input field
	builder.WriteString(m.textInput.View())
	builder.WriteString("\n\n")

	// Add help text
	builder.WriteString(helpStyle.Render("(esc to quit)"))

	// Wrap everything in the app style
	return appStyle.Render(builder.String())
}

func (m Model) Cleanup() error {
	return hosts.UnblockDomains(m.blockedDomains)
}

func (m *Model) formatMessage(message string, isError bool, domain string) string {
	prefix := "✓"
	style := blockedDomainStyle
	if isError {
		prefix = "✗"
		style = errorStyle
		return style.Render(fmt.Sprintf("%s Error blocking %s: %s", prefix, domain, message))
	}
	return style.Render(fmt.Sprintf("%s Blocked: %s", prefix, domain))
}

func (m Model) handleDomainInput() (Model, tea.Cmd) {
	domain := strings.TrimSpace(m.textInput.Value())
	if domain == "" {
		return m, nil
	}

	// Add domain to blocked list
	m.blockedDomains = append(m.blockedDomains, domain)

	// Format the new message
	var newMessage string
	if err := hosts.BlockDomain(domain); err != nil {
		m.err = err
		newMessage = m.formatMessage(err.Error(), true, domain)
	} else {
		newMessage = m.formatMessage("", false, domain)
	}

	// Update content and viewport
	if m.content == "Waiting for domains to block..." {
		m.content = newMessage
	} else {
		m.content += "\n" + newMessage
	}

	// Reset input and update viewport
	m.textInput.SetValue("")
	m.viewport.SetContent(m.content)
	m.viewport.GotoBottom()

	return m, nil
}
