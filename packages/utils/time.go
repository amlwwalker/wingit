/*

...

*/


package utils


import (
    "time"
)


// ============================================================================================================================


// PUBLIC

func DelaySecond(n time.Duration) {
    time.Sleep(n * time.Second)
}


// ============================================================================================================================


// EOF