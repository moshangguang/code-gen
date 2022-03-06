package Theme

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Theme int8

const (
	Light = 1
	Dark  = 2
)

func (t Theme) ToTheme() fyne.Theme {
	switch t {
	case Light:
		return theme.LightTheme()
	case Dark:
		return theme.DarkTheme()
	default:
		return theme.LightTheme()
	}
}
