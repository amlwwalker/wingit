/*

All cryptographic functions related to hashing

*/


package cryptography


import (
    "hash"
    "crypto"
    "crypto/sha256"
    "crypto/hmac"

)


// ============================================================================================================================


// Objects defined in objects.go


// ============================================================================================================================


// Our default hashing function.
func (c *Crypto) Hash(contentBytes []byte) ([]byte, crypto.Hash) {

    if c.Verbose {
        c.Logger("Hashing bytes...")
    }

    var sha256Hash hash.Hash

    // Calculate the hash for the content bytes
    sha256Hash = sha256.New()
    sha256Hash.Write(contentBytes)
    hashedFileBytes := sha256Hash.Sum(nil)

    if c.Verbose {
        c.Logger("Successfully hashed bytes.")
    }

    return hashedFileBytes, crypto.SHA256

} // end of Hash


// Create an HMAC
// https://golang.org/pkg/crypto/hmac/
func (c *Crypto) CreateHMAC(messageBytes []byte, password []byte) ([]byte) {

    if c.Verbose {
        c.Logger("Hashing bytes...")
    }

    mac := hmac.New(sha256.New, password)
    mac.Write(messageBytes)
    hmacBytes := mac.Sum(nil)

    if c.Verbose {
        c.Logger("Successfully hashed bytes.")
    }

    return hmacBytes

}


func (c *Crypto) VerifyHMAC(messageBytes []byte, hmacBytes []byte, password []byte) bool {

    mac := hmac.New(sha256.New, password)
    mac.Write(messageBytes)
    expectedBytes := mac.Sum(nil)

    return hmac.Equal(expectedBytes, hmacBytes)

}


// ============================================================================================================================


// EOF