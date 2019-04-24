/*

Handles all connections with the server.

*/

package server

import (
	"encoding/json"

	utils "github.com/amlwwalker/wingit/packages/utils"
)

// ============================================================================================================================

// Objects defined in objects.go

// ============================================================================================================================

// Get a list of all available keys from the server. - This will go in the future.
// NOTE: This is now listing all public keys.
//       Need to limit this to all keys for the logged-in user.
func (srv *Server) GetKeys(regexSearch, idToken string) ([]utils.KeyServer, error) { // NOTE: FUNCTIONALITY WILL BE REMOVED

	var err error
	var keys []utils.KeyServer

	if srv.Verbose {
		utils.PrintStatus("Retrieving all the public-keys from the server...")
	}

	contentBytes, err := srv.GetConnection("/keys/search", idToken+"&match="+regexSearch)
	if err != nil {
		utils.PrintErrorFull("Error retrieving all the public-keys from the server", err)
		return keys, err
	}

	if string(contentBytes) == "{}" {
		return keys, nil
	}

	if err := json.Unmarshal(contentBytes, &keys); err != nil {
		utils.PrintErrorFull("Error unmarshalling all the public-keys from the server", err)
		return keys, err
	}

	if srv.Verbose {
		utils.PrintSuccess("Successfully retrieved all the public-keys for")
	}

	return keys, nil

} // end of GetKeys

// Get  public-key from the server.
// NOTE: This can be an array!!!
func (srv *Server) GetKey(regexSearch, idToken string) (utils.KeyServer, error) {

	var key utils.KeyServer

	if srv.Verbose {
		// utils.PrintStatus("Retrieving public-key for `" + matchEmail + "` from the server...")
	}

	contentBytes, err := srv.GetConnection("/keys/retrieve", idToken+"&match="+regexSearch)
	if err != nil {
		// utils.PrintErrorFull("Error retrieving the public-key from the server", err)
		return key, err
	}

	if string(contentBytes) == "{}" {
		return key, nil
	}

	if err := json.Unmarshal(contentBytes, &key); err != nil {
		// utils.PrintErrorFull("Error unmarshalling the public-key from the server", err)
		return key, err
	}

	if srv.Verbose {
		// utils.PrintSuccess("Successfully retrieved the public-key for `" + matchEmail + "`")
	}

	return key, nil

} // end of GetKey

func (srv *Server) PostKeyBytes(keyBytes []byte, idToken string, modify bool) (string, error) {

	var str string
	var key utils.Key

	if srv.Verbose {
		// utils.PrintStatus("Posting public-key to the server...")
	}

	// Set the contents of the key
	key.Content = keyBytes

	// Marshall the struct into JSON
	keyPayloadJSON, err := json.Marshal(key)
	if err != nil {
		// utils.PrintErrorFull("Error marshalling key for upload", err)
		return str, err
	}
	//is this a mad approach?
	var approach string
	if modify {
		approach = "modify"
	} else {
		approach = "upload"
	}
	status, _, err := srv.PostConnection("/keys/"+approach, keyPayloadJSON, idToken)
	if err != nil {
		utils.PrintErrorFull("Error sending the key to the server", err)
		return str, err
	}

	if srv.Verbose {
		if status == "200 OK" {
			// utils.PrintSuccess("Successfully uploaded the public-key to the server...")
		} else {
			// utils.PrintError("Error sending the public-key, status code " + status + " received on attempt.")
		}
	}

	return status, nil

} // PostKeyBytes

// ============================================================================================================================

// EOF
