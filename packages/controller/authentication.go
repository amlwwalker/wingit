/*

...

*/

package controller

import (
	"encoding/json"
	"errors"
	"fmt"

	handlers "github.com/amlwwalker/wingit/packages/handlers"
)

func (c *CONTROLLER) Logout(callback func()) {
	//this just deletes the local data about the user (the token from the DB)
	//if user is logged in....
	//so this needs to clear the entry in the database
	//it needs to clear any reminance of c.User being anything
	//it needs to inform the front end
	c.User = &User{}
	c.RemoveCurrentUser() //something like this....
	//and responds to the front end to inform logout is complete
	callback()
}

func (c *CONTROLLER) CheckForExistingUser() error {
	if user, err := c.RetrieveUser(); err != nil {
		fmt.Println("couldn't retrieve user (perhaps there is none at this stage) ", err)
		return err
	} else {
		fmt.Println("user found: ", user)
		c.User = user
		if c.User.UserId.String() == "" {
			//no user was found, don't continue
			fmt.Println("no user found: ", c.User.UserId.String())
			return errors.New("no user was found")
		}
		c.StoreStatusMessage("doing user specific setup: ", 1)
		c.UserSpecificSetup(true)

		//return to the front end the user has authenticated
		//do this because this function shouldn't know what to do about the front end
		c.UserAuthenticated(c.User.Name, "")
		return nil
	}
}
func (c *CONTROLLER) Authorize(apiKey string) (User, error) {

	var userProfile User

	contentBytes, err := c.SERVER.GetConnection("/user/request", apiKey)
	if err != nil {
		// utils.PrintErrorFull("Error retrieving the public-key from the server", err)
		return userProfile, err
	}

	if string(contentBytes) == "{}" {
		return userProfile, nil
	}

	if err := json.Unmarshal(contentBytes, &userProfile); err != nil {
		// utils.PrintErrorFull("Error unmarshalling the public-key from the server", err)
		return userProfile, err
	}
	c.Logger("User seems logged in" + userProfile.Email)
	c.Logger("User seems logged in" + userProfile.Name)

	// Extract the id_token so that we can auth against our server
	c.User = &userProfile
	fmt.Println("id token: ", c.User.ApiKey)

	c.StoreUser(c.User)
	//last thing to do is update the UI with the information
	c.UserSpecificSetup(false) //modify currently DOES NOT work

	//return to the front end the user has authenticated
	//do this because this function shouldn't know what to do about the front end
	c.UserAuthenticated(c.User.Name, "")
	return userProfile, nil
} // end of authentication

func (c *CONTROLLER) UserSpecificSetup(modify bool) {
	c.StoreStatusMessage("starting setup: "+c.User.UserId.String(), 1)
	c.CRYPTO.InitRSAKeyPair(c.User.UserId.String())
	//modify ? over write keys on server or not (debug mode)
	handlers.SyncTemporaryStorageFolder(c.SERVER, c.User.ApiKey, modify)

	//confirm a user specific contact bucket has been created
	if err := CreateUserContactBucket(c.User.UserId.String(), c.Db); err != nil {
		fmt.Println("could not create a contact bucket for the user: ", err)
	}
}

// ============================================================================================================================

// EOF
