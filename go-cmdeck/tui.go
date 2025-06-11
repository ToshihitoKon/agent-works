package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TUI struct {
	executor *Executor
}

type model struct {
	executor     *Executor
	contexts     []Context
	cursor       int
	selected     map[int]struct{}
	currentView  string
	lastOutput   string
	showOutput   bool
	theme        ColorTheme
	width        int
	height       int
}

func getStyles(theme ColorTheme, width, height int) (titleStyle, selectedStyle, topPanelStyle, bottomPanelStyle, outputTitleStyle lipgloss.Style) {
	titleStyle = lipgloss.NewStyle().
		MarginLeft(2).
		Foreground(lipgloss.Color(theme.Title))

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Selected)).
		Bold(true)

	topHeight := height/2 - 2
	bottomHeight := height - topHeight - 4
	
	if topHeight < 8 {
		topHeight = 8
	}
	if bottomHeight < 5 {
		bottomHeight = 5
	}

	topPanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Border)).
		Padding(1, 2).
		Width(width - 4).
		Height(topHeight)

	bottomPanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Border)).
		Padding(1, 2).
		Width(width - 4).
		Height(bottomHeight)

	outputTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.OutputTitle)).
		Bold(true)

	return
}

func NewTUI(executor *Executor) *TUI {
	return &TUI{executor: executor}
}

func (t *TUI) Run() error {
	contexts := t.executor.listContexts()
	
	m := model{
		executor:    t.executor,
		contexts:    contexts,
		selected:    make(map[int]struct{}),
		currentView: "list",
		lastOutput:  "Ready to execute commands...",
		showOutput:  true,
		theme:       t.executor.config.Theme,
		width:       80,
		height:      24,
	}

	p := tea.NewProgram(&m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.contexts)-1 {
				m.cursor++
			}
		case " ":
			if len(m.contexts) > 0 {
				currentContextName := m.contexts[m.cursor].Name
				context := m.executor.config.Contexts[currentContextName]
				
				if runCmd, exists := context.Commands["run"]; exists {
					output, exitCode, err := m.executor.executeJobWithOutput(runCmd, context.Variables)
					
					result := &ExecutionResult{
						Timestamp: time.Now(),
						Success:   err == nil && exitCode == 0,
						ExitCode:  exitCode,
						Output:    output,
					}
					
					context.LastResult = result
					m.executor.config.Contexts[currentContextName] = context
					m.executor.config.save()
					
					oldCursor := m.cursor
					m.contexts = m.executor.listContexts()
					
					for i, ctx := range m.contexts {
						if ctx.Name == currentContextName {
							m.cursor = i
							break
						}
					}
					
					if m.cursor >= len(m.contexts) {
						m.cursor = oldCursor
						if m.cursor >= len(m.contexts) {
							m.cursor = len(m.contexts) - 1
						}
					}
				}
			}
		}
	}
	return m, nil
}

func (m *model) View() string {
	var topContent strings.Builder
	var bottomContent strings.Builder

	titleStyle, selectedStyle, topPanelStyle, bottomPanelStyle, outputTitleStyle := getStyles(m.theme, m.width, m.height)
	
	topHeight := m.height/2 - 2
	bottomHeight := m.height - topHeight - 4
	if topHeight < 8 {
		topHeight = 8
	}
	if bottomHeight < 5 {
		bottomHeight = 5
	}

	topContent.WriteString(titleStyle.Render("Job Deck"))
	topContent.WriteString("\n\n")

	if len(m.contexts) == 0 {
		topContent.WriteString("No contexts available.")
	} else {
		// Calculate available space for context list
		// topHeight is the content area height set by lipgloss
		// Account for: Title (1) + Empty line (1) + Help text (1) = 3
		availableLines := topHeight - 3
		if availableLines < 1 {
			availableLines = 1
		}
		
		// Show contexts around cursor position
		startIdx := 0
		endIdx := len(m.contexts)
		
		if len(m.contexts) > availableLines {
			startIdx = m.cursor - availableLines/2
			if startIdx < 0 {
				startIdx = 0
			}
			endIdx = startIdx + availableLines
			if endIdx > len(m.contexts) {
				endIdx = len(m.contexts)
				startIdx = endIdx - availableLines
				if startIdx < 0 {
					startIdx = 0
				}
			}
		}

		for i := startIdx; i < endIdx; i++ {
			context := m.contexts[i]
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			statusIcon := " "
			if context.LastResult != nil {
				if context.LastResult.Success {
					statusIcon = "✓"
				} else {
					statusIcon = "✗"
				}
			}

			style := lipgloss.NewStyle()
			if m.cursor == i {
				style = selectedStyle
			}

			line := fmt.Sprintf("%s [%s] %s", 
				cursor, statusIcon, context.Label)
			
			if context.Description != "" {
				line += fmt.Sprintf(" - %s", context.Description)
			}
			
			// Truncate long lines to fit within panel  
			// Width set in style - padding left/right (2*2=4)
			maxLineWidth := m.width - 4 - 4  // total width - borders - padding
			if maxLineWidth < 10 {
				maxLineWidth = 10
			}
			if len(line) > maxLineWidth {
				line = line[:maxLineWidth-3] + "..."
			}

			topContent.WriteString(style.Render(line))
			topContent.WriteString("\n")
		}
	}

	topContent.WriteString("\n↑/↓ or j/k: navigate • space: execute • q: quit")

	bottomContent.WriteString(outputTitleStyle.Render("Job Details"))
	bottomContent.WriteString("\n")
	bottomContent.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	var output string
	if len(m.contexts) > 0 && m.cursor < len(m.contexts) {
		selectedContext := m.contexts[m.cursor]
		
		// Show job information
		output = fmt.Sprintf("Name: %s\n", selectedContext.Name)
		output += fmt.Sprintf("Label: %s\n", selectedContext.Label)
		if selectedContext.Description != "" {
			output += fmt.Sprintf("Description: %s\n", selectedContext.Description)
		}
		
		if cmd, exists := selectedContext.Commands["run"]; exists {
			output += fmt.Sprintf("Command: %s\n", cmd)
		}
		
		if len(selectedContext.Variables) > 0 {
			output += "\nVariables:\n"
			for k, v := range selectedContext.Variables {
				output += fmt.Sprintf("  %s = %s\n", k, v)
			}
		}
		
		if selectedContext.LastResult != nil {
			output += "\nLast Execution:\n"
			output += fmt.Sprintf("  Time: %s\n", selectedContext.LastResult.Timestamp.Format("2006-01-02 15:04:05"))
			output += fmt.Sprintf("  Status: %s (Exit Code: %d)\n", 
				map[bool]string{true: "SUCCESS", false: "FAILED"}[selectedContext.LastResult.Success],
				selectedContext.LastResult.ExitCode)
			if selectedContext.LastResult.Output != "" {
				output += fmt.Sprintf("  Output:\n%s\n", selectedContext.LastResult.Output)
			}
		} else {
			output += "\nNever executed"
		}
	} else {
		output = "No job selected"
	}
	contentWidth := m.width - 4 - 4  // total width - borders - padding
	contentHeight := bottomHeight - 4  // title + separator + spacing + buffer
	
	if contentHeight < 1 {
		contentHeight = 1
	}
	
	// Split into lines and wrap long lines
	var processedLines []string
	for _, line := range strings.Split(output, "\n") {
		if len(line) <= contentWidth {
			processedLines = append(processedLines, line)
		} else {
			// Wrap long lines
			for i := 0; i < len(line); i += contentWidth {
				end := i + contentWidth
				if end > len(line) {
					end = len(line)
				}
				processedLines = append(processedLines, line[i:end])
			}
		}
	}
	
	// Limit to available height (show from beginning)
	if len(processedLines) > contentHeight {
		processedLines = processedLines[:contentHeight]
	}
	
	bottomContent.WriteString(strings.Join(processedLines, "\n"))

	topPanel := topPanelStyle.Render(topContent.String())
	bottomPanel := bottomPanelStyle.Render(bottomContent.String())

	return lipgloss.JoinVertical(lipgloss.Left, topPanel, bottomPanel)
}