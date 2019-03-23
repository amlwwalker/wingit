/*

CUI

*/

package main

import (
	"encoding/json"
	"fmt"

	controller "github.com/amlwwalker/wingit/packages/controller"
	cryptography "github.com/amlwwalker/wingit/packages/cryptography"
	srv "github.com/amlwwalker/wingit/packages/server"
	utils "github.com/amlwwalker/wingit/packages/utils"
	"github.com/atrox/homedir"
	"github.com/gobuffalo/packr"
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

	_, config := LoadConfiguration()

	VERBOSE := false

	PATH, err := homedir.Dir()
	if err != nil {
		fmt.Println("couldnt get home directory ", err)
		panic(err)
	}
	PATH = PATH + "/.wingit/"
	fmt.Println("path is: " + PATH)
	KEYFOLDER := PATH + ".keys/"
	DOWNLOADFOLDER := PATH + ".downloads/"
	SYNCFOLDER := PATH + ".sync/"

	// Initialise the controller. This will handle the user, server, crypto and ui.
	CONTROLLER := &controller.CONTROLLER{}
	CONTROLLER.DBPath = PATH //appends the db name once logged in (db name is per account)

	var server srv.Server
	//set up the server
	server.Ip = config.Host
	server.Address = "https://" + server.Ip
	server.Port = "8090" // Default port.
	server.Verbose = VERBOSE
	var crypto cryptography.Crypto

	CONTROLLER.SERVER = &server
	CONTROLLER.CRYPTO = &crypto
	PASSWORDLENGTH := 32 // in bytes

	//make sure the directories exist
	utils.InitiateDirectory(KEYFOLDER)
	utils.InitiateDirectory(SYNCFOLDER)
	utils.InitiateDirectory(DOWNLOADFOLDER)

	CONTROLLER.CRYPTO.Init(PASSWORDLENGTH, KEYFOLDER, SYNCFOLDER, VERBOSE)
	CONTROLLER.SERVER.Init(DOWNLOADFOLDER, SYNCFOLDER)

	// --------------------
	// SERVER
	utils.PrintLine()
	user := &controller.User{}
	CONTROLLER.User = user
	CONTROLLER.Contacts = &controller.Contacts{}
	CONTROLLER.Contacts.People = make(map[string]*controller.Person)
	CONTROLLER.SearchResults = &controller.SearchResults{}

	cui := &UIControl{}
	cui.CONTROLLER = CONTROLLER
	/*
	   once the cui knows about the controller
	   agnostic function calls the controller needs to be able to access
	   can be updated.
	*/
	CONTROLLER.Logger = cui.ExternalLogUpdate
	CONTROLLER.SERVER.Logger = CONTROLLER.Logger
	CONTROLLER.CRYPTO.Logger = CONTROLLER.Logger
	cui.CONTROLLER.UserAuthenticated = cui.UserAuthenticatedCallback

	//default the application to the user being logged out
	cui.LoggedIn = false
	//set an initial status message for the UI
	cui.Status = "welcome to Wingit"

	cui.Init() //init the UI

} // end of main
