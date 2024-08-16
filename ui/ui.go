package ui

import (
	"fmt"
	"time"

	"github.com/aayushbtw/monit/stats"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func Handler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)
	baseStyle := renderer.NewStyle()

	tbl := table.New(
		table.WithColumns([]table.Column{
			{Title: "PID", Width: 10},
			{Title: "Name", Width: 25},
			{Title: "CPU", Width: 12},
			{Title: "MEM", Width: 12},
			{Title: "Username", Width: 12},
			{Title: "Time", Width: 12},
		}),
		table.WithRows([]table.Row{}),
		table.WithHeight(10),
	)

	m := model{
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		baseStyle: baseStyle,
		tbl:       tbl,
	}

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type model struct {
	width  int
	height int

	tbl       table.Model
	baseStyle lipgloss.Style

	CpuUsage  cpu.TimesStat
	MemUsage  mem.VirtualMemoryStat
	SwapUsage mem.SwapMemoryStat
}

type TickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tickEvery()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case TickMsg:

		cpuStats, err := stats.GetCPUStats()
		if err != nil {
			log.Error("Could not get CPU info", "error", err)
		} else {
			m.CpuUsage = cpuStats
		}

		memStats, err := stats.GetMEMStats()
		if err != nil {
			log.Error("Could not get memory info", "error", err)
		} else {
			m.MemUsage = memStats
		}

		swapStats, err := stats.GetSWAPStats()
		if err != nil {
			log.Error("Could not get swap info", "error", err)
		} else {
			m.SwapUsage = swapStats
		}

		procs, err := stats.GetProcesses(10)
		if err != nil {
			log.Error("Could not get processes", "error", err)
		} else {
			rows := []table.Row{}
			for _, p := range procs {
				memString, memUnit := stats.ConvertBytes(p.Memory) // Use RSS for memory size
				rows = append(rows, table.Row{
					fmt.Sprintf("%d", p.PID),
					p.Name,
					fmt.Sprintf("%.2f%%", p.CPUPercent),
					fmt.Sprintf("%s %s", memString, memUnit), // Format memory value and unit
					p.Username,
					p.RunningTime,
				})
			}
			m.tbl.SetRows(rows)
		}

		//
		return m, tickEvery()
	}

	return m, nil
}

func (m model) View() string {
	content := m.baseStyle.
		Width(m.width).
		Height(m.height).
		Padding(4, 10).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				m.baseStyle.PaddingBottom(1).Render(m.ViewHeader()),
				m.tbl.View(),
			),
		)

	return content
}
