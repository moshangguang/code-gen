package tutorials

import (
	"code-gen/data"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"net/url"
)

func logoTextScreen(title string) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneScene)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(256, 256))
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		container.NewHBox(
			widget.NewLabel("\t\t"),
			widget.NewHyperlink("GitHub", parseURL("https://github.com/moshangguang")),
			widget.NewLabel("-"),
			widget.NewHyperlink("博客园", parseURL("https://www.cnblogs.com/beiluowuzheng")),
		),
		//widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))
}
func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}
