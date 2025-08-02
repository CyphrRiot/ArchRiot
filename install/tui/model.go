package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Color scheme
var (
	primaryColor = lipgloss.Color("#7aa2f7") // Tokyo Night blue
	accentColor  = lipgloss.Color("#bb9af7") // Tokyo Night purple
	successColor = lipgloss.Color("#9ece6a") // Tokyo Night green
	warningColor = lipgloss.Color("#e0af68") // Tokyo Night yellow
	errorColor   = lipgloss.Color("#f7768e") // Tokyo Night red
	bgColor      = lipgloss.Color("#1a1b26") // Tokyo Night background
	fgColor      = lipgloss.Color("#c0caf5") // Tokyo Night foreground
	dimColor     = lipgloss.Color("#565f89") // Tokyo Night comment
)

// ASCII Art
const ArchRiotASCII = `
  ▄▀█ █▀█ █▀▀ █ █ █▀█ █ █▀█ ▀█▀
  █▀█ █▀▄ █▄▄ █▀█ █▀▄ █ █▄█  █
`

// InstallModel represents the TUI model
type InstallModel struct {
	progress    float64
	message     string
	logs        []string
	maxLogs     int
	width       int
	height      int
	done        bool
	operation   string
	currentStep string
	inputMode   string // "git-username", "git-email", "reboot", ""
	inputValue  string // current typed input
	inputPrompt string // what we're asking for
	showConfirm    bool   // show YES/NO confirmation
	confirmPrompt  string // confirmation prompt text
	cursor         int    // 0 = YES, 1 = NO
}

// NewInstallModel creates a new installation model
func NewInstallModel() *InstallModel {
	return &InstallModel{
		logs:        make([]string, 0),
		maxLogs:     12,
		width:       80,
		height:      24,
		operation:   "ArchRiot Installation",
		currentStep: "Initializing...",
		inputMode:   "",
		inputValue:  "",
		inputPrompt: "",
		showConfirm:   false,
		confirmPrompt: "",
		cursor:        1, // Default to NO
	}
}

// Init implements tea.Model
func (m *InstallModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m *InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.showConfirm {
			switch msg.String() {
			case "left", "h":
				m.cursor = 0 // YES
				return m, nil
			case "right", "l":
				m.cursor = 1 // NO
				return m, nil
			case "enter", " ":
				return m.handleConfirmSelection()
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		} else if m.inputMode != "" {
			switch msg.String() {
			case "enter":
				return m.handleInputSubmit()
			case "ctrl+c":
				return m, tea.Quit
			case "backspace":
				if len(m.inputValue) > 0 {
					m.inputValue = m.inputValue[:len(m.inputValue)-1]
				}
				return m, nil
			default:
				if len(msg.String()) == 1 {
					m.inputValue += msg.String()
				}
				return m, nil
			}
		} else {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}

	case LogMsg:
		m.addLog(string(msg))
		return m, nil

	case ProgressMsg:
		m.setProgress(float64(msg))
		return m, nil

	case StepMsg:
		m.setCurrentStep(string(msg))
		return m, nil

	case DoneMsg:
		m.done = true
		m.setInputMode("reboot", "🔄 Reboot now? ")
		m.showConfirm = true
		m.confirmPrompt = "🔄 Reboot now?"
		m.cursor = 1 // Default to NO
		return m, nil

	case InputRequestMsg:
		m.setInputMode(msg.Mode, msg.Prompt)
		return m, nil

	case GitUsernameMsg:
		// Process git username input - handled by main
		return m, nil

	case GitEmailMsg:
		// Process git email input - handled by main
		return m, nil

	case GitConfirmMsg:
		// Git confirmation received, handled by main
		return m, nil

	case RebootMsg:
		if bool(msg) {
			return m, tea.Quit
		}
		return m, tea.Quit
	}

	return m, nil
}

// View implements tea.Model
func (m *InstallModel) View() string {
	var s strings.Builder

	// Clear screen on startup
	s.WriteString("\033[2J\033[H")


	// Header - ASCII + title + version (like Migrate) with spacing
	s.WriteString("\n") // Blank line before ASCII logo
	asciiStyle := lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	ascii := asciiStyle.Render(ArchRiotASCII)
	s.WriteString(ascii + "\n")

	titleStyle := lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	title := titleStyle.Render("ArchRiot Installer v"+GetVersion())
	s.WriteString(title + "\n")

	versionStyle := lipgloss.NewStyle().Foreground(dimColor)
	subtitle := versionStyle.Render("Charm with Bubbletea • Cypher Riot Themed")
	s.WriteString(subtitle + "\n\n")

	// Operation title
	operationStyle := lipgloss.NewStyle().Foreground(successColor).Bold(true)
	s.WriteString(operationStyle.Render("📦 "+m.operation) + "\n")

	// Info section - operation details
	infoStyle := lipgloss.NewStyle().Foreground(fgColor)
	logStyle := lipgloss.NewStyle().Foreground(dimColor)

	s.WriteString(infoStyle.Render("📋 Current Step:   "+m.currentStep) + "\n")
	s.WriteString(logStyle.Render("📝 Log File:       "+GetLogPath()) + "\n")

	// Progress bar
	s.WriteString("\n" + m.renderProgressBar() + "\n\n")

	// Scroll window - bordered content area
	s.WriteString(m.renderScrollWindow())

	// Confirmation below scroll window if shown
	if m.showConfirm {
		promptStyle := lipgloss.NewStyle().
			Foreground(fgColor).
			Bold(true)

		helpStyle := lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true)

		buttonRow := renderConfirmButtons(m.cursor)

		s.WriteString(fmt.Sprintf("\n\n%s  %s  %s",
			promptStyle.Render(m.confirmPrompt),
			buttonRow,
			helpStyle.Render("(← → to select, Enter to confirm)")))
	} else if m.inputMode != "" {
		s.WriteString("\n\n" + m.inputPrompt + m.inputValue + "_")
	} else {
		s.WriteString("\n\nPress 'q' to quit or 'ctrl+c' to exit")
	}

	return s.String()
}

// renderConfirmButtons creates styled YES/NO confirmation buttons
func renderConfirmButtons(cursor int) string {
	// Simple styled buttons that don't break layout
	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1a1b26")).
		Background(primaryColor).
		Padding(0, 2)

	unselectedStyle := lipgloss.NewStyle().
		Foreground(fgColor).
		Padding(0, 2)

	var yesButton, noButton string

	if cursor == 0 {
		yesButton = selectedStyle.Render("✓ YES")
		noButton = unselectedStyle.Render("✗ NO")
	} else {
		yesButton = unselectedStyle.Render("✓ YES")
		noButton = selectedStyle.Render("✗ NO")
	}

	// Create button row with proper spacing
	return lipgloss.JoinHorizontal(lipgloss.Center, yesButton, "   ", noButton)
}

func (m *InstallModel) renderProgressBar() string {
	progressStyle := lipgloss.NewStyle().Foreground(primaryColor).Bold(true)

	width := 50
	filled := int(m.progress * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	percentage := fmt.Sprintf("%.1f%%", m.progress*100)

	return progressStyle.Render(fmt.Sprintf("Progress: [%s] %s", bar, percentage))
}

// renderScrollWindow creates the bordered scroll window
func (m *InstallModel) renderScrollWindow() string {
	// Calculate responsive dimensions based on terminal size
	contentWidth := m.width - 4 // Account for borders and padding
	if contentWidth < 40 {
		contentWidth = 40 // Minimum width
	}

	// Calculate available height for scroll window
	// Account for: ASCII (7) + Title/subtitle (4) + Operation (2) + Progress (2) + Footer (3) + Buffer (2)
	usedHeight := 20
	availableHeight := m.height - usedHeight
	if availableHeight < 5 {
		availableHeight = 5 // Minimum scroll window height
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(contentWidth).
		Height(availableHeight)

	content := strings.Builder{}
	content.WriteString("Installation Log\n")
	separatorWidth := contentWidth - 2 // Account for padding
	if separatorWidth < 10 {
		separatorWidth = 10
	}
	content.WriteString(strings.Repeat("─", separatorWidth) + "\n")

	// Calculate how many log lines we can display
	maxLogLines := availableHeight - 2 // Account for header and separator
	if maxLogLines < 1 {
		maxLogLines = 1
	}

	// Show recent logs - limit to calculated space
	start := 0
	if len(m.logs) > maxLogLines {
		start = len(m.logs) - maxLogLines
	}

	actualLogCount := len(m.logs) - start
	if actualLogCount > maxLogLines {
		actualLogCount = maxLogLines
	}

	for i := start; i < start + actualLogCount; i++ {
		line := m.logs[i]
		maxLineWidth := contentWidth - 4 // Account for padding
		if maxLineWidth < 10 {
			maxLineWidth = 10
		}
		if len(line) > maxLineWidth {
			line = line[:maxLineWidth-3] + "..."
		}
		content.WriteString(line + "\n")
	}

	// Fill remaining lines in log area only
	for i := actualLogCount; i < maxLogLines; i++ {
		content.WriteString("\n")
	}

	return boxStyle.Render(content.String())
}

// addLog adds a new log entry
func (m *InstallModel) addLog(message string) {
	m.logs = append(m.logs, message)
}

// setProgress updates the progress value
func (m *InstallModel) setProgress(progress float64) {
	m.progress = progress
}

// setCurrentStep updates the current step
func (m *InstallModel) setCurrentStep(step string) {
	m.currentStep = step
}

// setInputMode sets the input mode and prompt
func (m *InstallModel) setInputMode(mode, prompt string) {
	if mode == "git-confirm" {
		m.showConfirm = true
		m.confirmPrompt = "🔧 Use these credentials?"
		m.cursor = 0 // Default to YES
	} else {
		m.inputMode = mode
		m.inputPrompt = prompt
		m.inputValue = ""
	}
}

// handleConfirmSelection processes YES/NO confirmation selection
func (m *InstallModel) handleConfirmSelection() (tea.Model, tea.Cmd) {
	if m.confirmPrompt == "🔄 Reboot now?" {
		// Reboot confirmation
		return m, tea.Quit
	} else if m.confirmPrompt == "🔧 Use these credentials?" {
		// Git credentials confirmation - send result back to main
		m.showConfirm = false
		m.confirmPrompt = ""
		// Signal completion through external callback
		if gitCompletionCallback != nil {
			gitCompletionCallback(m.cursor == 0) // YES = 0, NO = 1
		}
		return m, nil
	}
	return m, nil
}

// handleInputSubmit processes submitted input
func (m *InstallModel) handleInputSubmit() (tea.Model, tea.Cmd) {
	switch m.inputMode {
	case "git-username":
		// Send username back to main and clear input
		inputValue := m.inputValue
		m.inputMode = ""
		m.inputPrompt = ""
		m.inputValue = ""
		if gitUsernameCallback != nil {
			gitUsernameCallback(inputValue)
		}
		m.setInputMode("git-email", "Git Email: ")
		return m, nil
	case "git-email":
		// Send email back to main and clear input
		inputValue := m.inputValue
		m.inputMode = ""
		m.inputPrompt = ""
		m.inputValue = ""
		if gitEmailCallback != nil {
			gitEmailCallback(inputValue)
		}
		return m, nil
	}
	return m, nil
}

// Accessor methods for external packages
func (m *InstallModel) AddLog(message string) {
	m.addLog(message)
}

func (m *InstallModel) SetProgress(progress float64) {
	m.setProgress(progress)
}

func (m *InstallModel) SetCurrentStep(step string) {
	m.setCurrentStep(step)
}

func (m *InstallModel) SetInputMode(mode, prompt string) {
	m.setInputMode(mode, prompt)
}

func (m *InstallModel) IsDone() bool {
	return m.done
}
