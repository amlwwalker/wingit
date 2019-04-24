/*

This holds all the RSA functions

*/

package cryptography

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"

	// "fmt"
	"encoding/json"
	"errors"
	"hash"

	utils "github.com/amlwwalker/wingit/packages/utils"
)

// ============================================================================================================================

// Objects defined in objects.go

// ============================================================================================================================

// PUBLIC FUNCTIONS

func (c *Crypto) InitRSAKeyPair(id string) {

	var err error

	if c.Verbose {
		c.Logger("Looking for " + c.KeyFolder + id + "-private.key")
	}
	// Check if the current files (public.pem, private.key) exist.
	// If they do, load them. Otherwise create new ones.

	if utils.IsFile(c.KeyFolder+id+"-private.key") == true {

		err := c.loadRSAKeyPair(id) // Load the current key-pair
		if err != nil {
			panic(err)
		}

	} else {

		err = c.generateRSAKeyPair(id) // Create a new key-pair
		if err != nil {
			panic(err)
		}

		// Key-pair has been generated. Save the public-key to be synchronized.

		// Load public-key bytes.
		pubKeyBytes, err := c.GetPublicKeyBytes()
		if err != nil {
			panic(err)
		}

		// Save in the same structure as needed for uploading later.
		var key utils.Key
		key.Content = pubKeyBytes

		// Store the file in the temporary folder.
		keyJSONIndented, err := json.MarshalIndent(key, "", "  ")
		if err != nil {
			c.Logger("Error (temporary) storing the public-key file locally" + err.Error())
			panic(err)
		}

		err = utils.WriteToFile(keyJSONIndented, c.SyncFolder+"public-key-to-sync.pem")
		if err != nil {
			c.Logger("Fatal Error (temporary) storing the public-key file locally" + err.Error())
			panic(err)
		}
		if c.Verbose {
			utils.PrintSuccess("Successfully stored the public-key file locally as public-key-to-sync.pem")
		}
	}

	err = c.testRSAKeyPair()
	if err != nil {
		if c.Verbose {
			c.Logger("Error during the encryption (RSA-OAEP) test" + err.Error())
		}
	}

} // end of InitRSAKeyPair

// Wrapper function for EncryptRSAOaep
func (c *Crypto) EncryptRSA(publicKey *rsa.PublicKey, plaintext []byte) ([]byte, error) {

	var label []byte

	return c.EncryptRSAOaep(publicKey, plaintext, label)

} // end of Encrypt

// Wrapper function for DecryptRSAOaep
func (c *Crypto) DecryptRSA(ciphertext []byte) ([]byte, error) {

	var label []byte

	return c.DecryptRSAOaep(ciphertext, label)

} // end of Decrypt

/*
OAEP: Optimal Asymmetric Encryption Padding

https://en.wikipedia.org/wiki/Optimal_asymmetric_encryption_padding
Add an element of randomness which can be used to convert a deterministic encryption scheme
(e.g., traditional RSA) into a probabilistic scheme.
*/
func (c *Crypto) EncryptRSAOaep(publicKey *rsa.PublicKey, plaintext []byte, label []byte) ([]byte, error) {

	var err error
	var encrypted []byte
	var md5Hash hash.Hash

	md5Hash = md5.New()

	encrypted, err = rsa.EncryptOAEP(md5Hash, rand.Reader, publicKey, plaintext, label)
	if err != nil {
		if c.Verbose {
			c.Logger("Error RSA OAEP Encrypting" + err.Error())
		}
		return nil, err
	}

	return encrypted, nil

} // end of EncryptRSAOaep

func (c *Crypto) DecryptRSAOaep(ciphertext []byte, label []byte) ([]byte, error) {

	var err error
	var decrypted []byte
	var md5_hash hash.Hash

	md5_hash = md5.New()
	decrypted, err = rsa.DecryptOAEP(md5_hash, rand.Reader, c.PrivateKey, ciphertext, label)
	if err != nil {
		return nil, err
	}

	return decrypted, nil

} // end of DecryptRSAOaep

// Function to sign bytes using the RSA Private Key.
// https://golang.org/pkg/crypto/rsa/#SignPSS
func (c *Crypto) CreateSignature(contentBytes []byte) ([]byte, error) {

	if c.Verbose {
		c.Logger("Creating RSA Signature for bytes...")
	}

	hashedFileBytes, hashFunction := c.Hash(contentBytes)

	// Sign the contentBytes.
	signatureBytes, err := rsa.SignPSS(rand.Reader, c.PrivateKey, hashFunction, hashedFileBytes, nil)
	if err != nil {
		if c.Verbose {
			c.Logger("Error creating RSA Signature" + err.Error())
		}
		return nil, err
	}
	// fmt.Println("FINAL SIGNATURE: ", signatureBytes, " HASH BYTES ", hashedFileBytes)
	if c.Verbose {
		c.Logger("Successfully signed the bytes.")
	}

	return signatureBytes, nil

} // end of CreateSignature

// https://golang.org/pkg/crypto/rsa/#VerifyPSS
func (c *Crypto) VerifySignature(contentBytes []byte, signatureBytes []byte, publicKey *rsa.PublicKey) error {

	if c.Verbose {
		c.Logger("Verifying RSA Signature for bytes...")
	}

	hashedFileBytes, hashFunction := c.Hash(contentBytes)
	// fmt.Println("VERIFYING SIGNATURE: ", signatureBytes, " HASH BYTES ", hashedFileBytes)
	// Verify the signature
	//func VerifyPSS(pub *PublicKey, hash crypto.Hash, hashed []byte, sig []byte, opts *PSSOptions) error
	err := rsa.VerifyPSS(publicKey, hashFunction, hashedFileBytes, signatureBytes, nil)
	if err != nil {
		if c.Verbose {
			c.Logger("Error verifying RSA Signature" + err.Error())
		}
		return err
	}

	if c.Verbose {
		c.Logger("Successfully verified the signature.")
	}

	return nil

} // end of VerifySignature

// ============================================================================================================================

// PRIVATE FUNCTIONS

// Function to generate a new RSA public/private keypair
func (c *Crypto) generateRSAKeyPair(id string) error {

	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey
	var err error

	if c.Verbose {
		c.Logger("RSA: Generating a public/private keypair...")
	}

	// Generate Private Key
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		if c.Verbose {
			c.Logger("Error generating RSA key-pair" + err.Error())
		}
		return err
	}

	// Precompute some calculations
	// Calculations that speed up private key operations in the future
	privateKey.Precompute()

	// Validate the private-key -- Sanity checks on the key
	if err = privateKey.Validate(); err != nil {
		if c.Verbose {
			c.Logger("Error validating RSA public-key" + err.Error())
		}
		return err
	}

	// Obtain the public-key from the private-key
	publicKey = &privateKey.PublicKey

	// Store in the struct.
	c.PrivateKey = privateKey
	c.PublicKey = publicKey

	// Save the private-key to a .key file
	if err = c.saveRSAPrivateKey(id); err != nil {
		if c.Verbose {
			c.Logger("Error saving RSA private key" + err.Error())
		}
		return err
	}
	// Save the public-key part to a .pem file.
	if err = c.saveRSAPublicKey(c.KeyFolder, id, publicKey); err != nil {
		if c.Verbose {
			c.Logger("Error saving public-key" + err.Error())
		}
		return err
	}

	if c.Verbose {
		c.Logger("Successfully generated a new RSA key-pair.")
	}

	return nil

} // end of generateRSAKeyPair

// Test the RSA key-pair.
// We do this by encrypting a string, decrypting that string and checking if they match.
func (c *Crypto) testRSAKeyPair() error {

	var err error
	var label, plaintext, ciphertext, decrypted []byte

	if c.Verbose {
		c.Logger("Testing the key-pair for encryption (RSA-OAEP)...")
	}

	// Create a test text-sample.
	plaintext = []byte("The quick brown fox jumps over the zebra!")

	// Encryption & Decryption
	ciphertext, err = c.EncryptRSAOaep(c.PublicKey, plaintext, label)
	if err != nil {
		return err
	}
	decrypted, err = c.DecryptRSAOaep(ciphertext, label)
	if err != nil {
		return err
	}

	// Assert & Print for debugging.
	if string(plaintext) != string(decrypted) {
		// fmt.Printf("\tOAEP Encrypted [%s] to \n[%x]\n", string(plaintext), ciphertext)
		// fmt.Printf("\tOAEP Decrypted [%x] to \n[%s]\n", ciphertext, decrypted)
		// Create error and return
		err = errors.New("Error in testing the RSA keypair encryption. Encryption failed!")
		return err
	}

	if c.Verbose {
		c.Logger("Successfully passed encryption (RSA-OAEP) test.")
	}

	return nil

} // end of testRSAKeyPair

// ============================================================================================================================

// EOF
