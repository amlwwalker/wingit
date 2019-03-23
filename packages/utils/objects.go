/*

...

*/

package utils

import (
	"time"
)

// ============================================================================================================================

type Config struct {
	Author   string `json:"author"`
	Path     string `json:"path"`
	Date     string `json:"date"`
	Mode     string `json:"mode"`
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
	Version  string `json:"version"`
	Port     string `json:"port"`
	Hotload  bool   `json:"hotload"`
	Verbose  bool   `json:"verbose"`
	Apikey   string `json:"apikey"`
}

type File struct {
	FileNameEnc string `json:"name"` // base64 encoded.
	FileName    string `json:"-"`
	FilePath    string
	ContentEnc  []byte `json:"content"` // not encoded?
	Content     []byte `json:"-"`
	PasswordEnc string `json:"password"`  // base64 encoded.
	Signature   string `json:"signature"` // base64 encoded.
	HMAC        string `json:"HMAC"`      // base64 encoded.
	UserID      string `json:"userID"`    // ? needed here?
	FileSize    int    `json:"fileSize"`
	// Does the server not do anything with empty fields?
	ID     int       `json:"ID"`
	Expiry time.Time `json:"expiry"`
	Sender string    `json:"sender"`
}

type Key struct {
	Content []byte `json:"content"`
}

type KeyServer struct {
	ID      int    `json:"ID"`
	Content []byte `json:"content"`
	UserID  string `json:"userID"`
}

// ============================================================================================================================

// EOF
