/*

...

*/


package cryptography


import (
	"crypto/rand"

)


// ============================================================================================================================


// Objects defined in objects.go


// ============================================================================================================================


func (c *Crypto) GenerateRandomBytes(n int) ([]byte, error) {

    // Create the byte-array.
    b := make([]byte, n)

    // Generate the random bytes.
    _, err := rand.Read(b)
    if err != nil {
        if c.Verbose {
    	   c.Logger("Error generating random bytes" + err.Error())
        }
        return nil, err
    }

    return b, nil

} // end of GenerateRandomBytes


// ============================================================================================================================


// EOF