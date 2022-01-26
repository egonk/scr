//go:build go1.18
// +build go1.18

package scr

// Must panics if err is not nil. v is returned if err is nil.
//  f := scr.Must(os.Create("example"))
//  defer scr.Close(f)
//  // write to f
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
