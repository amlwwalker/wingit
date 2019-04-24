package main

import (
	"fmt"
	"strconv"
	"time"

	controller "github.com/amlwwalker/wingit/packages/controller"
	cryptography "github.com/amlwwalker/wingit/packages/cryptography"
	srv "github.com/amlwwalker/wingit/packages/server"
	utils "github.com/amlwwalker/wingit/packages/utils"
	"github.com/atrox/homedir"
)

//this handles interfacing with any business logic occuring elsewhere
type BusinessInterface struct {
	CONTROLLER                  *controller.CONTROLLER
	notifier                    NotificationHandler
	systray                     *Tray
	PendingUploadFileName       string //the file the user wants to upload (could be selected or dragged)
	PendingUpload               string //holds the filepath that is dragged over the drop area
	CurrentSelectedContact      string //their name
	CurrentSelectedContactIndex int    //their index (if needed)
	LoggedIn                    bool
	Status                      string
	pModel                      *PersonModel //list of contacts
	sModel                      *PersonModel //list of searched for
	fModel                      *FileModel   //list of files for a user
	dModel                      *FileModel   //list of  downloaded files
}

//handles the interface between the backend architecture
//and the bridge
func (b *BusinessInterface) configureInterface(signalLogin func(string, string), updateProcessStatus func(float64, bool), config utils.Config) {
	fmt.Printf("%+v\r\n", config)
	b.pModel = NewPersonModel(nil)
	b.sModel = NewPersonModel(nil)
	b.fModel = NewFileModel(nil)
	b.dModel = NewFileModel(nil)
	//configure the systray, do we therefore need the notifier
	b.systray = NewTray("systray_32x32.png")

	var server srv.Server
	var crypto cryptography.Crypto

	server.Ip = config.Host
	server.Address = config.Protocol + server.Ip
	fmt.Println("using: " + server.Address)
	server.Port = config.Port
	server.Verbose = config.Verbose
	fmt.Printf("%+v\r\n", server)
	PATH, err := homedir.Dir()
	if err != nil {
		fmt.Println("couldnt get home directory ", err)
		panic(err)
	}
	PATH = PATH + "/.wingit/"
	fmt.Println("path is: " + PATH)
	KEYFOLDER := PATH + ".keys/"
	DOWNLOADFOLDER := PATH + "downloads/" //no need to hide this
	SYNCFOLDER := PATH + ".sync/"

	// Initialise the controller. This will handle the user, server, crypto and ui.
	b.CONTROLLER = &controller.CONTROLLER{}
	b.CONTROLLER.DBPath = PATH //appends the db name once logged in (db name is per account)

	b.CONTROLLER.SERVER = &server
	b.CONTROLLER.CRYPTO = &crypto

	PASSWORDLENGTH := 32 // in bytes

	//need to implement logger in gui app if we want it to do anything front end
	logger := func(msg string) {
		fmt.Println("log: ", msg)
	}
	b.CONTROLLER.Logger = logger
	b.CONTROLLER.SERVER.Logger = b.CONTROLLER.Logger
	b.CONTROLLER.CRYPTO.Logger = b.CONTROLLER.Logger
	//make sure the directories exist
	utils.InitiateDirectory(KEYFOLDER)
	utils.InitiateDirectory(SYNCFOLDER)
	utils.InitiateDirectory(DOWNLOADFOLDER)
	b.CONTROLLER.CRYPTO.Init(PASSWORDLENGTH, KEYFOLDER, SYNCFOLDER, config.Verbose)
	b.CONTROLLER.SERVER.Init(DOWNLOADFOLDER, SYNCFOLDER)
	////setup the ConfigureDatabase
	b.CONTROLLER.DBDisabled = false
	//need to check to close the db here
	//first set the database connection up
	if err := b.CONTROLLER.InitialiseDatabase(); err != nil {
		fmt.Println("we cannot setup the database " + err.Error())
		b.CONTROLLER.DBDisabled = true
	}
	//make sure that generic storage buckets are created
	if err := controller.CreateGenericBuckets(b.CONTROLLER.Db); err != nil {
		fmt.Println("we cannot setup the database " + err.Error())
		b.CONTROLLER.DBDisabled = true
	} else {
		b.CONTROLLER.StoreStatusMessage("path set to: "+PATH, 1)
		b.CONTROLLER.StoreStatusMessage("database configured: ", 1)
		b.CONTROLLER.DBDisabled = false
	}
	//when a user logs in, we can then create/attempt to create a bucket
	//to store their contacts in
	// --------------------
	// SERVER
	user := &controller.User{}
	b.CONTROLLER.User = user
	b.CONTROLLER.Contacts = &controller.Contacts{}
	b.CONTROLLER.Contacts.People = make(map[string]*controller.Person)
	b.CONTROLLER.SearchResults = &controller.SearchResults{}

	//setting the indicator for the download progress
	//however will need access to the progress indicator
	b.CONTROLLER.SERVER.DownloadProgress = func(percentage float64) {
		fmt.Println("downloaded...: ", percentage, "%")
		if percentage > 0.98 {
			//turning this notificaiton off now there is synchronisation
			b.notifier.Push("Status", "Finished downloading")
		}
		updateProcessStatus(percentage, false)
	}
	//we dont know how long this will take
	//it would be nice to show a percentage of upload
	b.CONTROLLER.SERVER.UploadProgress = func(percentage float64, err error) {
		fmt.Println("uploading...")
		updateProcessStatus(percentage, true)
		if err != nil {
			b.notifier.Push("Upload Status", "Upload failed due to "+err.Error())
			return
		}
		if percentage > 0.98 {
			b.notifier.Push("Upload Status", "Upload complete")
		}
	}
	//need front end funtion to handle once a user is authd
	//special case requirement, as the front end will need to adapt to a login
	b.CONTROLLER.UserAuthenticated = func(user, err string) {
		//this firstly needs to update the UI
		//the needs to update the contacts to the currently logged in
		//user
		// so we have a user, lets pass it to the front end
		b.CONTROLLER.StoreStatusMessage("received yser: "+b.CONTROLLER.User.Name, 1)
		fmt.Println("received user: ", user)
		fmt.Println("user: ", b.CONTROLLER.User.Name)
		signalLogin(b.CONTROLLER.User.Name, err)
		if err != "" {
			//the user didn't auth properly, so dont bother contining
			return
		}

		//now we need to read the users contacts from the DB
		people, e := b.CONTROLLER.SyncPeople()
		if e != nil {
			b.CONTROLLER.Logger(e.Error())
		}
		fmt.Println("received: ", people)
		for k, v := range people {
			fmt.Println("adding ", k)
			var p = NewPerson(nil)
			p.SetFirstName("")
			p.SetLastName("")
			p.SetEmail(k)
			p.SetFileCount(strconv.Itoa(v.Len))
			b.pModel.AddPerson(p)
		}
		//at this stage we could also sync files
		//this could all be on a thread to make sure there is no blocking of the UI
	}
	b.CONTROLLER.StoreStatusMessage("business configured: ", 1)
}

func (b *BusinessInterface) SynchronizeWithServer(setMessage func(message string)) {
	fmt.Println("syncing files")
	b.CONTROLLER.RetrieveFilesForUser()

	//just reload all the contacts incase any new people have been added
	b.pModel.ClearPeople()
	for _, v := range b.CONTROLLER.Contacts.People {
		fmt.Println("after updating files, adding: ", v.UserId)
		addPersonToList(*v, b.pModel)
	}
	//here we should look for more keys

	//no contact is currently selected so no frontend file list to update
	if b.CurrentSelectedContact == "" {
		//no contact selected so we can't update the files
		return //don't let it try and update the files list
	}

	b.fModel.ClearFiles()
	for _, v := range b.CONTROLLER.Contacts.People[b.CurrentSelectedContact].Files {
		//now we add each one to the file list
		var f = NewFile(nil)
		f.SetFilePath(v.FileName)
		f.SetFileSize(strconv.Itoa(v.FileSize))
		f.SetFileSource(v.UserID)
		b.fModel.AddFile(f)
	}

	//this shows the user all the downloads they have
	if files, err := b.CONTROLLER.GetDownloadedFiles(); err != nil {
		//couldn't retrieve errors
		fmt.Println("error retrieving previously downloaded files", err)
	} else {
		b.dModel.ClearFiles()
		for _, v := range files {
			var f = NewFile(nil)
			f.SetFilePath(v.FileName)
			f.SetFileSize(strconv.Itoa(v.FileSize))
			b.dModel.AddFile(f)
		}
	}
	setMessage("sync at " + time.Now().String())
	// q.PopUpToast("sync at " + time.Now().String())

}
func (b *BusinessInterface) searchForMatches(regex string, informant func(float64, bool)) {
	//can do any preprocessing before it goes to the backend
	modelUpdater := func(c float64, indeterminate bool) {
		//if the logic is complete, then we need to update our model
		//with the search results
		//otherwise just inform the front end of progress
		if 1.0 == c { //complete
			//this is where you would add the contacts if you want async
			//updates to the UI
			fmt.Println("search is complete!")
		}
		//updates the front end
		informant(c, true) //we don't know how long it will take
	}
	modelUpdater = modelUpdater //dummy so it doesnt complain about lack of use
	//if we care about informing the front end how long this will take,
	//we can use the model updater above, but requires to update the cui aswell
	b.sModel.ClearPeople()
	go func() {
		if sResults, err := b.CONTROLLER.SearchForContacts(regex); err != nil {
			fmt.Println("error searching for contact: ", err)
		} else {
			//now these need adding to the search results
			for _, v := range sResults {
				//clicking on a search result adds them to the contacts
				fmt.Println("adding search result ", v)
				//TODO: Only show contacts not already in your list
				addPersonToList(v, b.sModel)
			}
			fmt.Println("updating front end that task is complete")
		}
		informant(1.0, true)
	}()
}

//example
//the interface needs to know how to inform the front end on progress
//so takes a function that takes a value that the front end will use
func (b *BusinessInterface) startAsynchronousRoutine(informant func(float64, bool)) {
	//on a go routine, count up to 10
	//each tick, inform the front end of your percentage progress
	//when it reaches ten, inform the front end it is complete

	go func() {
		var c float64
		c = 0.0
		for c < 1.0 {
			informant(c, false) //we know how long this process will take
			time.Sleep(1 * time.Second)
			c = c + 0.1
		}
		return
	}()
}
