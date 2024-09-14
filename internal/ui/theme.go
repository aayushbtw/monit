package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Primary   lipgloss.AdaptiveColor
	Secondary lipgloss.AdaptiveColor
	Highlight lipgloss.AdaptiveColor
	Border    lipgloss.AdaptiveColor
	Green     lipgloss.AdaptiveColor
	Red       lipgloss.AdaptiveColor
}

var Color = Theme{
	Primary:   lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"},
	Secondary: lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"},
	Highlight: lipgloss.AdaptiveColor{Light: "#8b2def", Dark: "#8b2def"},
	Border:    lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"},
	Green:     lipgloss.AdaptiveColor{Light: "#00FF00", Dark: "#00FF00"},
	Red:       lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"},
}
