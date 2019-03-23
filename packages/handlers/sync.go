/*

Handler function for the encryption and uploading of a file...

*/


package handlers


import (
    "strings"
    "os"
    "path/filepath"
    "fmt"
    srv "github.com/amlwwalker/wingit/packages/server"
    utils "github.com/amlwwalker/wingit/packages/utils"
)


// ============================================================================================================================


func SyncTemporaryStorageFolder(SERVER *srv.Server, idToken string, modify bool) (error) {

    // utils.PrintStatus("Attempting to sync all locally stored files...")

    // Check if there is internet-connectivity and the server responds?


    // Load all files from the temp-folder with the *.enc extension.
    filepath.Walk( SERVER.SyncFolder, func(path string, info os.FileInfo, err error) (error) {

        // Make sure this is not a folder.
        if !info.Mode().IsDir() {
            fmt.Println("attemping to sync " + path)
            // Check if this is a temporary stored encrypted file.
            if strings.Contains(path, ".enc") {
                err := SERVER.SyncLocalFile(path, idToken)
                if err != nil {
                    utils.PrintError("Aborting file sync...")
                }
            }

            // Check if this is a public-key that was temporarily stored.
            if strings.Contains(path, ".pem") {
                err := SERVER.SyncLocalKey(path, idToken, modify)
                if err != nil {
                    utils.PrintError("Aborting public-key sync...")
                }
            }

        }

        return nil

    }) // end of filepath.Walk

    return nil

}





// ============================================================================================================================


// EOF