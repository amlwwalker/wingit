/*

...

*/


package cryptography


import (
)


// ============================================================================================================================


// Objects defined in objects.go


// ============================================================================================================================


// PUBLIC

func (c *Crypto) Init(passwordLength int, keysFolder, syncFolder string, verbosity bool) {

    // Set the verbosity.
    c.Verbose = verbosity

    c.PasswordLength = passwordLength

    // Set the keys folder.
    c.KeyFolder = keysFolder
    c.SyncFolder = syncFolder

    if c.Verbose {
        c.Logger("Initialising crypto handlers...")
    }
}




// ============================================================================================================================


// PRIVATE

// ...


// ============================================================================================================================


// EOF