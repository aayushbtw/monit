package ui

import (
	"fmt"
	"time"

	"github.com/aayushbtw/monit/internal/stats"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

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
		case "esc":
			if m.processTable.Focused() {
				m.tableStyle.Selected = m.baseStyle
				m.processTable.SetStyles(m.tableStyle)
				m.processTable.Blur()
			} else {
				m.tableStyle.Selected = m.tableStyle.Selected.Background(Color.Highlight)
				m.processTable.SetStyles(m.tableStyle)
				m.processTable.Focus()
			}
		case "up", "k":
			if m.processTable.Focused() {
				m.processTable.MoveUp(1)
			}
		case "down", "j":
			if m.processTable.Focused() {
				m.processTable.MoveDown(1)
			}
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

		currentNetStats, err := stats.GetNetworkStats()
		if err != nil {
			log.Error("Could not get network info", "error", err)
		} else {
			if m.PrevNetworkStats.BytesSent != 0 || m.PrevNetworkStats.BytesRecv != 0 {
				uploadSpeed := float64(currentNetStats.BytesSent-m.PrevNetworkStats.BytesSent) / 1.0
				downloadSpeed := float64(currentNetStats.BytesRecv-m.PrevNetworkStats.BytesRecv) / 1.0

				m.NetworkUploadSpeed = uploadSpeed
				m.NetworkDownloadSpeed = downloadSpeed
			}
			m.PrevNetworkStats = currentNetStats
			m.NetworkStats = currentNetStats
		}

		procs, err := stats.GetProcesses(40)
		if err != nil {
			log.Error("Could not get processes", "error", err)
		} else {
			rows := []table.Row{}
			for _, p := range procs {
				memString, memUnit := stats.ConvertBytes(p.Memory)
				rows = append(rows, table.Row{
					fmt.Sprintf("%d", p.PID),
					p.Name,
					fmt.Sprintf("%.2f%%", p.CPUPercent),
					fmt.Sprintf("%s %s", memString, memUnit),
					p.Username,
					p.RunningTime,
				})
			}
			m.processTable.SetRows(rows)
		}

		//
		return m, tickEvery()
	}

	return m, nil
}
