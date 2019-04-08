/*

Handles all file related endpoint connections with the server.

*/

package server

import (
	// "fmt"
	"encoding/json"

	utils "github.com/amlwwalker/wingit/packages/utils"
)

// ============================================================================================================================

// Objects defined in objects.go

// ============================================================================================================================

// Get a list of all files, that are shared with you, from the server.
// NOTE: Returns encrypted filenames.
func (srv *Server) GetFileList(idToken string) ([]utils.File, error) {

	var err error
	var files []utils.File

	if srv.Verbose {
		srv.Logger("Retrieving all files from the server...")
	}

	contentBytes, err := srv.GetConnection("/files/list", idToken)
	if err != nil {
		utils.PrintErrorFull("Error in receiving the file-list from the server", err)
		return files, err
	}

	if srv.Verbose {
		srv.Logger("Received the content for all files from the endpoint.")
	}

	// Disgusting handling of empty response object
	if string(contentBytes) == "{}" {
		return files, nil
	}

	// Unmarshall all the information into the files.
	if err := json.Unmarshal(contentBytes, &files); err != nil {
		utils.PrintErrorFull("Error in unmarshalling the content of the filenames", err)
		return files, err
	}

	srv.Logger("Received all files from the endpoint.")

	// This should return an array of all the files including the filenames from the bucket
	// NOTE! These filenames are encrypted. It does not return any encrypted file-content.
	return files, nil

} // end of GetFileList

// Get a single file from the server.
// It returns the encrypted file-contents, name, password, signature, hmac.
func (srv *Server) GetFile(fileNameEnc, idToken string) (utils.File, error) {

	var err error
	var fileEnc utils.File

	if srv.Verbose {
		srv.Logger("Retrieving encrypted file " + fileNameEnc + " from the server...")
	}

	contentBytes, err := srv.GetConnection("/files/retrieve", idToken+"&filename="+fileNameEnc)
	if err != nil {
		utils.PrintErrorFull("Error retrieving the file from the server", err)
		return fileEnc, err
	}

	if err := json.Unmarshal(contentBytes, &fileEnc); err != nil {
		utils.PrintErrorFull("Error unmarshalling the retrieved file", err)
		return fileEnc, err
	}

	if srv.Verbose {
		srv.Logger("Successfully retrieved encrypted file from the server.")
	}

	return fileEnc, nil

} // end of GetFile

// // Wrapper for PostFile
// func (srv *Server) UploadFile(fileEnc utils.FileEnc) (string, error) {
//     return srv.PostFile(fileEnc)
// }

/* the data is in this form
var payload utils.File
payload.ContentEnc = encryptedBytes
payload.FileNameEnc = utils.EncodeBase64(fileNameEncryptedBytes)
payload.PasswordEnc = utils.EncodeBase64(passwordEncryptedBytes)
payload.Signature = utils.EncodeBase64(signatureBytes)
payload.HMAC = utils.EncodeBase64(hmacBytes)
payload.UserID = toUserId
*/

// Function to post an encrypted file to the server.
func (srv *Server) PostFile(fileEnc utils.File, idToken string) (string, string, error) {

	var str string

	if srv.Verbose {
		srv.Logger("Attempting to post file " + fileEnc.FileNameEnc + " to the server...")
	}

	payloadJSON, err := json.Marshal(fileEnc)
	if err != nil {
		utils.PrintErrorFull("Error marshalling encrypted file data", err)
		return "", str, err
	}

	metaData := map[string]string{}
	metaData["FileNameEnc"] = fileEnc.FileNameEnc
	metaData["PasswordEnc"] = fileEnc.PasswordEnc
	metaData["Signature"] = fileEnc.Signature
	metaData["HMAC"] = fileEnc.HMAC
	metaData["UserID"] = fileEnc.UserID
	// fmt.Println("metadata set for ", metaData["UserID"])
	// Post to the server
	status, fileUrl, err := srv.MultiPartUpload("/files/upload/", payloadJSON, metaData, idToken)
	if err != nil {
		utils.PrintErrorFull("Error posting the (encrypted) file-bytes", err)
		return "", status, err
	}

	if srv.Verbose {
		if status == "201 OK" {
			srv.Logger("Successfully uploaded " + fileEnc.FileNameEnc + " to the server...")
		} else {
			utils.PrintError("Error sending file, status code " + status + " received on attempt.")
		}
	}

	return status, fileUrl, nil

} // end of PostFile

// ============================================================================================================================

// EOF
