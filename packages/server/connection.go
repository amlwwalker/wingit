/*

Handles all connections with the server.

*/

package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	tus "github.com/amlwwalker/go-tus"
	utils "github.com/amlwwalker/wingit/packages/utils"
)

// ============================================================================================================================

// Objects defined in objects.go

// ============================================================================================================================

// Function to make get-connection with the server. It will return the contents.
func (srv *Server) GetConnection(path string, idToken string) ([]byte, error) {

	var err error

	URL := srv.Address + ":" + srv.Port + path + "?id_token=" + idToken
	fmt.Println(URL)
	if srv.Verbose {
		// srv.Logger("Making GET request to " + URL)
		// utils.PrintStatus("Making GET connection to: " + URL)
	}

	// NOTE: the path needs to contain the leading slash
	response, err := http.Get(URL)
	if err != nil {
		utils.PrintErrorFull("Error making GET request to server", err)
		return nil, err
	} else {
		// Create our progress reporter and pass it to be used alongside our writer
		readerpt := &PassThru{Reader: response.Body, length: response.ContentLength}
		//perhaps pointer is better, this is a pass
		readerpt.DownloadProgress = srv.DownloadProgress

		defer response.Body.Close()
		respBody, err := ioutil.ReadAll(readerpt)
		if err != nil {
			utils.PrintErrorFull("Error reading response body after GET request", err)
			// os.Exit(1) // This should error properly.
			return nil, err
		}
		return respBody, nil
	}

	return nil, nil // Redundancy.

} // end of GetConnection

func (srv *Server) MultiPartUpload(path string, data []byte, metaData map[string]string, idToken string) (string, string, error) {

	URL := srv.Address + ":" + srv.Port + path + "?id_token=" + idToken

	if srv.Verbose {
		srv.Logger("Making POST connection to: " + URL)
	}

	// fmt.Println("URL conneting with ", URL)
	// create the tus client.
	client, _ := tus.NewClient(URL, nil)

	// create an upload from a file.
	upload := tus.NewUploadFromBytes(data)

	upload.Metadata = metaData
	// fmt.Println("the metadata on send ", upload.Metadata)
	var uploader *tus.Uploader
	var err error
	// create the uploader.
	if uploader, err = client.CreateUpload(upload); err != nil {
		fmt.Println("there was an error uploading ", err)
		return "504", "", err
	}
	// start the uploading process.
	if err := uploader.Upload("?id_token=" + idToken); err != nil {
		return "503", "", err
	}

	return "201 OK", uploader.Url(), nil
} // end of Multipart upload

func (srv *Server) PostConnection(path string, data []byte, idToken string) (string, []byte, error) {

	var err error
	var emptyString string

	URL := srv.Address + ":" + srv.Port + path + "?id_token=" + idToken

	if srv.Verbose {
		srv.Logger("Making POST connection to: " + URL)
	}

	req, err := http.NewRequest("POST", URL, bytes.NewReader(data))
	if err != nil {
		utils.PrintErrorFull("Error setting up the POST request", err)
		return emptyString, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.PrintErrorFull("Error making POST request", err)
		return emptyString, nil, err
	}
	defer resp.Body.Close()
	// Create our progress reporter and pass it to be used alongside our writer
	srv.Logger("data size " + string(len(data)))
	// readerpt := &PassThru{Reader: resp.Body, length: int64(len(data))}
	//perhaps pointer is better, this is a pass
	// readerpt.DownloadProgress = srv.DownloadProgress

	// respBody, _ := ioutil.ReadAll(readerpt)
	respBody, _ := ioutil.ReadAll(resp.Body)
	// srv.Logger("body: " + string(respBody))
	if srv.Verbose {
		// srv.Logger(" ~~~ POST Response Status: ", resp.Status)
		// srv.Logger(" ~~~ POST Response Headers: ", resp.Header)
	}
	srv.UploadProgress(1.0, nil)
	return resp.Status, respBody, nil

} // end of PostConnection

// ============================================================================================================================

// EOF
