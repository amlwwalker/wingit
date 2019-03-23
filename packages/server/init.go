/*

Handles all communications with the server.

*/

package server

import (
	utils "github.com/amlwwalker/wingit/packages/utils"
)

// ============================================================================================================================

// Objects defined in objects.go

// ============================================================================================================================

// Function to initialiase the server setup.
func (srv *Server) Init(downloadFolder, syncFolder string) {
	srv.DownloadFolder = downloadFolder
	srv.SyncFolder = syncFolder
	utils.PrintStatus("Initialising download location...")
	//this should just initialise the storage of the download files (by creating the directory)

} // end of Init

// ============================================================================================================================

// EOF
