package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Tray struct {
	obj      *QSystemTrayIconWithCustomSlot
	rootmenu *widgets.QMenu
	first    *widgets.QAction
	quit     *widgets.QAction
}

type QSystemTrayIconWithCustomSlot struct {
	widgets.QSystemTrayIcon
	_ func() `slot:"triggerSlot"`

	_ func()                                     `constructor:"init"`
	_ func(text string, action *widgets.QAction) `signal:"setTextForAction"`
}

func (t *QSystemTrayIconWithCustomSlot) init() {
	t.ConnectSetTextForAction(func(text string, action *widgets.QAction) { action.SetText(text) })

}

//icons must be in qml/pictures dir
func NewTray(iconName string) *Tray {

	var systray = NewQSystemTrayIconWithCustomSlot(nil)
	var systrayMenu = widgets.NewQMenu(nil)
	var icon *gui.QIcon
	icon = gui.NewQIcon()
	icon.AddPixmap(gui.NewQPixmap5(":/qml/pictures/"+iconName, "", core.Qt__AutoColor), gui.QIcon__Normal, gui.QIcon__Off)
	systray.SetIcon(icon)
	systray.SetContextMenu(systrayMenu)
	t := &Tray{
		obj:      systray,
		rootmenu: systrayMenu,
	}

	return t
}

func (t *Tray) build(showUI func(bool)) {

	t.addAction("Show UI", showUI)
	// t.addAction("more-padding", func(bool) {})
	// t.first = t.addAction("updated", func(bool) {
	// 	log.Println("frist clicked")
	// })
	t.quit = t.addAction("quit", func(bool) {
		log.Println("exit")
		os.Exit(0)
	})

	t.obj.Show()
}

func (t *Tray) addAction(str string, fn func(bool)) *widgets.QAction {
	a := widgets.NewQAction2(str, nil)
	a.ConnectTriggered(fn)
	t.rootmenu.AddActions([]*widgets.QAction{a})
	return a
}
func (t *Tray) update() {
	count := 0
	time.Sleep(2 * time.Second)
	for {
		time.Sleep(5 * time.Second)
		t.obj.SetTextForAction(fmt.Sprintf("updated %d", count), t.first)
		log.Println("Set", t.first.Text())
		count++
	}
}
