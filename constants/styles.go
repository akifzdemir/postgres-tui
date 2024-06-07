package constants

import "github.com/charmbracelet/lipgloss"

const (
	HotPink    = lipgloss.Color("#FF06B7")
	DarkGray   = lipgloss.Color("#767676")
	Red        = lipgloss.Color("160")
	BlueViolet = lipgloss.Color("57")
	White      = lipgloss.Color("15")
	Yellow     = lipgloss.Color("226")
	Black      = lipgloss.Color("0")
)

var (
	InputStyle  = lipgloss.NewStyle().Background(BlueViolet).Foreground(White).Width(20)
	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(BlueViolet).
			Padding(2, 4)
)
