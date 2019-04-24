/*

Handler function for the encryption and uploading of a file...

*/

package handlers

import (
	"fmt"
	"path/filepath"

	cryptography "github.com/amlwwalker/wingit/packages/cryptography"
	srv "github.com/amlwwalker/wingit/packages/server"
	utils "github.com/amlwwalker/wingit/packages/utils"
)

// ============================================================================================================================

func ProcessFileEncryptAndUpload(filePath string, toUserId string, idToken string, CRYPTO *cryptography.Crypto, SERVER *srv.Server) (string, error) {

	// Get the file-name and file-bytes from the path.
	fileName := filepath.Base(filePath)
	filePath = utils.StripFilePathBase(filePath)
	SERVER.Logger("fileName: " + fileName)
	SERVER.Logger("filePath: " + filePath)
	//     fName := path.Base(filePath)
	// extName := path.Ext(inFName)
	// bName := fName[:len(fName)-len(extName)]

	fileBytes, err := utils.ReadFromFile(filePath)
	if err != nil {
		SERVER.Logger("Error in Upload, Process getting the file-bytes" + err.Error())
	}
	SERVER.Logger("about to upload for " + toUserId)
	url, err := UploadFileFromBytes(fileName, fileBytes, toUserId, idToken, CRYPTO, SERVER)
	if err != nil {
		SERVER.Logger("Error in Upload, Process" + err.Error())
		return "", err
	}

	return url, nil

}

func UploadFileFromBytes(fileName string, fileBytes []byte, toUserId string, idToken string, CRYPTO *cryptography.Crypto, SERVER *srv.Server) (string, error) {

	// Setup.
	fileNameBytes := []byte(fileName)

	// Generate password
	passwordBytes, err := CRYPTO.GeneratePassword()
	if err != nil {
		utils.PrintErrorFull("Error generating a password", err)
		return "", err
	}

	// Encrypt file with password
	encryptedBytes, err := CRYPTO.EncryptFileBytes(fileBytes, passwordBytes)
	if err != nil {
		utils.PrintErrorFull("Error encrypting the file", err)
		return "", err
	}
	// encryptedBytes = encryptedBytes
	// fmt.Println("WARNING - BYTES NOT BEING ENCRYPTED AT THIS STAGE")
	var url string
	if url, err = UploadEncryptedFile(fileNameBytes, passwordBytes, encryptedBytes, toUserId, idToken, CRYPTO, SERVER); err != nil {
		utils.PrintErrorFull("Error uploading the file", err)
		return "", err
	}

	return url, nil

}

// func GenerateEncryptedBytesFromFile(filePath string, fileName string, toUserId string, idToken string, CRYPTO *cryptography.Crypto, SERVER *srv.Server) (error) {

//     // Setup.
//     pathToFilePlain := filePath + fileName
//     fileNameBytes := []byte(fileName)

//     // Generate password
//     passwordBytes, err := CRYPTO.GeneratePassword()
//     if err != nil {
//         utils.PrintErrorFull("Error generating a password", err)
//         return err
//     }

//     // Encrypt file with password
//     encryptedBytes, err := CRYPTO.EncryptFile(pathToFilePlain, passwordBytes)
//     if err != nil {
//         utils.PrintErrorFull("Error encrypting the file", err)
//         return err
//     }
//     if err := UploadEncryptedFile(fileNameBytes, passwordBytes, encryptedBytes, toUserId, idToken, CRYPTO, SERVER); err != nil {
//         utils.PrintErrorFull("Error uploading the file", err)
//     }

//     return nil

// }

func UploadEncryptedFile(fileNameBytes []byte, passwordBytes []byte, encryptedBytes []byte, toUserId string, idToken string, CRYPTO *cryptography.Crypto, SERVER *srv.Server) (string, error) {

	// Encrypt filename with public key
	// fileNameEncryptedBytes, err := CRYPTO.EncryptRSA(CRYPTO.PublicKey, fileNameBytes)
	fileNameEncryptedBytes, err := CRYPTO.AESEncrypt(fileNameBytes, passwordBytes)
	if err != nil {
		utils.PrintErrorFull("Error encrypting the file name", err)
		return "", err
	}
	if !utils.IsFile(CRYPTO.KeyFolder + toUserId + ".pem") {
		downloadedKey, err := SERVER.GetKey(toUserId, idToken)
		if err != nil {
			SERVER.Logger("Error fetching the required public-key from the server" + err.Error())
			return "", err
		}
		CRYPTO.SavePublicKeys([]utils.KeyServer{downloadedKey})
	}
	// Load recipients public key...
	fmt.Println("loading senders public key " + toUserId + ".pem")
	recipientPublicKey, err := CRYPTO.LoadPublicKey(CRYPTO.KeyFolder, toUserId+".pem")
	if err != nil {
		utils.PrintErrorFull("Error loading the sender's public key", err)
		return "", err
	}

	// Encrypt password with recipients public key
	passwordEncryptedBytes, err := CRYPTO.EncryptRSA(recipientPublicKey, passwordBytes)
	if err != nil {
		utils.PrintErrorFull("Error encrypting the password", err)
		return "", err
	}

	// Sign file using private key
	signatureBytes, err := CRYPTO.CreateSignature(encryptedBytes)
	if err != nil {
		utils.PrintErrorFull("Error signing the encrypted contents", err)
		return "", err
	}

	// test signature
	fmt.Println("testing signature")
	err = CRYPTO.VerifySignature(encryptedBytes, signatureBytes, CRYPTO.PublicKey)
	if err != nil {
		SERVER.Logger("Error verifying the signature of the sender" + err.Error())
	}

	// HMAC file
	hmacBytes := CRYPTO.CreateHMAC(encryptedBytes, passwordBytes)

	// Create JSON Payload to send to server...
	var payload utils.File
	payload.ContentEnc = encryptedBytes
	payload.FileNameEnc = utils.EncodeBase64(fileNameEncryptedBytes)
	payload.PasswordEnc = utils.EncodeBase64(passwordEncryptedBytes)
	payload.Signature = utils.EncodeBase64(signatureBytes)
	payload.HMAC = utils.EncodeBase64(hmacBytes)
	payload.UserID = utils.EncodeBase64([]byte(toUserId))

	fmt.Println("payload is for ", payload.UserID, " signature ", payload.Signature)

	// payloadJSONIndented, _ := json.MarshalIndent(payload, "", "  ")
	// fmt.Println("PAYLOAD FOR ENCRYPTION ", string(payloadJSONIndented))

	// make post request... --- fallback to temporary storage?

	status, url, err := SERVER.PostFile(payload, idToken)
	if err != nil {
		utils.PrintErrorFull("Error posting file", err)
		return "", err
	}
	//just don't thnk its necessary to store it locally for future. Just do it again
	// 	// Try alternatively to temporary store the file for later sync.
	// 	err = StoreEncryptedFileTemporary(&payload, CRYPTO)
	// 	if err != nil {
	// 		utils.PrintErrorFull("Error storing temporary", err)
	// 		return "", err
	// 	}
	if status != "201 OK" {
		fmt.Println("File didn't upload correctly for some reason. Check Status ")
		return status, nil
	}

	return url, nil

}

// ============================================================================================================================

// EOF
