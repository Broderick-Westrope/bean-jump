package main

import (
	"doodle-jump/internal/tui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := tui.NewModel()

	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
