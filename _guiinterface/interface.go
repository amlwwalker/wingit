package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/skratchdot/open-golang/open"

	controller "github.com/amlwwalker/wingit/packages/controller"
	utils "github.com/amlwwalker/wingit/packages/utils"
	"github.com/therecipe/qt/core"
)

type QmlUser struct {
	core.QObject

	_ string `property:"name"`
	_ string `property:"email"`
	_ string `property:"picture"`
	_ string `property:"apiKey"`
}

/*
	type User struct {
		Id      string `"json":"user_id"`
		Email   string `"json":"email"`
		Name    string `"json":"name"`
		Picture string `"json":"picture"`
		Locale  string `"json":"locale"`
		ApiKey  string `"json":"-"`
	}

*/
type QmlBridge struct {
	core.QObject
	hotLoader     HotLoader
	business      BusinessInterface
	RefreshTicker *time.Ticker
	TickerEnabled bool     `property:"tickerEnabled"`
	User          *QmlUser `property:"qmlUser"`
	//qml properties:
	// UserName string `property:"username"`
	// UsersEmail string `property:"usersEmail"`
	// UserPic string `property:"userPic"`
	_ func(p string) `signal:"updateLoader"`

	//updates the progress bars on the front end
	_ func(c float64, indeterminate bool) `signal:"updateProcessStatus"`
	_ func(message string)                `signal:"popUpToast"`
	_ func(message string)                `signal:"setMessage"`

	_ func(username, error string) `signal:"signalLogin"`
	_ func()                       `signal:"signalLogout"`

	//requests from qml
	//searchFor processes the searching request from the front end
	_ func(regex string)                   `slot:"searchFor"`
	_ func()                               `slot:"checkForUser"`
	_ func() bool                          `slot:"toggleAutoSync"`
	_ func() bool                          `slot:"toggleState"`
	_ func(title, message string)          `slot:"createNotification"`
	_ func(pendingUser, fileUrl string)    `slot:"updatePendingUpload"`
	_ func(searchResult string, index int) `slot:"addResultToContacts"`
	_ func(userId string)                  `slot:"updatePendingUser"`
	_ func(fileName string)                `slot:"beginFileDownload"`
	_ func()                               `slot:"getDownloads"`
	_ func(path string)                    `slot:"openFile"`
	_ func()                               `slot:"openDowloadsDirectory"`
	_ func()                               `slot:"syncFiles"`
	_ func(apiKey string)                  `slot:"loginUser"`
	_ func()                               `slot:"logout"`
}

//setup functions to communicate between front end and back end

//example of receiving data from frontend and returning a result
func (q *QmlBridge) ConfigureBridge(config utils.Config) {
	//1. configure the hotloader
	q.business = BusinessInterface{}
	q.business.configureInterface(q.SignalLogin, q.UpdateProcessStatus, config)
	q.hotLoader = HotLoader{} //may not need it, specified in main.go
	q.business.CONTROLLER.StoreStatusMessage("business configured", 1)
	q.RefreshTicker = time.NewTicker(60 * time.Second)
	q.User = NewQmlUser(nil) //create a person to add to the contacts
	//periodical functionality
	go func() {
		for {
			<-q.RefreshTicker.C
			fmt.Println("auto sync enabled ? ", q.TickerEnabled)
			if q.TickerEnabled {
				q.business.SynchronizeWithServer(q.SetMessage)
			}
		}
	}()

	q.ConnectToggleAutoSync(func() bool {
		q.TickerEnabled = !q.TickerEnabled
		return q.TickerEnabled
	})
	q.ConnectToggleState(func() bool {
		return q.TickerEnabled
	})
	//2. Configure signals
	q.ConnectSearchFor(func(regex string) {
		//in here we are going to add matches to the search model
		//that way the front end will be updated live
		//inform front end work has begun
		fmt.Println("requested search for: " + regex)
		q.UpdateProcessStatus(0.0, true) //required as no idea where it is so needs a start point
		q.business.searchForMatches(regex, q.UpdateProcessStatus)
	})
	q.ConnectAddResultToContacts(func(searchResult string, index int) {
		//in here we are going to add matches to the search model
		//that way the front end will be updated live
		//inform front end work has begun
		fmt.Println("adding to contacts: " + searchResult)
		p := q.business.CONTROLLER.CreatePersonFromIdentifier(searchResult)
		q.business.CONTROLLER.AddContactToList(p)
		fmt.Println("deleting key for "+q.business.CONTROLLER.CRYPTO.KeyFolder, p.Name+".pem")
		q.business.CONTROLLER.CRYPTO.DeletePublickKeyForContact(q.business.CONTROLLER.CRYPTO.KeyFolder, p.Name+".pem")
		//when we add a search result we want to remove them from the search results
		//and add them to our cotacts
		//we need to clear the list and redisplay it so we don't get repeats
		q.business.sModel.removePerson(index)
		q.business.pModel.ClearPeople()
		for _, v := range q.business.CONTROLLER.Contacts.People {
			// fmt.Println("after updating files, adding: ", v.UserId)
			addPersonToList(*v, q.business.pModel)
		}

		// q.business.pModel.AddPerson(qP)
		// q.business.sModel.removePerson(index)
	})
	q.ConnectSyncFiles(func() {
		q.TickerEnabled = !q.TickerEnabled
		//manuall synchronising with the server will toggle auto sync
		q.business.SynchronizeWithServer(q.SetMessage)
		q.SetMessage("SyncK'd!")
	})
	q.ConnectUpdatePendingUser(func(pendingUser string) {
		q.business.CurrentSelectedContact = pendingUser
		//if the files have been synced already, we just need to update
		//the file list to the files owned by the currently selected contact
		//need to clear the list first
		q.business.fModel.ClearFiles()
		fmt.Println("pending user: " + pendingUser)
		fmt.Println("files: ", q.business.CONTROLLER.Contacts.People[pendingUser])
		for _, v := range q.business.CONTROLLER.Contacts.People[pendingUser].Files {
			//now we add each one to the file list
			var f = NewFile(nil)
			f.SetFilePath(v.FileName)
			f.SetFileSize(strconv.Itoa(v.FileSize))
			f.SetFileSource(v.UserID)
			q.business.fModel.AddFile(f)
		}
	})
	q.ConnectUpdatePendingUpload(func(pendingUser, fileUrl string) {
		q.business.CurrentSelectedContact = pendingUser
		q.business.PendingUpload = fileUrl
		fmt.Println("received: ", q.business.PendingUpload)
		fmt.Println("for: ", q.business.CurrentSelectedContact)
		//however now we just need to upload the file to the server
		fmt.Println("beginning upload")
		q.business.CONTROLLER.SERVER.UploadProgress(0.0)
		q.business.CONTROLLER.UploadFileToContact(pendingUser, fileUrl)
	})
	q.ConnectBeginFileDownload(func(fileName string) {
		fmt.Println("beginning file download")
		q.PopUpToast("File Download Beginning")
		q.UpdateProcessStatus(0.0, false) //to show the progress bar
		q.business.CONTROLLER.DownloadFileFromContact(q.business.CurrentSelectedContact, fileName)
	})
	q.ConnectOpenFile(func(path string) {
		q.business.CONTROLLER.OpenDownloadedFile(path)
	})
	q.ConnectOpenDowloadsDirectory(func() {
		fmt.Println("opening downloads directory at " + q.business.CONTROLLER.SERVER.DownloadFolder)
		open.Run(q.business.CONTROLLER.SERVER.DownloadFolder)
	})

	q.ConnectGetDownloads(func() {
		//this shows the user all the downloads they have
		if files, err := q.business.CONTROLLER.GetDownloadedFiles(); err != nil {
			//couldn't retrieve errors
			fmt.Println("error retrieving previously downloaded files", err)
		} else {
			q.business.dModel.ClearFiles()
			for _, v := range files {
				var f = NewFile(nil)
				f.SetFilePath(v.FileName)
				f.SetFileSize(strconv.Itoa(v.FileSize))
				q.business.dModel.AddFile(f)
			}
		}
	})
	q.ConnectCheckForUser(func() {
		//here we can test if the database has a user in it.
		//if it does, signal the user to the front end
		//if it finds a result, it will call UserAuthenticated
		q.business.CONTROLLER.StoreStatusMessage("checking if user exists: ", 1)
		if err := q.business.CONTROLLER.CheckForExistingUser(); err != nil {
			//there was no user, return
			return
		}

		fmt.Println("Synchronising with server on auto login")
		q.business.SynchronizeWithServer(q.SetMessage)
		// q.QmlUser = q.business.CONTROLLER.User
		// q.QmlUser

		/*
			type User struct {
				Id      string `"json":"user_id"`
				Email   string `"json":"email"`
				Name    string `"json":"name"`
				Picture string `"json":"picture"`
				Locale  string `"json":"locale"`
				ApiKey  string `"json":"-"`
			}

		*/
		q.User.SetName(q.business.CONTROLLER.User.Name)
		q.User.SetEmail(q.business.CONTROLLER.User.Email)
		q.User.SetPicture(q.business.CONTROLLER.User.Picture)
		q.User.SetApiKey(q.business.CONTROLLER.User.ApiKey)

		q.TickerEnabled = true
	})
	q.ConnectLoginUser(func(apiKey string) {
		//this calls the server to open the login
		//but once the server has authd the user
		//an asynchronous response will tell the user that
		//login was successful

		//this function results in a callback running
		//so that this signal doesn't end up dealing with all post auth
		//logic
		q.business.notifier.Push("Authentication", "Attempting Login")
		fmt.Println("api key received: ", apiKey)
		if u, err := q.business.CONTROLLER.Authorize(apiKey); err != nil {
			fmt.Println("Auth failed; ", err)
		} else {
			fmt.Println("user: ", u)
			//at this point, we can create the users contact bucket
			if err := controller.CreateUserContactBucket(u.Id, q.business.CONTROLLER.Db); err != nil {
				fmt.Println("could not create a contact bucket for the user: ", err)
			}
			fmt.Println("Synchronising with server now logged in")
			q.business.SynchronizeWithServer(q.SetMessage)
			q.TickerEnabled = true
		}
	})
	q.ConnectCreateNotification(func(title, message string) {
		q.business.notifier.Push(title, message)
	})

	//we might need to clear the contacts and files
	q.ConnectLogout(func() {
		//this calls the server to open the login
		//but once the server has authd the user
		//an asynchronous response will tell the user that
		//login was successful
		q.business.CONTROLLER.Logout(func() {
			fmt.Println("logout called back")
			//in reality this calls the front end, which
			//updates the stack view to the home page
			//blank the user
			fmt.Println("cleared user: ", q.business.CONTROLLER.User)
			//need to clear current user settings out of the DB
			// q.UserName = "" //clear for frontend
			//be good to get user pics and store them in db
			// q.UserPic = ""
			//clear all front end lists
			q.business.pModel.ClearPeople()
			q.business.sModel.ClearPeople()
			q.business.fModel.ClearFiles()
			q.business.dModel.ClearFiles()
			//tell the UI logout has occured
			q.SignalLogout()
		})
	})
	q.business.CONTROLLER.StoreStatusMessage("bridge configured", 1)
}

// helper function to add people to one of the interface models
func addPersonToList(tmp controller.Person, model *PersonModel) {
	fmt.Println("about to add: ", tmp.UserId, " to list with ", tmp.Len, " files")
	var p = NewPerson(nil)
	p.SetFirstName(tmp.Name)
	p.SetLastName("")
	p.SetEmail(tmp.UserId)
	p.SetFileCount(strconv.Itoa(tmp.Len))
	fmt.Println("adding person to front end: " + tmp.UserId)
	model.AddPerson(p)
}
