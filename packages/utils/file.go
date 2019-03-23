/*

Utility functions...

*/

package utils

import (
	"io/ioutil"
	"os"
	"strings"
)

// ============================================================================================================================

// PUBLIC

func WriteToFile(data []byte, pathToFile string) error {
	err := ioutil.WriteFile(pathToFile, data, 777)
	return err
}

func StoreFileFromDownload(f File, path string) error {
	// WriteToFile(f.Content, path + f.Name) // permissions...
	err := ioutil.WriteFile(path+f.FileName, f.Content, 0755)
	return err
}

func ReadFromFile(pathToFile string) ([]byte, error) {
	data, err := ioutil.ReadFile(pathToFile)
	return data, err
}

func IsFile(pathToFile string) bool {
	if s, err := os.Stat(pathToFile); os.IsNotExist(err) || s.IsDir() || err != nil {
		return false
	}
	return true
}

func DeleteFile(pathToFile string) error {
	err := os.Remove(pathToFile)
	return err
}

func StripFilePathBase(pathToFile string) string {
	return strings.Replace(pathToFile, "file://", "", -1)
}

// ============================================================================================================================

// EOF
