/*

Objects available.

*/

package server

import (
	"fmt"
	"io"
)

// ============================================================================================================================

// ...
type Server struct {
	Ip               string
	Address          string
	Port             string
	Verbose          bool
	Logger           func(message string)
	DownloadProgress func(progress float64, sync bool, err error)
	UploadProgress   func(progress float64, err error)
	DownloadFolder   string
	SyncFolder       string
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
// PassThru wraps an existing io.Reader.
//
// It simply forwards the Read() call, while displaying
// the results from individual calls to it.
type PassThru struct {
	io.Reader
	total            int64 // Total # of bytes transferred
	length           int64 // Expected length
	progress         float64
	Sync             bool
	DownloadProgress func(progress float64, sync bool, err error)
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy(). We simply
// use it to keep track of byte counts and then forward the call.
func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if n > 0 {
		pt.total += int64(n)
		percentage := float64(pt.total) / float64(pt.length)
		fmt.Println(percentage)
		if percentage-pt.progress > 0.02 {
			if !pt.Sync {
				go pt.DownloadProgress(percentage, pt.Sync, nil) //can we just pass it through
			}
			pt.progress = percentage
		}
	}

	return n, err
}
