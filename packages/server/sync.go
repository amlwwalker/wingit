/*

Handles all connections with the server.

*/

package server

import (
	"encoding/json"
	"errors"
	"fmt"

	utils "github.com/amlwwalker/wingit/packages/utils"
)

// ============================================================================================================================

// Objects defined in objects.go

// ============================================================================================================================

func (srv *Server) SyncLocalFile(pathToFile, idToken string) error {

	srv.Logger("Loading locally stored file `" + pathToFile + "`...")

	contentBytes, err := utils.ReadFromFile(pathToFile)
	if err != nil {
		srv.Logger("Error loading an encrypted file from temporary storage" + err.Error())
		return err
	}

	var file utils.File
	err = json.Unmarshal(contentBytes, &file)
	if err != nil {
		srv.Logger("Error unmarshalling the file" + err.Error())
		return err
	}

	// Send to server.
	status, url, err := srv.PostFile(file, idToken)
	if err != nil {
		srv.Logger("Error uploading file to server" + err.Error())
		return err
	}
	// Is this check sufficient to check for success?
	if status == "201 OK" {
		// If successful, delete file.
		srv.Logger("Deleting file from local storage.")
		_ = utils.DeleteFile(pathToFile)
	} else {
		err = errors.New("Error sending to the server. ")
		srv.Logger("Not deleting local file." + err.Error())
		return err
	}
	fmt.Println("url for synced file is ", url)

	return nil

} // end of SyncLocalFile

func (srv *Server) SyncLocalKey(pathToFile, idToken string, modify bool) error {

	contentBytes, err := utils.ReadFromFile(pathToFile)
	if err != nil {
		srv.Logger("Error loading a public-key from temporary storage" + err.Error())
		return err
	}

	var key utils.Key
	err = json.Unmarshal(contentBytes, &key)
	if err != nil {
		srv.Logger("Error unmarshalling the key" + err.Error())
		return err
	}

	// Send to server.
	status, err := srv.PostKeyBytes(key.Content, idToken, modify)
	if err != nil {
		srv.Logger("Error uploading key to server" + err.Error())
		return err
	}
	// Is this check sufficient to check for success?
	if status == "200 OK" {
		// If successful, delete file.
		srv.Logger("Deleting public-key from local storage.")
		_ = utils.DeleteFile(pathToFile)
	} else {
		err = errors.New("Error sending to the server. ")
		srv.Logger("Not deleting local file." + err.Error())
		return err
	}

	return nil

}

// ============================================================================================================================

// EOF
