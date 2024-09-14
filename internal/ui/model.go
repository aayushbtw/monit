package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

type model struct {
	width  int
	height int

	processTable table.Model
	tableStyle   table.Styles
	baseStyle    lipgloss.Style
	viewStyle    lipgloss.Style

	CpuUsage             cpu.TimesStat
	MemUsage             mem.VirtualMemoryStat
	SwapUsage            mem.SwapMemoryStat
	NetworkStats         net.IOCountersStat
	PrevNetworkStats     net.IOCountersStat
	NetworkUploadSpeed   float64
	NetworkDownloadSpeed float64
}

type TickMsg time.Time
