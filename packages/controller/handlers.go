package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	handlers "github.com/amlwwalker/wingit/packages/handlers"
	utils "github.com/amlwwalker/wingit/packages/utils"
	"github.com/skratchdot/open-golang/open"
)

// ============================================================================================================================
func (c *CONTROLLER) UploadFileToContact(pendingContact, pendingFile string) (string, error) {
	if c.User.ApiKey == "" {
		//there is no logged in user.
		return "", errors.New("No logged in user")
	}
	// Make sure its a valid file-path.
	c.Logger("uploading file" + pendingFile + " for " + pendingContact)
	if !strings.Contains(pendingFile, "file:///") {
		pendingFile = "file:///" + pendingFile
		// return errors.New("File Name does not start correctly")
	}
	c.Logger("uploading file" + pendingFile + " for " + pendingContact)
	go func() {
		if url, err := handlers.ProcessFileEncryptAndUpload(pendingFile, pendingContact, c.User.ApiKey, c.CRYPTO, c.SERVER); err != nil {
			//if this is due to a file location failure, should clear contact out of database
			c.Logger("uploading failed due to " + err.Error())
		} else {
			fmt.Println("Success! File uploaded lives at ", url)
		}
		c.SERVER.UploadProgress(1.0)

	}()
	return "", nil
}
func (c *CONTROLLER) DownloadFileFromContact(pendingContact, pendingFile string) error {
	if c.User.ApiKey == "" {
		//there is no logged in user.
		return errors.New("No logged in user")
	}
	file := c.Contacts.People[pendingContact].Files[pendingFile]
	c.Logger("downloading file: " + file.FileName)

	go func() {
		err := handlers.DownloadAndDecryptFile(file.FileNameEnc, pendingContact, c.User.ApiKey, c.CRYPTO, c.SERVER)
		if err != nil {
			c.Logger("Error'd downloading file: " + err.Error())
		}
	}()

	return nil
}

func (c *CONTROLLER) OpenDownloadedFile(filePath string) error {
	fmt.Println("opening file: ", c.SERVER.DownloadFolder+filePath)
	err := open.Run(c.SERVER.DownloadFolder + filePath)
	return err
}
func (c *CONTROLLER) GetDownloadedFiles() ([]utils.File, error) {
	//scan the file system based on the file download location
	//get file name and file size
	//if a user clicks, we are going to open the file if we can
	var files []utils.File
	fileList, err := ioutil.ReadDir(c.SERVER.DownloadFolder)
	if err != nil {
		return files, err
	}
	//just for debugging
	for _, f := range fileList {
		var tmp utils.File
		tmp.FileName = f.Name()
		tmp.FileSize = int(f.Size())
		files = append(files, tmp)
	}
	return files, nil
}

func (c *CONTROLLER) RetrieveFilesForUser() error {
	if c.User.ApiKey == "" {
		//there is no logged in user.
		return errors.New("No logged in user")
	}
	// Trying to grab keys from the server
	// utils.PrintLine()
	c.Logger("OAuth successful! Trying to get keys from the server...")

	files, err := c.SERVER.GetFileList(c.User.ApiKey) // Returns an []utils.FileEnc
	// Error handling for the server request.
	if err != nil {
		c.Logger("Error fetching the files" + err.Error())
	} else if len(files) == 0 {
		c.Logger("No files found")
	} else {
		for _, file := range files {
			// fmt.Printf("retrieved file %+v\r\n", file)
			fileName, err := handlers.DecryptFileName(file.FileNameEnc, file.PasswordEnc, c.CRYPTO)
			if err != nil {
				c.Logger("Could not decrypt file name: " + file.FileNameEnc)
			} else {
				file.FileName = fileName
				// fmt.Println("file.Sender: ", file.Sender)
				c.AddFileToContact(file.Sender, file)
			}

		}
	}
	return nil

} // end of OAuth2Callback

// func (c *CONTROLLER) RetrieveAllKeys() error {
// 	if c.User.ApiKey == "" {
// 		//there is no logged in user.
// 		return errors.New("No logged in user")
// 	}

// 	keys, err := c.SERVER.GetKeys(c.User.ApiKey) // Returns an []utils.KeyServer
// 	// Error handling for the server request.
// 	if err != nil {
// 		c.Logger("Error fetching the keys" + err.Error())
// 	} else if len(keys) == 0 {
// 		utils.PrintError("No keys found")
// 	} else {
// 		for _, k := range keys {
// 			var p Person
// 			p.Name = k.UserID //don't have a name at this time
// 			p.UserId = k.UserID
// 			p.Len = len(p.Files)
// 			p.KeyId = k.ID

// 			c.AddContactToList(&p)
// 		}
// 	}
// 	return nil

// } // end of OAuth2Callback

// ============================================================================================================================

// EOF
