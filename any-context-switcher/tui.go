package main

import (
	"fmt"
	"strings"

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
}

func getStyles(theme ColorTheme) (titleStyle, selectedStyle, topPanelStyle, bottomPanelStyle, outputTitleStyle lipgloss.Style) {
	titleStyle = lipgloss.NewStyle().
		MarginLeft(2).
		Foreground(lipgloss.Color(theme.Title))

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Selected)).
		Bold(true)

	topPanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Border)).
		Padding(1, 2).
		Height(15)

	bottomPanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Border)).
		Padding(1, 2).
		Height(10)

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
				
				if context.Status == "active" {
					context.Status = "inactive"
					context.LastError = false
					m.lastOutput = "Deactivated: " + context.Label
				} else {
					if activateCmd, exists := context.Commands["activate"]; exists {
						output, err := m.executor.executeCommandWithOutput(activateCmd, context.Variables)
						if err != nil {
							context.LastError = true
							m.lastOutput = fmt.Sprintf("Error executing command: %v\n\nOutput:\n%s", err, output)
						} else {
							context.Status = "active"
							context.LastError = false
							m.lastOutput = fmt.Sprintf("Activated: %s\n\nCommand output:\n%s", context.Label, output)
						}
					} else {
						context.Status = "active"
						context.LastError = false
						m.lastOutput = "Activated: " + context.Label
					}
				}
				
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
	return m, nil
}

func (m *model) View() string {
	var topContent strings.Builder
	var bottomContent strings.Builder

	titleStyle, selectedStyle, topPanelStyle, bottomPanelStyle, outputTitleStyle := getStyles(m.theme)

	topContent.WriteString(titleStyle.Render("Context Switcher"))
	topContent.WriteString("\n\n")

	if len(m.contexts) == 0 {
		topContent.WriteString("No contexts available.")
	} else {
		for i, context := range m.contexts {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			checkbox := "[ ]"
			if context.Status == "active" {
				checkbox = "[x]"
			}

			errorIcon := ""
			if context.LastError {
				errorIcon = " ✗"
			}

			style := lipgloss.NewStyle()
			if m.cursor == i {
				style = selectedStyle
			}

			line := fmt.Sprintf("%s %s%s %s", 
				cursor, checkbox, errorIcon, context.Label)
			
			if context.Description != "" {
				line += fmt.Sprintf(" - %s", context.Description)
			}

			topContent.WriteString(style.Render(line))
			topContent.WriteString("\n")
		}
	}

	topContent.WriteString("\n↑/↓ or j/k: navigate • space: toggle • q: quit")

	bottomContent.WriteString(outputTitleStyle.Render("Command Output"))
	bottomContent.WriteString("\n")
	bottomContent.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	bottomContent.WriteString(m.lastOutput)

	topPanel := topPanelStyle.Render(topContent.String())
	bottomPanel := bottomPanelStyle.Render(bottomContent.String())

	return lipgloss.JoinVertical(lipgloss.Left, topPanel, bottomPanel)
}