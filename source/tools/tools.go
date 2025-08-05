package tools

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"archriot-installer/logger"
)

// Tool represents an individual tool
type Tool struct {
	ID          string
	Name        string
	Description string
	Category    string
	ExecuteFunc func() error
	Advanced    bool
	Available   bool
}

// ToolsModel represents the tools selection TUI model
type ToolsModel struct {
	tools       []Tool
	cursor      int
	selected    map[int]bool
	viewport    ViewportModel
	width       int
	height      int
	showDetails bool
	currentTool *Tool
}

// ViewportModel handles scrolling
type ViewportModel struct {
	top    int
	height int
}

const ArchRiotASCII = `
▄  ▄▀█ █▀█ █▀▀ █ █ █▀█ █ █▀█ ▀█▀  ▄
▄  █▀█ █▀▄ █▄▄ █▀█ █▀▄ █ █▄█  █   ▄
`

// Styles matching main installer (Tokyo Night theme)
var (
	primaryColor = lipgloss.Color("#7aa2f7") // Tokyo Night blue
	accentColor  = lipgloss.Color("#bb9af7") // Tokyo Night purple
	successColor = lipgloss.Color("#9ece6a") // Tokyo Night green
	warningColor = lipgloss.Color("#e0af68") // Tokyo Night yellow
	errorColor   = lipgloss.Color("#f7768e") // Tokyo Night red
	bgColor      = lipgloss.Color("#1a1b26") // Tokyo Night background
	fgColor      = lipgloss.Color("#c0caf5") // Tokyo Night foreground
	dimColor     = lipgloss.Color("#565f89") // Tokyo Night comment

	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	itemStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Padding(0, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(bgColor).
				Background(primaryColor).
				Padding(0, 2).
				Bold(true)

	descriptionStyle = lipgloss.NewStyle().
				Foreground(dimColor).
				Italic(true)

	categoryStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	asciiStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true).
			MarginTop(2)
)

// GetAvailableTools returns the registry of Go-implemented tools
func GetAvailableTools() ([]Tool, error) {
	tools := []Tool{
		{
			ID:          "secure-boot",
			Name:        "🛡️  Secure Boot Setup",
			Description: "Clean UEFI Secure Boot implementation using sbctl + shim-signed",
			Category:    "Security",
			ExecuteFunc: executeSecureBoot,
			Advanced:    true,
			Available:   checkSecureBootAvailable(),
		},
		{
			ID:          "memory-optimizer",
			Name:        "🧠 Memory Optimizer",
			Description: "Advanced memory management and optimization",
			Category:    "Performance",
			ExecuteFunc: executeMemoryOptimizer,
			Advanced:    false,
			Available:   true,
		},
		{
			ID:          "performance-tuner",
			Name:        "⚡ Performance Tuner",
			Description: "System performance optimization and tuning",
			Category:    "Performance",
			ExecuteFunc: executePerformanceTuner,
			Advanced:    true,
			Available:   true,
		},
		{
			ID:          "dev-environment",
			Name:        "💻 Development Environment",
			Description: "Complete development environment setup",
			Category:    "Development",
			ExecuteFunc: executeDevEnvironment,
			Advanced:    false,
			Available:   true,
		},
	}

	return tools, nil
}

// checkSecureBootAvailable checks if secure boot tools are available
func checkSecureBootAvailable() bool {
	// Check if sbctl is installed
	_, err := exec.LookPath("sbctl")
	return err == nil
}

// Tool execution functions
func executeSecureBoot() error {
	logger.LogMessage("INFO", "🚧 Work In Progress")
	logger.LogMessage("INFO", "This tool is not yet implemented")
	return fmt.Errorf("work in progress - not yet implemented")
}

func executeMemoryOptimizer() error {
	logger.LogMessage("INFO", "🚧 Work In Progress")
	logger.LogMessage("INFO", "This tool is not yet implemented")
	return fmt.Errorf("work in progress - not yet implemented")
}

func executePerformanceTuner() error {
	logger.LogMessage("INFO", "🚧 Work In Progress")
	logger.LogMessage("INFO", "This tool is not yet implemented")
	return fmt.Errorf("work in progress - not yet implemented")
}

func executeDevEnvironment() error {
	logger.LogMessage("INFO", "🚧 Work In Progress")
	logger.LogMessage("INFO", "This tool is not yet implemented")
	return fmt.Errorf("work in progress - not yet implemented")
}

// NewToolsModel creates a new tools model
func NewToolsModel() (*ToolsModel, error) {
	tools, err := GetAvailableTools()
	if err != nil {
		return nil, fmt.Errorf("failed to get available tools: %w", err)
	}

	return &ToolsModel{
		tools:       tools,
		cursor:      0,
		selected:    make(map[int]bool),
		viewport:    ViewportModel{top: 0, height: 10},
		width:       80,
		height:      24,
		showDetails: false,
	}, nil
}

// Init implements tea.Model
func (m ToolsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ToolsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.height = msg.Height - 8 // Reserve space for header and footer

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.viewport.top {
					m.viewport.top--
				}
			}

		case "down", "j":
			if m.cursor < len(m.tools)-1 {
				m.cursor++
				if m.cursor >= m.viewport.top+m.viewport.height {
					m.viewport.top++
				}
			}

		case "enter", " ":
			if len(m.tools) > 0 && m.cursor < len(m.tools) {
				tool := m.tools[m.cursor]
				if !tool.Available {
					logger.LogMessage("ERROR", fmt.Sprintf("Tool '%s' is not available", tool.Name))
					return m, nil
				}
				return m, m.executeTool(tool)
			}

		case "i", "d":
			m.showDetails = !m.showDetails
			if m.showDetails && len(m.tools) > 0 {
				m.currentTool = &m.tools[m.cursor]
			}

		case "r":
			// Refresh tools list
			tools, err := GetAvailableTools()
			if err == nil {
				m.tools = tools
				if m.cursor >= len(m.tools) {
					m.cursor = len(m.tools) - 1
					if m.cursor < 0 {
						m.cursor = 0
					}
				}
			}
		}
	}

	return m, nil
}

// View implements tea.Model
func (m ToolsModel) View() string {
	if len(m.tools) == 0 {
		return m.renderEmptyState()
	}

	if m.showDetails {
		return m.renderDetails()
	}

	return m.renderMainView()
}

// renderMainView renders the main tools selection view like Migrate
func (m ToolsModel) renderMainView() string {
	var s strings.Builder

	// Header with ASCII art like Migrate
	ascii := asciiStyle.Render(ArchRiotASCII)
	s.WriteString(ascii + "\n")

	title := titleStyle.Render("🔧 ArchRiot Tools")
	s.WriteString(title + "\n\n")

	// Menu options with beautiful styling like Migrate
	for i, tool := range m.tools {
		// Status indicator
		status := "✓"
		if !tool.Available {
			status = "✗"
		}

		// Build choice text
		choice := fmt.Sprintf("%s %s", status, tool.Name)

		if m.cursor == i {
			s.WriteString(selectedItemStyle.Render("❯ "+choice) + "\n")
		} else {
			s.WriteString(itemStyle.Render("  "+choice) + "\n")
		}
	}

	// Help text like Migrate
	help := helpStyle.Render("↑/↓: navigate • enter: select • i: details • r: refresh • q: quit")
	s.WriteString("\n" + help)

	// Center the content with border like Migrate
	content := borderStyle.Width(safeRenderWidth(m.width)).Render(s.String())
	return safeCenterContent(m.width, m.height, content)
}

// renderDetails renders the detailed view for a tool like Migrate
func (m ToolsModel) renderDetails() string {
	if m.currentTool == nil {
		return "No tool selected"
	}

	var s strings.Builder
	tool := *m.currentTool

	// Header with ASCII art like Migrate
	ascii := asciiStyle.Render(ArchRiotASCII)
	s.WriteString(ascii + "\n")

	title := titleStyle.Render("🔧 Tool Details")
	s.WriteString(title + "\n\n")

	// Tool information
	s.WriteString(headerStyle.Render("Tool: ") + tool.Name + "\n")
	s.WriteString(headerStyle.Render("Category: ") + categoryStyle.Render(tool.Category) + "\n")
	s.WriteString(headerStyle.Render("Description: ") + tool.Description + "\n")
	s.WriteString(headerStyle.Render("Implementation: ") + "Native Go function\n")

	if tool.Available {
		s.WriteString(headerStyle.Render("Status: ") + successStyle.Render("Ready to run") + "\n")
	} else {
		s.WriteString(headerStyle.Render("Status: ") + warningStyle.Render("Missing dependencies") + "\n")
	}

	if tool.Advanced {
		s.WriteString("\n")
		s.WriteString(warningStyle.Render("⚠ ADVANCED TOOL WARNING ⚠") + "\n")
		s.WriteString(warningStyle.Render("This tool modifies critical system components.") + "\n")
		s.WriteString(warningStyle.Render("Ensure you have backups before proceeding.") + "\n")
	}

	// Help text like Migrate
	help := helpStyle.Render("i/d: return to menu • enter: run tool • q: quit")
	s.WriteString("\n" + help)

	// Center the content with border like Migrate
	content := borderStyle.Width(safeRenderWidth(m.width)).Render(s.String())
	return safeCenterContent(m.width, m.height, content)
}

// renderEmptyState renders the empty state when no tools are found
func (m ToolsModel) renderEmptyState() string {
	var s strings.Builder

	// Header with ASCII art like Migrate
	ascii := asciiStyle.Render(ArchRiotASCII)
	s.WriteString(ascii + "\n")

	title := titleStyle.Render("🔧 ArchRiot Tools")
	s.WriteString(title + "\n\n")

	s.WriteString(warningStyle.Render("❌ No tools found in registry") + "\n\n")

	s.WriteString("🔧 Available tool categories:\n")
	s.WriteString(descriptionStyle.Render("   • Security     - UEFI Secure Boot, system hardening") + "\n")
	s.WriteString(descriptionStyle.Render("   • Performance  - Memory optimization, CPU tuning") + "\n")
	s.WriteString(descriptionStyle.Render("   • Development  - Development environment setup") + "\n")
	s.WriteString(descriptionStyle.Render("   • System       - General system configuration") + "\n")

	// Help text like Migrate
	help := helpStyle.Render("r: refresh • q: quit")
	s.WriteString("\n" + help)

	// Center the content with border like Migrate
	content := borderStyle.Width(safeRenderWidth(m.width)).Render(s.String())
	return safeCenterContent(m.width, m.height, content)
}

// executeTool runs the selected tool
func (m ToolsModel) executeTool(tool Tool) tea.Cmd {
	return func() tea.Msg {
		if !tool.Available {
			logger.LogMessage("ERROR", fmt.Sprintf("Tool is not available: %s", tool.Name))
			return nil
		}

		logger.LogMessage("INFO", fmt.Sprintf("Executing tool: %s", tool.Name))

		// Execute the Go function
		if err := tool.ExecuteFunc(); err != nil {
			logger.LogMessage("ERROR", fmt.Sprintf("Tool execution failed: %v", err))
		} else {
			logger.LogMessage("SUCCESS", fmt.Sprintf("Tool completed successfully: %s", tool.Name))
		}

		return nil
	}
}

// RunToolsInterface starts the tools selection interface
func RunToolsInterface() error {
	model, err := NewToolsModel()
	if err != nil {
		return fmt.Errorf("failed to create tools model: %w", err)
	}

	p := tea.NewProgram(*model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run tools interface: %w", err)
	}

	return nil
}

// safeRenderWidth calculates a safe border width for consistent rendering across terminals.
func safeRenderWidth(termWidth int) int {
	var margin int
	if termWidth < 70 {
		margin = 4
	} else if termWidth < 100 {
		margin = 6
	} else {
		margin = 8
	}

	width := termWidth - margin

	minWidth := 40
	if termWidth >= 60 {
		minWidth = 50
	}
	if termWidth >= 80 {
		minWidth = 60
	}

	if width < minWidth {
		width = minWidth
	}

	return width
}

// safeCenterContent safely centers content within the terminal bounds.
func safeCenterContent(termWidth, termHeight int, content string) string {
	lines := strings.Count(content, "\n") + 1

	if lines > termHeight-2 {
		return content
	}

	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, content)
}
