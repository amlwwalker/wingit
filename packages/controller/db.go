package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	utils "github.com/amlwwalker/wingit/packages/utils"
	"github.com/boltdb/bolt"
)

var (
	KEYS      = "keys"
	CONTACTS  = "contacts"
	USER      = "user"
	DOWNLOADS = "downloads"
	ERRORS    = "errors"
	STATUS    = "status"
)

type Error struct {
	Code    int
	Message string
	Time    time.Time
}
type Status struct {
	Code    int
	Message string
	Time    time.Time
}

func encode(store interface{}) ([]byte, error) {
	enc, err := json.Marshal(store)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func decode(store interface{}, data []byte) error {
	err := json.Unmarshal(data, &store)
	if err != nil {
		return err
	}
	return nil
}
func (c *CONTROLLER) InitialiseDatabase() error {

	//db is not up, so open it
	path := c.DBPath + ".wingit.db"
	var err error
	c.Db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return err
	} else if c.Db == nil {
		return errors.New("There is no db setup")
	}

	if s := c.Db.Path(); s != path {
		// c.Logger("path of db doesnt match " + s)
		return errors.New("unexpected path: " + s)
	}

	return nil
}
func CreateGenericBuckets(db *bolt.DB) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(ERRORS))
		tx.CreateBucketIfNotExists([]byte(STATUS))
		tx.CreateBucketIfNotExists([]byte(USER))
		tx.CreateBucketIfNotExists([]byte(CONTACTS))
		return nil
	}); err != nil {
		utils.PrintStatus("Error configuring db: " + err.Error())
		return err
	}
	return nil
}

//attempts to create a bucket for a user to store contacts in
func CreateUserContactBucket(userBucket string, db *bolt.DB) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		contacts, _ := tx.CreateBucketIfNotExists([]byte(CONTACTS))
		//nested bucket
		contacts.CreateBucketIfNotExists([]byte(userBucket))
		return nil
	}); err != nil {
		return err
	}
	return nil
}
func (c *CONTROLLER) CloseDb() error {
	if err := c.Db.Close(); err != nil {
		return err
	}
	return nil
}
func (c *CONTROLLER) RetrievePerson(userBucket, userId string) (*Person, error) {
	// if !c.Db.Open {
	//     return nil, fmt.Errorf("db must be opened before saving!")
	// }
	var p *Person
	if c.DBDisabled {
		fmt.Println("not storing as cant access db")
		return p, nil
	}
	err := c.Db.View(func(tx *bolt.Tx) error {
		var err error
		c := tx.Bucket([]byte(CONTACTS))
		u := c.Bucket([]byte(userBucket))

		k := []byte(userId)
		err = decode(p, u.Get(k))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// c.Logger("Could not get Person ID %s", userId)
		return nil, err
	}
	return p, nil
}
func (c *CONTROLLER) RetrieveAllPeople(userBucket string) ([]Person, error) {

	var people []Person

	if c.DBDisabled {
		fmt.Println("not storing as cant access db")
		return people, nil
	}

	err := c.Db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(CONTACTS))
		u := b.Bucket([]byte(userBucket))
		c := u.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			// fmt.Printf("key=%s, value=%s\n", k, v)
			var tmp Person
			err := decode(&tmp, u.Get(k))
			if err != nil {
				// fmt.Println("could not decode P" + err.Error())
				return err
			} else {
				people = append(people, tmp)
			}
		}
		return nil
	})
	return people, err
}
func (c *CONTROLLER) StorePerson(userBucket string, p *Person) error {
	if c.DBDisabled {
		fmt.Println("not storing as cant access db")
		return nil
	}
	if err := c.Db.Update(func(tx *bolt.Tx) error {
		c, _ := tx.CreateBucketIfNotExists([]byte(CONTACTS))
		b, _ := c.CreateBucketIfNotExists([]byte(userBucket))
		enc, err := encode(*p) //dereference the pointer
		if err != nil {
			return err
		}

		if err := b.Put([]byte(p.UserId), enc); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c *CONTROLLER) RetrieveUser() (*User, error) {
	// if !c.Db.Open {
	//     return nil, fmt.Errorf("db must be opened before saving!")
	// }
	var p User

	if c.DBDisabled {
		fmt.Println("not storing as cant access db")
		return &p, nil
	}

	err := c.Db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(USER))
		k := []byte("CURRENT_USER")
		if err = decode(&p, b.Get(k)); err != nil {
			fmt.Println("error: ", err)
		} else {
			fmt.Println("p: ", p)
			return nil

		}

		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// c.Logger("Could not get Person ID %s", userId)
		return nil, err
	}
	return &p, nil
}

func (c *CONTROLLER) StoreUser(u *User) error {
	if c.DBDisabled {
		fmt.Println("not storing as cant access db")
		return nil
	}
	// Insert until we get above the minimum 4MB size.
	if err := c.Db.Update(func(tx *bolt.Tx) error {
		user, _ := tx.CreateBucketIfNotExists([]byte(USER))
		enc, err := encode(*u) //dereference the pointer
		if err != nil {
			return err
		}

		if err := user.Put([]byte("CURRENT_USER"), enc); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
func (c *CONTROLLER) RemoveCurrentUser() error {
	if c.DBDisabled {
		fmt.Println("not storing as cant access db")
		return nil
	}
	// Insert until we get above the minimum 4MB size.
	if err := c.Db.Update(func(tx *bolt.Tx) error {
		user, _ := tx.CreateBucketIfNotExists([]byte(USER))
		if err := user.Delete([]byte("CURRENT_USER")); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
func (c *CONTROLLER) StoreErrorMessage(message string, code int) error {
	if c.DBDisabled {
		fmt.Println("not storing as cant access db")
		return nil
	}
	var e Error
	e.Code = code
	e.Message = message
	e.Time = time.Now()
	// Insert until we get above the minimum 4MB size.
	if err := c.Db.Update(func(tx *bolt.Tx) error {
		errors, _ := tx.CreateBucketIfNotExists([]byte(ERRORS))
		enc, err := encode(e)
		if err != nil {
			return err
		}

		if err := errors.Put([]byte(e.Time.Format(time.RFC3339)), enc); err != nil {
			return err
		}
		return nil
		// if err := errors.Put([]byte(e.Time.Format(time.RFC3339)), []byte(e)); err != nil {
		// 	return err
		// }
		// return nil
	}); err != nil {
		return err
	}
	return nil
}
func (c *CONTROLLER) StoreStatusMessage(message string, code int) error {
	if c.DBDisabled {
		fmt.Println("not storing as cant access db")
		return nil
	}
	var s Status
	s.Code = code
	s.Message = message
	s.Time = time.Now()
	// Insert until we get above the minimum 4MB size.
	if err := c.Db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(STATUS))
		enc, err := encode(s)
		if err != nil {
			return err
		}

		if err := b.Put([]byte(s.Time.Format(time.RFC3339)), enc); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
