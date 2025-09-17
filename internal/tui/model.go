package tui

import (
	"bean-jump/internal/game"
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	game         *game.Game
	width        int
	height       int
	leftPressed  bool
	rightPressed bool
}

type tickMsg time.Time

func NewModel() Model {
	return Model{
		game:   game.NewGame(),
		width:  80,
		height: 24,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg { // 16ms = 60FPS (rounded)
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, 80)
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "left", "a":
			m.leftPressed = true
			m.rightPressed = false
		case "right", "d":
			m.rightPressed = true
			m.leftPressed = false
		case "r":
			if m.game.GameOver {
				m.game = game.NewGame()
			}
		}
		return m, nil

	case tickMsg:
		if !m.game.GameOver {
			m.game.Update(m.leftPressed, m.rightPressed)
		}
		// Reset movement flags after update
		m.leftPressed = false
		m.rightPressed = false
		return m, tickCmd()
	}

	return m, nil
}

func (m Model) View() string {
	if m.game.GameOver {
		return m.renderGameOver()
	}
	return lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).Render(m.renderGame())
}

func (m Model) renderGame() string {
	// Create a 2D grid to render the game
	grid := make([][]rune, m.height)
	colorGrid := make([][]string, m.height) // Track colors for each position
	for i := range grid {
		grid[i] = make([]rune, m.width)
		colorGrid[i] = make([]string, m.width)
		for j := range grid[i] {
			grid[i][j] = ' '
			colorGrid[i][j] = ""
		}
	}

	// Calculate scale factors to map game coordinates to terminal
	scaleX := float64(m.width) / game.GameWidth
	scaleY := float64(m.height) / game.GameHeight

	// Render platforms
	for _, platform := range m.game.Platforms {
		// Transform platform position relative to camera
		relativeY := platform.Position.Y - m.game.Camera.Y

		// Skip platforms outside view
		if relativeY < 0 || relativeY > game.GameHeight {
			continue
		}

		startX := int(platform.Position.X * scaleX)
		endX := int((platform.Position.X + platform.Width) * scaleX)
		y := int(relativeY * scaleY)

		if y >= 0 && y < m.height {
			var platformColor string
			platformChar := '-'
			if platform.Boost > 0 {
				platformColor = "2"
				platformChar = '='
			}
			for x := startX; x <= endX && x < m.width; x++ {
				if x >= 0 {
					grid[y][x] = platformChar
					colorGrid[y][x] = platformColor
				}
			}

			// Render boost level above the platform if it has boost
			if platform.Boost != 0 {
				centerX := int((platform.Position.X + platform.Width/2) * scaleX)
				powerUpY := y - 1 // Place boost level above platform
				if centerX >= 0 && centerX < m.width && powerUpY >= 0 && powerUpY < m.height {
					grid[powerUpY][centerX] = '⇧'
					grid[powerUpY][centerX+1] = []rune(strconv.Itoa(int(platform.Boost)))[0]
				}
			}
		}
	}

	// Render player
	playerRelativeY := m.game.Player.Position.Y - m.game.Camera.Y
	playerX := int(m.game.Player.Position.X * scaleX)
	playerY := int(playerRelativeY * scaleY)

	if playerX >= 0 && playerX < m.width && playerY >= 0 && playerY < m.height {
		grid[playerY][playerX] = 'O'
	}

	// Convert grid to string with colors
	result := ""
	for i, row := range grid {
		lineResult := ""
		for j, char := range row {
			if colorGrid[i][j] != "" {
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(colorGrid[i][j]))
				lineResult += style.Render(string(char))
			} else {
				lineResult += string(char)
			}
		}
		result += lineResult + "\n"
	}

	// Add score and instructions
	scoreStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	score := scoreStyle.Render(fmt.Sprintf("Score: %d", m.game.Score))

	instructions := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(
		"Use ← → or A/D to move | Q to quit")

	return score + "\n" + result + instructions
}

func (m Model) renderGameOver() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("204")).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width)

	gameOverText := style.Render("GAME OVER!")
	scoreText := style.Render(fmt.Sprintf("Final Score: %d", m.game.Score))
	restartText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center).
		Width(m.width).
		Render("Press R to restart | Q to quit")

	return "\n\n" + gameOverText + "\n" + scoreText + "\n\n" + restartText
}
