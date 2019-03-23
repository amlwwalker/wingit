package main

import (
	"encoding/json"
	"fmt"

	"os"
	"path/filepath"
	"strings"

	utils "github.com/amlwwalker/wingit/packages/utils"
	"github.com/gobuffalo/packr"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/widgets"
)

func LoadConfiguration() (error, utils.Config) {
	var config utils.Config

	//lets compile the config.json file into the binary so its easily accessible
	box := packr.NewBox("./configfiles")
	if configFile, err := box.MustBytes("config.json"); err != nil {
		return err, config
	} else {
		json.Unmarshal(configFile, &config)
		return nil, config
	}
}
func main() {
	//0. set any required env vars for qt
	os.Setenv("QT_QUICK_CONTROLS_STYLE", "material") //set style to material
	//1. the hotloader needs a path to the qml files highest directory
	// change this if you are working elsewhere
	var topLevel = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "amlwwalker", "wingit", "_guiinterface", "qml")

	_, config := LoadConfiguration()
	//2. load the configuration file
	fmt.Println("real config: ", config)
	//3. Create a bridge to the frontend
	var qmlBridge = NewQmlBridge(nil)
	//hi def scaling
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	//4. Configure the qml binding and create an application
	widgets.NewQApplication(len(os.Args), os.Args)

	//5. configure the bridge and the systray
	qmlBridge.ConfigureBridge(config)
	qmlBridge.business.CONTROLLER.StoreStatusMessage("config loaded", 1)

	//create a view
	var view = quick.NewQQuickView(nil)
	view.SetTitle("WingIt!")
	//enable notifiers now, so that they can be used from front end if needs be, later
	qmlBridge.business.notifier.Initialise()
	//configure the view to know about the bridge
	//this needs to happen before anything happens on another thread
	//else the thread might beat the context property to setup
	view.RootContext().SetContextProperty("QmlBridge", qmlBridge)
	view.RootContext().SetContextProperty("QmlUser", qmlBridge.User)
	view.RootContext().SetContextProperty("ContactsModel", qmlBridge.business.pModel)
	view.RootContext().SetContextProperty("SearchModel", qmlBridge.business.sModel)
	view.RootContext().SetContextProperty("FilesModel", qmlBridge.business.fModel)
	view.RootContext().SetContextProperty("DownloadsModel", qmlBridge.business.dModel)

	qmlBridge.business.CONTROLLER.StoreStatusMessage("contexts in place", 1)
	//5. Configure hotloading
	//configure the loader to handle updating qml live
	loader := func(p string) {
		fmt.Println("changed:", p)
		view.SetSource(core.NewQUrl())
		view.Engine().ClearComponentCache()
		view.SetSource(core.NewQUrl3(topLevel+"/loader.qml", 0))
		if !strings.Contains(p, "/loader.qml") {
			relativePath := strings.Replace(p, topLevel+"/", "", -1)
			qmlBridge.UpdateLoader(relativePath)
		}
	}
	//decide whether to enable hotloading (must be disabled for deployment)
	if !config.Hotload {
		fmt.Println("compiling qml into binary...")
		view.SetSource(core.NewQUrl3("qrc:/qml/loader-production.qml", 0))
	} else {
		view.SetSource(core.NewQUrl3(topLevel+"/loader-production.qml", 0))
		go qmlBridge.hotLoader.startWatcher(loader)
	}
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
	qmlBridge.business.CONTROLLER.StoreStatusMessage("resizing and show off components", 1)

	// view.SetFlags(core.Qt__WindowTitleHint | core.Qt__CustomizeWindowHint)
	view.SetFlags(core.Qt__WindowTitleHint | core.Qt__CustomizeWindowHint | core.Qt__WindowMinimizeButtonHint)
	// view.SetFlags(core.Qt__MSWindowsFixedSizeDialogHint)
	//6. Complete setup, and start the UI
	view.SetMaximumWidth(view.Width())
	// view.SetFixedSize(view.width(), view.height())
	qmlBridge.business.systray.build(func(bool) {
		fmt.Println("show view from systray")
		//could do a check here, but at the moment just do both
		//bring to front
		view.Raise()
		//show if minimized
		view.Show()
	})
	view.Show()

	widgets.QApplication_Exec()

}
