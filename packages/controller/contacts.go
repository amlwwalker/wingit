package controller

import (
	"errors"
	"fmt"

	utils "github.com/amlwwalker/wingit/packages/utils"
)

func (c *CONTROLLER) SearchForContacts(searchQuery string) ([]Person, error) {
	//use the server to search for a query
	//and return the people associated with that key
	var people []Person
	if c.User.ApiKey == "" {
		//there is no logged in user.
		return []Person{}, errors.New("No logged in user")
	}
	fmt.Println("search query: " + searchQuery)
	if keys, err := c.SERVER.GetKeys(searchQuery, c.User.ApiKey); err != nil {
		return []Person{}, err
	} else {
		err := c.CRYPTO.SavePublicKeys(keys) //TODO: This should really only save keys for contacts they want, not all!
		if err != nil {
			c.Logger("Error saving the public keys" + err.Error())
		}
		for _, v := range keys {
			var tmp Person
			tmp.Name = v.UserID
			// tmp.Files = []utils.File
			tmp.UserId = v.UserID
			fmt.Println("adding user: " + v.UserID)
			//at this point we have the key, but we aren't actually storing it?
			people = append(people, tmp)
		}
		fmt.Println("returning: ", len(people), " as search results")
		return people, nil
	}
}

func (c *CONTROLLER) AddFileToContact(sender string, file utils.File) error {
	if c.User.ApiKey == "" {
		//there is no logged in user.
		return errors.New("No logged in user")
	}
	if _, ok := c.Contacts.People[sender]; !ok {
		//do something here
		var p Person
		p.Name = sender //don't have a name at this time
		p.UserId = sender
		p.Files = make(map[string]*utils.File)
		//at this stage we know we haven't got their key, so we should by rights request it
		if keys, err := c.SERVER.GetKey(sender, c.User.ApiKey); err != nil {
			c.Logger("Couldn't get key for " + sender + " " + err.Error())
		} else {
			err := c.CRYPTO.SavePublicKeys([]utils.KeyServer{keys})
			if err != nil {
				c.Logger("Error saving the public keys" + err.Error())
			}
			c.AddContactToList(&p)

		}
	}

	c.Logger("storing file for " + sender)
	c.Contacts.People[sender].Files[file.FileName] = &file
	c.Contacts.People[sender].Len = len(c.Contacts.People[sender].Files)
	return nil
	// }
	// }
	//can only have got here if there was no match
	//create the contact
	// var p *Person
	// p.Name = sender //don't have a name at this time
	// p.UserId = sender
	// p.Files = append(p.Files, file)
	// p.Len = len(p.Files)

	// c.AddContactToList(p)
	// return
}

func (c *CONTROLLER) SyncPeople() (map[string]*Person, error) {
	if c.User.ApiKey == "" {
		//there is no logged in user.
		return map[string]*Person{}, errors.New("No logged in user")
	}
	//retrieve from DB and store here in model
	people, err := c.RetrieveAllPeople(c.User.UserId.String())
	if err != nil {
		c.Logger("error retrieving people" + err.Error())
	}
	//we need to clean them up so they are just strings
	// var names []string
	for k, _ := range people {
		// var p Person
		// p.Name = v.UserId //don't have a name at this time
		// p.UserId = v.UserId
		// p.Files = make(map[string]*utils.File)
		// p.Len = 0
		//messy because this happens below aswell
		//cant do both as below will store it
		//this is a nasty catch
		//v only exists for one iteration
		//and gets over written each time
		//but we are pointing to the memory location
		//so they all end up at the same memory location
		//which means we have to point to the original data
		c.Contacts.People[people[k].UserId] = &people[k]
		// names = append(names, v.UserId)
	}
	fmt.Println("people now are: ", c.Contacts.People)
	c.Contacts.Len = len(c.Contacts.People)
	return c.Contacts.People, nil
}
func (c *CONTROLLER) CreatePersonFromIdentifier(identifier string) *Person {
	var p Person
	p.Name = identifier
	p.UserId = identifier
	p.Files = make(map[string]*utils.File)
	return &p
}

func (c *CONTROLLER) DeleteContactFromList(identifier string) {
	fmt.Println("deleting contact ", identifier)
	//this deletes it from the list
	delete(c.Contacts.People, identifier) //memory deletion
	//we now also need to delete the contact from the database
	c.deleteContactForUser(identifier) //database deletion

}
func (c *CONTROLLER) AddContactToList(p *Person) {

	//first create a file owned by the sender
	// var f File
	// f.FileName = fileNameDec
	// f.FileNameEnc = fileNameEnc

	// //the create the sender who owns the files
	// var p Person
	// p.Name = contactName
	// p.Files = append(p.Files, f)
	// p.Len = len(p.Files)

	//finally attach them to the ctrl object to send the the frontend
	c.Contacts.People[p.UserId] = p //append(c.Contacts.People, p)
	c.Contacts.Len = len(c.Contacts.People)
	c.StorePerson(c.User.UserId.String(), p)
	// fmt.Println(ctrl.Contacts.People)
	// qml.Changed(ctrl, &ctrl.Contacts)

}
