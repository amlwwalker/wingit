/*

This holds all the RSA functions

*/

package cryptography

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	utils "github.com/amlwwalker/wingit/packages/utils"
)

// ============================================================================================================================

// Objects defined in objects.go

// ============================================================================================================================

// PUBLIC FUNCTIONS

func (c *Crypto) DeletePublickKeyForContact(path string, fileName string) error {
	if err := utils.DeleteFile(path + fileName); err != nil {
		fmt.Println("error deleting key " + path + fileName)
		return err
	}
	return nil
}

// Loads a public-key using a file-name. Assumes they are stored in the keys folder.
// Returns a pointer to that public-key.
func (c *Crypto) LoadPublicKey(path string, fileName string) (*rsa.PublicKey, error) {

	var publicKey *rsa.PublicKey

	if c.Verbose {
		c.Logger("Loading private-key file...")
	}

	publicKeyBytes, err := utils.ReadFromFile(path + fileName)
	if err != nil {
		c.Logger("Error opening the public-key file" + err.Error())
		return publicKey, err
	}

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		c.Logger("Error decoding public-key file. block == nil")
	}

	pubK, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		c.Logger("Error pasing (PKIX) public-key file" + err.Error())
	}

	switch pubK := pubK.(type) {
	case *rsa.PublicKey:
		if c.Verbose {
			c.Logger("Successfully loaded private-key file.")
		}
		return pubK, nil
	// case *dsa.PublicKey:
	//     panic(err)
	// case *ecdsa.PublicKey:
	//     utils.PrintError("Error loading public-key. We currently dont support ECDSA.")
	//     err = errors.New("Dont support ECDSA")
	//     return publicKey, err
	default:
		err = errors.New("Unknown keytype")
		c.Logger("Error assessing public-key type" + err.Error())
		return publicKey, err
	}

	if c.Verbose {
		c.Logger("Successfully loaded public-key.")
	}

	return publicKey, nil

} // end of LoadPublicKey

// Function that returns the bytes for the publicKey so it can be sent to the server
func (c *Crypto) GetPublicKeyBytes() ([]byte, error) {

	PubASN1Bytes, err := x509.MarshalPKIXPublicKey(c.PublicKey)
	if err != nil {
		c.Logger("Error marshalling the RSA public key" + err.Error())
		return nil, err
	}

	return PubASN1Bytes, nil

} // end of GetPublicKeyBytes

// Function to save public-keys that have been retrieved from the server
func (c *Crypto) SavePublicKeys(keyList []utils.KeyServer) error {

	if c.Verbose {
		c.Logger("Attempting to save all public-keys...")
	}

	// The input is an array...
	for _, key := range keyList {

		keyName := key.UserID
		keyBytes := key.Content

		// Create the file.
		pubPem, err := os.Create(c.KeyFolder + keyName + ".pem")
		if err != nil {
			c.Logger("Error creating the placeholder RSA public key file" + err.Error())
			return err
		}

		var pubKey = &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: keyBytes,
		}
		err = pem.Encode(pubPem, pubKey)
		if err != nil {
			c.Logger("Error saving RSA public key" + err.Error())
			return err
		}

		pubPem.Close()

	}

	if c.Verbose {
		c.Logger("Successfully saved all public-keys.")
	}

	return nil

} // end of SavePublicKeys

// ============================================================================================================================

// PRIVATE FUNCTIONS

func (c *Crypto) loadRSAKeyPair(id string) error {

	var err error
	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey

	if c.Verbose {
		c.Logger("Loading local key-pair files...")
	}

	privateKey, err = c.loadPrivateKey(id)
	if err != nil {
		c.Logger("Error loading the private-key file" + err.Error())
		return err
	}

	privateKey.Precompute() // Speeds up some calculations later on

	publicKey = &privateKey.PublicKey

	// Store in the struct.
	c.PrivateKey = privateKey
	c.PublicKey = publicKey

	if c.Verbose {
		c.Logger("Successfully loaded local key-pair files.")
	}

	return nil

} // end of loadRSAKeyPair

func (c *Crypto) saveRSAPrivateKey(id string) error {

	if c.Verbose {
		c.Logger("Saving private-key file...")
	}

	// Create the file.
	privateKeyFile, err := os.Create(c.KeyFolder + id + "-private.key")
	if err != nil {
		return err
	}

	//
	privateKeyEncoder := gob.NewEncoder(privateKeyFile)
	privateKeyEncoder.Encode(c.PrivateKey)
	privateKeyFile.Close()

	if c.Verbose {
		c.Logger("Successfully saved private-key file.")
	}

	return nil

} // end of saveRSAPrivateKey

// Saves the public pem to share with others or on a server
func (c *Crypto) saveRSAPublicKey(path string, id string, publicKey *rsa.PublicKey) error {

	if c.Verbose {
		c.Logger("Saving public-key file...")
	}

	// Create the file.
	pubPem, err := os.Create(path + id + "-public.pem")
	if err != nil {
		c.Logger("Error creating the placeholder RSA public key file" + err.Error())
		return err
	}

	// From: http://stackoverflow.com/questions/13555085/save-and-load-crypto-rsa-privatekey-to-and-from-the-disk
	PubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		c.Logger("Error marshalling the RSA public key" + err.Error())
		return err
	}

	// http://golang.org/pkg/encoding/pem/#Block
	var pubKey = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: PubASN1,
	}
	err = pem.Encode(pubPem, pubKey)
	if err != nil {
		c.Logger("Error saving RSA public key" + err.Error())
		return err
	}

	pubPem.Close()

	if c.Verbose {
		c.Logger("Successfully saved public-key file.")
	}

	return nil

} // end of saveRSAPublicKey

// Loads a private key if necessary from a .key file
func (c *Crypto) loadPrivateKey(id string) (*rsa.PrivateKey, error) {

	var privateKey rsa.PrivateKey

	if c.Verbose {
		c.Logger("Loading private-key file...")
	}

	privateKeyFile, err := os.Open(c.KeyFolder + id + "-private.key")
	if err != nil {
		c.Logger("Eror loading the private-key file" + err.Error())
		return &privateKey, err
	}

	decoder := gob.NewDecoder(privateKeyFile)
	err = decoder.Decode(&privateKey)
	if err != nil {
		c.Logger("Eror decoding the private-key file" + err.Error())
		return &privateKey, err
	}

	privateKeyFile.Close()

	if c.Verbose {
		c.Logger("Successfully loaded private-key file.")
	}

	return &privateKey, nil

} // end of loadPrivateKey

// ============================================================================================================================

// EOF
