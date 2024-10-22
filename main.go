package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/4ster-light/phocus/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check for required privileges
	if runtime.GOOS != "windows" && os.Geteuid() != 0 {
		fmt.Println("This program requires root privileges to modify the hosts file and flush DNS.")
		fmt.Println("Please run it with sudo.")
		os.Exit(1)
	}

	// Initialize and run the application
	model := app.NewModel()
	program := tea.NewProgram(model)

	finalModel, err := program.Run()
	if err != nil {
		fmt.Println(app.ErrorStyle().Render("Error running program: " + err.Error()))
		os.Exit(1)
	}

	// Cleanup on exit
	if m, ok := finalModel.(app.Model); ok {
		if err := m.Cleanup(); err != nil {
			fmt.Println(app.ErrorStyle().Render("Error during cleanup: " + err.Error()))
		} else {
			fmt.Println(app.SuccessStyle().Render("\n✨ All domains have been successfully unblocked! ✨"))
		}
	}
}
