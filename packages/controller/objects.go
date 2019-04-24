//

package controller

import (
	"time"

	cryptography "github.com/amlwwalker/wingit/packages/cryptography"
	srv "github.com/amlwwalker/wingit/packages/server"
	utils "github.com/amlwwalker/wingit/packages/utils"
	"github.com/boltdb/bolt"
	"github.com/gobuffalo/packr"
	uuid "github.com/satori/go.uuid"
)

// ============================================================================================================================

type CONTROLLER struct {
	SERVER            *srv.Server
	CRYPTO            *cryptography.Crypto
	User              *User
	UserAuthenticated func(user string, error string) //when oauth is complete can call this
	Contacts          *Contacts
	Logger            func(message string)
	SearchResults     *SearchResults
	Db                *bolt.DB
	DBPath            string
	DBDisabled        bool
	WebAssets         packr.Box
}

// type UserProfile struct {
//     Id string `"json":"user_id"`
//     Email string `"json":"email"`
//     Name string `"json":"name"`
//     Picture string `"json":"picture"`
//     Locale string `"json":"locale"`
//     ApiKey string `"json":"apiKey"`
// }

// Child of the CONTROLLER object
// type User struct {
// 	Id      string `"json":"user_id"`
// 	Email   string `"json":"email"`
// 	Name    string `"json":"name"`
// 	Picture string `"json":"picture"`
// 	Locale  string `"json":"locale"`
// 	ApiKey  string `"json":"-"`
// 	Service string `"json":"service"`
// }

type User struct {
	UserId *uuid.UUID `json:"user_id"` //our id
	// Id        string     `json:"id"`      //a services id (e.g google)
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"-" gorm:"not null;unique"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	Locale    string    `json:"-"`
	ApiKey    string    `json:"apiKey"`
	Service   string    `json:"service"`
	CreatedAt time.Time `json:"createdAt"`
	Verified  bool      `json:"verified"` //cannot send files just receive files if not verified (anonymous)
	// Expiry    time.Time `json:"expiry"`
	Enabled bool `json:"enabled"`
}

type SearchResults struct {
	Results []string
	Len     int
}

type Contacts struct {
	People map[string]*Person
	Len    int
}

type Person struct {
	Name   string
	UserId string //email
	Files  map[string]*utils.File
	Len    int
	KeyId  int
}

// type File struct {
//     FileName    string
//     FileNameEnc string
//     FileId      string
//     FileSize    int
// }
// type File struct {
//     FileNameEnc         string  `json:"name"` // base64 encoded.
//     FileName    string  `json:"-"`
//     ContentEnc      []byte      `json:"content"` // not encoded?
//     PasswordEnc     string      `json:"password"` // base64 encoded.
//     Signature       string      `json:"signature"` // base64 encoded.
//     HMAC            string      `json:"HMAC"` // base64 encoded.
//     UserID          string      `json:"userID"` // ? needed here?
//     FileSize        int `json:"fileSize"`
//     // Does the server not do anything with empty fields?
//     ID              int         `json:"ID"`
//     Expiry          time.Time   `json:"expiry"`
//     Sender          string      `json:"sender"`
// }
// ============================================================================================================================

// EOF
