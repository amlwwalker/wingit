/*

Base64 Encoding and Decoding functions...

*/


package utils


import (
    "encoding/base64"
)


// ============================================================================================================================


// PUBLIC

// Wrapper for EncodeBase64FromBytes
func EncodeBase64(payloadBytes []byte) (string) {
    return EncodeBase64FromBytes(payloadBytes)
} // end of EncodeBase64


func EncodeBase64FromBytes(payloadBytes []byte) (string) {
    return base64.URLEncoding.EncodeToString(payloadBytes)
} // end of EncodeBase64FromBytes


// Wrapper for DecodeBase64ToBytes
func DecodeBase64(payloadString string) ([]byte, error) {
    return DecodeBase64ToBytes(payloadString)
} // end of DecodeBase64


func DecodeBase64ToBytes(payloadString string) ([]byte, error) {
    return base64.URLEncoding.DecodeString(payloadString)
} // end of DecodeBase64ToBytes


// Convenience function.
func EncodeBase64FromString(payloadString string) (string) {
    return base64.URLEncoding.EncodeToString([]byte(payloadString))
} // end of EncodeBase64FromString


// Convenience function.
func DecodeBase64ToString(payloadString string) (string, error) {
	dec, err := base64.URLEncoding.DecodeString(payloadString)
    return string(dec), err
} // end of DecodeBase64ToString


// ============================================================================================================================


// EOF