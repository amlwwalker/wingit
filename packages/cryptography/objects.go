/*

...

*/


package cryptography


import (
    "crypto/rsa"
)


// ============================================================================================================================


type Crypto struct {
    PasswordLength      int // in bytes
    KeyFolder           string
    SyncFolder           string
    Logger func(message string)
    Verbose             bool

    PrivateKey          *rsa.PrivateKey
    PublicKey           *rsa.PublicKey
}


// ============================================================================================================================


// EOF