package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

func Handler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)

	tableStyle := table.DefaultStyles()
	tableStyle.Selected = renderer.NewStyle().Background(Color.Highlight)
	processTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "PID", Width: 10},
			{Title: "Name", Width: 25},
			{Title: "CPU", Width: 12},
			{Title: "MEM", Width: 12},
			{Title: "Username", Width: 12},
			{Title: "Time", Width: 12},
		}),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(20),
		table.WithWidth(pty.Window.Width),
		table.WithStyles(tableStyle),
	)

	m := model{
		width:        pty.Window.Width,
		height:       pty.Window.Height,
		processTable: processTable,
		tableStyle:   tableStyle,
		baseStyle:    renderer.NewStyle(),
		viewStyle: renderer.NewStyle().
			// Background(lipgloss.Color("#33333")). // For Debug
			Width(pty.Window.Width),
	}

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
