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

func (c *Crypto) GeneratePassword() ([]byte, error) {

    if c.Verbose {
	    c.Logger("Creating a password...")
    }

    randomBytes, err := c.GenerateRandomBytes(c.PasswordLength)
    if err != nil {
    	if c.Verbose {
	    	c.Logger("Error generating random bytes" + err.Error())
	    }
        return nil, err
    }

    if c.Verbose {
	    c.Logger("Generated password.")
    }

    return randomBytes, nil

} // end of createPassword


// ============================================================================================================================


// EOF