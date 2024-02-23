package theme2

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

var (
	//go:embed fonts/msyhl.ttc
	NotoSansSC []byte
)

type MyTheme struct {
	dark bool
}

type Option func(*MyTheme)

func ApplyDark(dark bool) Option {
	return func(myTheme *MyTheme) {
		myTheme.dark = dark
	}
}

func NewMyTheme(options ...Option) *MyTheme {
	m := new(MyTheme)
	for _, option := range options {
		option(m)
	}
	return m
}

var _ fyne.Theme = (*MyTheme)(nil)

// StaticName 为 fonts 目录下的 ttf 类型的字体文件名
func (m MyTheme) Font(fyne.TextStyle) fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "NotoSansSC.ttf",
		StaticContent: NotoSansSC,
	}
}

func (m *MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if m.dark {
		return theme.DarkTheme().Color(n, v)
	} else {
		return theme.LightTheme().Color(n, v)
	}

}

func (*MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
