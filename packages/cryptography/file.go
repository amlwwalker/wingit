/*

Functions for file handling...

*/


package cryptography


import (
    utils "github.com/amlwwalker/wingit/packages/utils"
)


// ============================================================================================================================


// Objects defined in objects.go


// ============================================================================================================================

func (c *Crypto) EncryptFile(pathToFile string, password []byte) ([]byte, error) {

    if c.Verbose {
        c.Logger("Starting the encryption of a file...")
    }

    // Load the file as bytes.
    fileBytes, err := utils.ReadFromFile(pathToFile)
    if err != nil {
        if c.Verbose {
            c.Logger("Error loading file to encrypt" + err.Error())
        }
        return nil, err
    }

    encrypted, err := c.AESEncrypt(fileBytes, password)
    if err != nil {
        if c.Verbose {
            c.Logger("Error ecrypting file contents" + err.Error())
        }
        return nil, err
    }

    if c.Verbose {
        c.Logger("Successfully encrypted the file contents.")
    }

    return encrypted, nil

}

func (c *Crypto) EncryptFileBytes(fileBytes []byte, password []byte) ([]byte, error) {

    if c.Verbose {
        c.Logger("Starting the encryption of a file from byte array...")
    }

    encrypted, err := c.AESEncrypt(fileBytes, password)
    if err != nil {
        c.Logger("Error ecrypting file contents" + err.Error())
        return nil, err
    }

    if c.Verbose {
        c.Logger("Successfully encrypted the file contents.")
    }

    return encrypted, nil

}


func (c *Crypto) DecryptFile(pathToFile string, password []byte) ([]byte, error) {

    if c.Verbose {
        c.Logger("Starting the decryption of a file...")
    }

    fileBytes, err := utils.ReadFromFile(pathToFile)
    if err != nil {
            c.Logger("Error loading file to decrypt" + err.Error())
        return nil, err
    }

    decrypted, err := c.AESDecrypt(fileBytes, password)
    if err != nil {
        c.Logger("Error decrypting file contents" + err.Error())
        return nil, err
    }

    if c.Verbose {
        c.Logger("Successfully decrypted the file contents.")
    }

    return decrypted, nil

}


// ============================================================================================================================


// EOF