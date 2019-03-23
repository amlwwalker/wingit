// /*

// Handler function for the encryption of a file...

// */


package handlers


// import (
//     utils "bitbucket.org/encryptoclient/packages/utils"
//     cryptography "bitbucket.org/encryptoclient/packages/cryptography"
// )


// // ============================================================================================================================


// func TestRSA(CRYPTO *cryptography.Crypto) {

//     utils.PrintLine()
//     utils.PrintStatus("Testing the RSA crypto-functions...")

//     CRYPTO.InitRSAKeyPair() // Will create new keys...
//     CRYPTO.InitRSAKeyPair() // Will load existing keys...

//     _, err := CRYPTO.LoadPublicKey(CRYPTO.KeyFolder, "public.pem")
//     if err != nil {
//         panic(err)
//     } 

//     utils.PrintStatus("Done testing the RSA crypto-functions...\n")

// }


// func TestAES(CRYPTO *cryptography.Crypto) {

//     utils.PrintLine()
//     utils.PrintStatus("Testing the AES crypto-functions...")

//     pathToFilePlain := "tmp/lipsum.txt"
//     pathToFileEnc := "tmp/lipsum.cry"

//     CRYPTO.InitRSAKeyPair()

//     // Generate a password
//     passwordBytes, err := CRYPTO.GeneratePassword()
//     if err != nil {
//         utils.PrintErrorFull("Error generating password", err)
//         panic(err)
//     }

//     // For testing purposes, make sure lipsum.enc is deleted
//     _ = utils.DeleteFile(pathToFileEnc)

//     // Load and encrypt the file, save as enc.
//     encryptedBytes, err := CRYPTO.EncryptFile(pathToFilePlain, passwordBytes)
//     if err != nil {
//         panic(err)
//     }

//     err = utils.WriteToFile(encryptedBytes, pathToFileEnc)
//     if err != nil {
//         panic(err)
//     }

//     // Load and decrypt the enc. file
//     decryptedBytes, err := CRYPTO.DecryptFile(pathToFileEnc, passwordBytes)
//     if err != nil {
//         panic(err)
//     }

//     // Check the two files are the same...
//     plainBytes, err := utils.ReadFromFile(pathToFilePlain)
//     if string(plainBytes) == string(decryptedBytes) {
//         utils.PrintSuccess("SUCCESS! Decrypted version and plain-text match\n")
//     } else {
//         utils.PrintError("ERROR! Decrypted version and plain-text do not match\n")
//     }

//     // Lets check the signature functions.
//     utils.PrintStatus("Trying the generate and verify the Signature for the encrypted bytes...")
//     signatureBytes, err := CRYPTO.CreateSignature(encryptedBytes)
//     if err != nil {
//         panic(err)
//     }
//     err = CRYPTO.VerifySignature(encryptedBytes, signatureBytes, CRYPTO.PublicKey)
//     if err != nil {
//         panic(err)
//     }

//     // Lets check the HMAC function. Will do this on the encrypted bytes so Encrypt-then-Mac
//     utils.PrintStatus("Trying the generate and verify the HMAC for the encrypted bytes...")
//     hmacBytes := CRYPTO.CreateHMAC(encryptedBytes, passwordBytes)
//     if CRYPTO.VerifyHMAC(encryptedBytes, hmacBytes, passwordBytes) {
//         utils.PrintSuccess("Successfully passed the HMAC test.")
//     } else {
//         utils.PrintError("Error: Failed on the HMAC test.")
//     }

//     utils.PrintStatus("Done testing the AES crypto-functions.\n")
    
// }


// // ============================================================================================================================


// // EOF 