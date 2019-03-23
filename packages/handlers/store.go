/*

Handler function for the encryption and uploading of a file...

*/


package handlers


import (
    "encoding/json"

    cryptography "github.com/amlwwalker/wingit/packages/cryptography"
    utils "github.com/amlwwalker/wingit/packages/utils"
)


// ============================================================================================================================


func StoreEncryptedFileTemporary(payload *utils.File, CRYPTO *cryptography.Crypto) (error) {

    utils.PrintStatus("Attempting local storage of encrypted file...")

    // Store the file in the temporary folder.
    payloadJSONIndented, err := json.MarshalIndent(payload, "", "  ")
    if err != nil {
        utils.PrintErrorFull("Error (temporary) storing the encrypted file locally", err)
        return err
    }

    // Use a hashed name to make sure there is no leakage.
    nameHashBytes, _ := CRYPTO.Hash([]byte(payload.FileNameEnc))

    err = utils.WriteToFile(payloadJSONIndented, "tmp/" + string(nameHashBytes)[0:8] + ".enc")
    if err != nil {
        utils.PrintErrorFull("Error (temporary) storing the encrypted file locally", err)
        return err
    }

    utils.PrintSuccess("Successfully stored the encrypted file locally as " + string(nameHashBytes)[0:8] + ".enc")

    return nil

}


// ============================================================================================================================


// EOF