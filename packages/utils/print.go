/*

...

*/


package utils


import (
    "fmt"
)


// ============================================================================================================================


// PUBLIC

// https://groups.google.com/forum/#!topic/Golang-nuts/nluStAtr8NA

func PrintError(str string) {
    fmt.Println("\t \x1b[31;1m  *** " + str + " \x1b[0m")
}


func PrintErrorFull(str string, err error) {
    fmt.Println("\t \x1b[31;1m  *** " + str + ": %s \x1b[0m", err)
}


func PrintStatus(str string) {
    fmt.Println("\x1b[32;1m  ~~~ " + str + "\x1b[0m")
}


func PrintSuccess(str string) {
    fmt.Println("\t \x1b[35;1m  \u2713\u2713\u2713 " + str + "\x1b[0m")
}


func PrintLine() {
    fmt.Println("\x1b[33;1m \n  ---------------------------------------------------------------- \x1b[0m")
}


// ============================================================================================================================


// EOF