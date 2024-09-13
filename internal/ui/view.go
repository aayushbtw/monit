package ui

import (
	"fmt"
	"strings"

	"github.com/aayushbtw/monit/internal/stats"
	"github.com/charmbracelet/lipgloss"
)

func ProgressBar(percentage float64, baseStyle lipgloss.Style) string {
	totalBars := 20
	fillBars := int(percentage / 100 * float64(totalBars))
	filled := baseStyle.
		Foreground(Color.Green).
		Render(strings.Repeat("|", fillBars))
	empty := baseStyle.
		Foreground(Color.Secondary).
		Render(strings.Repeat("|", totalBars-fillBars))

	return baseStyle.Render(fmt.Sprintf("%s%s%s%s", "[", filled, empty, "]"))
}

func (m model) ViewBanner() string {
	return m.baseStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render("Monit")
}

func (m model) ViewHeader() string {
	list := m.baseStyle.
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(Color.Border).
		Height(4).
		Padding(0, 1)

	listHeader := m.baseStyle.Bold(true).Render

	listItem := func(key string, value string, suffix ...string) string {
		finalSuffix := ""
		if len(suffix) > 0 {
			finalSuffix = suffix[0]
		}

		listItemValue := m.baseStyle.Align(lipgloss.Right).Render(fmt.Sprintf("%s%s", value, finalSuffix))

		listItemKey := func(key string) string {
			return m.baseStyle.Render(key + ":")
		}

		return fmt.Sprintf("%s %s", listItemKey(key), listItemValue)
	}

	return m.viewStyle.Render(
		lipgloss.JoinHorizontal(lipgloss.Top,
			// Progress Bars
			list.Render(
				lipgloss.JoinVertical(lipgloss.Left,
					listHeader("Monit."),
					listItem("CPU", fmt.Sprintf("%s %.1f", ProgressBar(100-m.CpuUsage.Idle, m.baseStyle), 100-m.CpuUsage.Idle), "%"),
					listItem("MEM", fmt.Sprintf("%s %.1f", ProgressBar(m.MemUsage.UsedPercent, m.baseStyle), m.MemUsage.UsedPercent), "%"),
					listItem("SWAP", fmt.Sprintf("%s %.1f", ProgressBar(m.SwapUsage.UsedPercent, m.baseStyle), m.SwapUsage.UsedPercent), "%"),
				),
			),
			//

			// CPU
			list.Border(lipgloss.NormalBorder(), false).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					listHeader("CPU"),
					listItem("user", fmt.Sprintf("%.1f", m.CpuUsage.User), "%"),
					listItem("sys", fmt.Sprintf("%.1f", m.CpuUsage.System), "%"),
					listItem("idle", fmt.Sprintf("%.1f", m.CpuUsage.Idle), "%"),
				),
			),
			list.Border(lipgloss.NormalBorder(), false).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					listHeader(""),
					listItem("nice", fmt.Sprintf("%.1f", m.CpuUsage.Nice), "%"),
					listItem("iowait", fmt.Sprintf("%.1f", m.CpuUsage.Iowait), "%"),
					listItem("irq", fmt.Sprintf("%.1f", m.CpuUsage.Irq), "%"),
				),
			),
			list.Render(
				lipgloss.JoinVertical(lipgloss.Left,
					listHeader(""),
					listItem("softirq", fmt.Sprintf("%.1f", m.CpuUsage.Softirq), "%"),
					listItem("steal", fmt.Sprintf("%.1f", m.CpuUsage.Steal), "%"),
					listItem("guest", fmt.Sprintf("%.1f", m.CpuUsage.Guest), "%"),
				),
			),
			//

			// MEM
			list.Border(lipgloss.NormalBorder(), false).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					listHeader("MEM"),
					func() string {
						value, unit := stats.ConvertBytes(m.MemUsage.Total)
						return listItem("total", value, unit)
					}(),
					func() string {
						value, unit := stats.ConvertBytes(m.MemUsage.Used)
						return listItem("used", value, unit)
					}(),
					func() string {
						value, unit := stats.ConvertBytes(m.MemUsage.Available)
						return listItem("free", value, unit)
					}(),
				),
			),
			list.Render(
				lipgloss.JoinVertical(lipgloss.Left,
					listHeader(""),
					func() string {
						value, unit := stats.ConvertBytes(m.MemUsage.Active)
						return listItem("active", value, unit)
					}(),
					func() string {
						value, unit := stats.ConvertBytes(m.MemUsage.Buffers)
						return listItem("buffers", value, unit)
					}(),
					func() string {
						value, unit := stats.ConvertBytes(m.MemUsage.Cached)
						return listItem("cached", value, unit)
					}(),
				),
			),
			//

			// SWAP
			list.Render(
				lipgloss.JoinVertical(lipgloss.Left,
					listHeader("SWAP"),
					func() string {
						value, unit := stats.ConvertBytes(m.SwapUsage.Total)
						return listItem("total", value, unit)
					}(),
					func() string {
						value, unit := stats.ConvertBytes(m.SwapUsage.Used)
						return listItem("used", value, unit)
					}(),
					func() string {
						value, unit := stats.ConvertBytes(m.SwapUsage.Free)
						return listItem("free", value, unit)
					}(),
				),
			),
			//
			// Network
			list.Border(lipgloss.NormalBorder(), false).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					listHeader("Network"),
					func() string {
						value, unit := stats.ConvertBytes(uint64(m.NetworkUploadSpeed))
						return listItem("up", value, unit)
					}(),
					func() string {
						value, unit := stats.ConvertBytes(uint64(m.NetworkDownloadSpeed))
						return listItem("down", value, unit)
					}(),
				),
			),
		),
	)
}

func (m model) ViewProcess() string {
	return m.viewStyle.Render(m.processTable.View())
}

func (m model) View() string {
	column := m.baseStyle.Width(m.width).Padding(1, 0, 0, 0).Render
	content := m.baseStyle.
		Width(m.width).
		Height(m.height).
		// Padding(4, 10).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				column(m.ViewHeader()),
				column(m.ViewProcess()),
			),
		)

	return content
}
