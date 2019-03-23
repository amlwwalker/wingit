/*

Utility functions...

*/


package utils


import (
    "os"
)


// ============================================================================================================================

func InitiateDirectory(directory string) {
    // For the keys-folder we need to check if the folder exists...
    checkDir, err := IsDirectory(directory)
    if err != nil {
        PrintErrorFull("Error checking for " + directory + " directory", err)
        panic(err)
    }

    if checkDir == true {
        PrintStatus(directory + " already exists.")
    } else {
        // Create the directory.
        PrintStatus("Creating " + directory + "...")
        err = CreateDirectory(directory)
        if err != nil {
            PrintErrorFull("Error creating the folder.", err)
        } else {
            PrintSuccess("folder created.")
        }
    }
}

// PUBLIC

func IsDirectory(path string) (bool, error) {

    s, err := os.Stat(path) // returns an error if the path does not exist.
    if err != nil {
        if os.IsNotExist(err) {
            return false, nil
        } else {
            return false, err // Different error...?
        }
    }

    if s.IsDir() {
        return true, nil
    }

    return false, nil // Redundancy

}


func CreateDirectory(path string) (error) {

    // Assumes checks have been done on if the directory exists...
    err := os.MkdirAll(path, os.ModePerm)
    if err != nil {
        return err
    }

    return nil // Redundancy

}


func DeleteDirectory(path string) (error) {

    err := os.RemoveAll(path)
    return err

}


// ============================================================================================================================


// EOF