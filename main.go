package main

import (
	"code-gen/pkg/initialize"
	"code-gen/pkg/tutorials"
	theme2 "code-gen/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

var topWindow fyne.Window

func main() {
	initialize.Init()
	a := app.NewWithID("code-gen")
	w := a.NewWindow("一个平平无奇的代码生成器")
	topWindow = w
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	content := container.NewMax()
	setTutorial := func(t tutorials.Tutorial) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}
	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setTutorial, false))
	} else {
		split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

const preferenceCurrentTutorial = "currentTutorial"

func makeNav(setTutorial func(tutorial tutorials.Tutorial), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()
	a.Settings().SetTheme(theme2.NewMyTheme(theme2.ApplyDark(true)))
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return tutorials.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := tutorials.TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := tutorials.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := tutorials.Tutorials[uid]; ok {
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}
	tree.OpenAllBranches()
	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "MySQL2")
		tree.Select(currentPref)
	}
	return container.NewBorder(nil, nil, nil, nil, tree)
}
func init() {

	fontPaths := findfont.List()
	for _, path := range fontPaths {
		//fmt.Println(path)
		//楷体:simkai.ttf
		//黑体:simhei.ttf
		//微软雅黑:msyh.ttc
		if strings.Contains(path, "msyh.ttc") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}
