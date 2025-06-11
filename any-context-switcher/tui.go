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
}

var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(lipgloss.Color("86"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	currentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true).
			Background(lipgloss.Color("57"))
)

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
		case "enter", " ":
			if len(m.contexts) > 0 {
				contextName := m.contexts[m.cursor].Name
				if err := m.executor.switchContext(contextName); err != nil {
					return m, tea.Quit
				}
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Any Context Switcher"))
	b.WriteString("\n\n")

	if len(m.contexts) == 0 {
		b.WriteString("No contexts available.\n")
		b.WriteString("\nPress q to quit.")
		return b.String()
	}

	current := m.executor.getCurrentContext()
	
	for i, context := range m.contexts {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		style := lipgloss.NewStyle()
		if current != nil && current.Name == context.Name {
			style = currentStyle
		} else if m.cursor == i {
			style = selectedStyle
		}

		line := fmt.Sprintf("%s %s (%s) - %s", 
			cursor, context.Label, context.Name, context.Status)
		
		if context.Description != "" {
			line += fmt.Sprintf(" | %s", context.Description)
		}

		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString("↑/↓ or j/k: navigate • enter/space: switch • q: quit")

	return b.String()
}