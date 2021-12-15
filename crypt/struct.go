package crypt

import "context"

// RangeSeeker is the interface that wraps the RangeSeek method.
//
// Some of the returns from Object.Open() may optionally implement
// this method for efficiency purposes.
type RangeSeeker interface {
	// RangeSeek behaves like a call to Seek(offset int64, whence
	// int) with the output wrapped in an io.LimitedReader
	// limiting the total length to limit.
	//
	// RangeSeek with a limit of < 0 is equivalent to a regular Seek.
	RangeSeek(ctx context.Context, offset int64, whence int, length int64) (int64, error)
}
