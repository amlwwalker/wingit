package main

import (
	"github.com/jroimartin/gocui"
	controller "github.com/amlwwalker/wingit/packages/controller"
)

type UIControl struct {
	CONTROLLER *controller.CONTROLLER
	Views map[string]*gocui.View
	gui *gocui.Gui
    PendingUploadFileName     string //the file the user wants to upload (could be selected or dragged)
    PendingUploadFilePath         string //holds the filepath that is dragged over the drop area
    CurrentSelectedContact      string //their name
    CurrentSelectedContactIndex int //their index (if needed)
    // downloadFile        int
	LoggedIn bool
	Status string
}