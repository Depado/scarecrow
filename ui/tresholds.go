package ui

import "github.com/gdamore/tcell/v2"

var (
	DefaultThreshold = tcell.StyleDefault.Foreground(tcell.ColorWhite)
	GreenThreshold   = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)
	YelloThreshold   = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorYellowGreen)
	OrangeThreshold  = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorOrange)
	RedThreshold     = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
)

func GetIntStyle(perc int) tcell.Style {
	switch {
	case perc < 40:
		return GreenThreshold
	case perc < 60:
		return YelloThreshold
	case perc < 80:
		return OrangeThreshold
	case perc >= 80:
		return RedThreshold
	}
	return DefaultThreshold
}

func GetFloatStyle(perc float64) tcell.Style {
	switch {
	case perc < 40:
		return GreenThreshold
	case perc < 60:
		return YelloThreshold
	case perc < 80:
		return OrangeThreshold
	case perc >= 80:
		return RedThreshold
	}
	return DefaultThreshold
}
