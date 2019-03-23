/*

Functions for AES Encryption and Decryption

*/


package cryptography


import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
    "errors"

)


// ============================================================================================================================


// Objects defined in objects.go


// ============================================================================================================================


func (c *Crypto) AESEncrypt(plaintext []byte, password []byte) ([]byte, error) {

    if c.Verbose {
        c.Logger("Encrypting AES Encrypted payload...")
    }

    key := []byte(password) // Why?

    // Create the AES cipher
    block, err := aes.NewCipher(key)
    if err != nil {
        c.Logger("Error creating an AES Cipher for encryption" + err.Error())
        return nil, err
    }

    // Empty array of 16 + plaintext length
    // Include the IV at the beginning
    ciphertext := make([]byte, aes.BlockSize + len(plaintext))

    iv := ciphertext[:aes.BlockSize] // Slice of first 16 bytes

    // Write 16 rand bytes to fill iv
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        c.Logger("Error creating random IV for AES Encryption" + err.Error())
        return nil, err
    }

    // Return an encrypted stream
    stream := cipher.NewCFBEncrypter(block, iv)

    // Encrypt bytes from plaintext to ciphertext
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    if c.Verbose {
        c.Logger("Successfully encrypted the payload.")
    }

    return ciphertext, nil

} // end of AESEncrypt


func (c *Crypto) AESDecrypt(ciphertext []byte, password []byte) ([]byte, error) {

    if c.Verbose {
        c.Logger("Decrypting AES Encrypted payload...")
    }

    key := []byte(password) // Why?

    // Create the AES cipher
    block, err := aes.NewCipher(key)
    if err != nil {
        c.Logger("Error creating an AES Cipher for decryption" + err.Error())
        return nil, err
    }

    // Before even testing the decryption,
    // if the text is too small, then it is incorrect
    if len(ciphertext) < aes.BlockSize {
        err = errors.New("Length of decryption block does not correspond with aes-block-size")
        return nil, err
    }

    iv := ciphertext[:aes.BlockSize] // Get the 16 byte IV
    ciphertext = ciphertext[aes.BlockSize:] // Remove the IV from the ciphertext

    // Return a decrypted stream
    stream := cipher.NewCFBDecrypter(block, iv)

    // Decrypt bytes from ciphertext
    stream.XORKeyStream(ciphertext, ciphertext)

    if c.Verbose {
        c.Logger("Successfully decrypted the payload.")
    }

    return ciphertext, nil

} // end of AESDecrypt


// ============================================================================================================================


// EOF