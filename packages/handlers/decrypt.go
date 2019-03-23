/*

Handler function for the downloading and decryption of a file...

*/

package handlers

import (
	"errors"
	"fmt"
	"reflect"

	cryptography "github.com/amlwwalker/wingit/packages/cryptography"
	srv "github.com/amlwwalker/wingit/packages/server"
	utils "github.com/amlwwalker/wingit/packages/utils"
)

// ============================================================================================================================

func DownloadAndDecryptFile(fileNameEnc string, pendingContact string, idToken string, CRYPTO *cryptography.Crypto, SERVER *srv.Server) error {

	// Download the file
	fileEncrypted, err := SERVER.GetFile(fileNameEnc, idToken)
	if err != nil {
		SERVER.Logger("Error downloading the file" + err.Error())
	}
	if reflect.DeepEqual(fileEncrypted, []utils.File{}) {
		fmt.Println("no file was returned. Search results came back empty")
		return errors.New("No file was returned with this name for this id")
	}
	// fmt.Printf("Was sent file %+v\r\n", fileEncrypted)

	// Verify the signature (requires sender's public-key)
	// Check if that public key exists now... If not, download.
	// Verify the signature.
	if !utils.IsFile(CRYPTO.KeyFolder + pendingContact + ".pem") {
		downloadedKeys, err := SERVER.GetKey(pendingContact, idToken)
		if err != nil {
			SERVER.Logger("Error fetching the required public-key from the server" + err.Error())
			return err
		}
		// fmt.Println("saving key for " + pendingContact + idToken)
		CRYPTO.SavePublicKeys(downloadedKeys)
	}
	fmt.Println("verifying signature with key for ", pendingContact)
	// The key should now exist. Load the key.
	senderPublicKey, err := CRYPTO.LoadPublicKey(CRYPTO.KeyFolder, pendingContact+".pem")
	if err != nil {
		SERVER.Logger("1 Error loading the sender's public key" + err.Error())
		return err
	}
	// Verify the signature.
	signatureBytes, err := utils.DecodeBase64(fileEncrypted.Signature)
	if err != nil {
		SERVER.Logger("1 Error decoding the signature from base64" + err.Error())
		return err
	}

	err = CRYPTO.VerifySignature(fileEncrypted.ContentEnc, signatureBytes, senderPublicKey)
	if err != nil {
		SERVER.Logger("1 Error verifying the signature of the sender" + err.Error())
		// return err
	}
	// Get the plain-text password. Decrypt using private-key.
	passwordEncBytes, err := utils.DecodeBase64(fileEncrypted.PasswordEnc)
	if err != nil {
		SERVER.Logger("Error decoding the password from base64" + err.Error())
		return err
	}

	passwordBytes, err := CRYPTO.DecryptRSA(passwordEncBytes)
	if err != nil {
		SERVER.Logger("Error decrypting the password " + err.Error())
		return err
	}
	// fmt.Println("password bytes " + string(passwordBytes))
	// Verify the HMAC (requires password)
	hmacBytes, err := utils.DecodeBase64(fileEncrypted.HMAC)
	if err != nil {
		SERVER.Logger("Error decoding the HMAC from base64" + err.Error())
		return err
	}
	checkHmac := CRYPTO.VerifyHMAC(fileEncrypted.ContentEnc, hmacBytes, passwordBytes)
	if !checkHmac {
		SERVER.Logger("Error: HMAC can not be verified!")
		// panic("ERROR: HMAC EN DECRYPT")
	}

	// Decrypt the filename (requires password).
	fileNameEncBytes, err := utils.DecodeBase64(fileEncrypted.FileNameEnc)
	if err != nil {
		SERVER.Logger("Error decoding the file name from base64" + err.Error())
		return err
	}
	fileNameBytes, err := CRYPTO.AESDecrypt(fileNameEncBytes, passwordBytes)
	if err != nil {
		SERVER.Logger("Error decrypting the file name" + err.Error())
		return err
	}

	// Decrypt the contents (requires password)
	contentBytes, err := CRYPTO.AESDecrypt(fileEncrypted.ContentEnc, passwordBytes)
	if err != nil {
		SERVER.Logger("Error decrypting the file contents" + err.Error())
		return err
	}
	// fmt.Printf("FileEncrypted %+v\r\n", fileEncrypted)
	// fmt.Println("encrData ", fileEncrypted.ContentEnc, "contentBytes ", contentBytes)
	var newFile utils.File
	newFile.FileName = string(fileNameBytes)
	newFile.Content = contentBytes
	// fmt.Printf("THE FILE STRAIGHT FROM SERVER AS STRING FOR %+v\r\n", newFile)
	utils.StoreFileFromDownload(newFile, SERVER.DownloadFolder)

	return nil

}

func DecryptFileName(nameEncBase64 string, pwEncBase64 string, CRYPTO *cryptography.Crypto) (string, error) {

	// Decode both.
	passwordEncBytes, err := utils.DecodeBase64(pwEncBase64)
	if err != nil {
		CRYPTO.Logger("Error decoding the password from base64" + err.Error())
	}
	fileNameEncBytes, err := utils.DecodeBase64(nameEncBase64)
	if err != nil {
		CRYPTO.Logger("Error decoding the file name from base64" + err.Error())
	}

	// fmt.Println("pw64 ", pwEncBase64, " pwBytes", string(passwordEncBytes))
	// Decrypt both.
	// fmt.Println("password to decript pre download ", string(passwordEncBytes))
	passwordBytes, err := CRYPTO.DecryptRSA(passwordEncBytes)
	if err != nil {
		CRYPTO.Logger("Error decrypting the password for the filename" + err.Error())
	}
	fileNameBytes, err := CRYPTO.AESDecrypt(fileNameEncBytes, passwordBytes)
	if err != nil {
		CRYPTO.Logger("Error decrypting the file name" + err.Error())
	}

	return string(fileNameBytes), err

}

// ============================================================================================================================

// EOF
