package scr_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/egonk/scr"
)

func expectRecover(t *testing.T, expect interface{}) {
	if r := recover(); !reflect.DeepEqual(r, expect) {
		t.Fatal(r)
	}
}

func expectRecoverErr(t *testing.T, expect string) {
	r := recover()
	if err, _ := r.(error); err == nil || err.Error() != expect {
		t.Fatal(r)
	}
}

type closer struct {
	err error
}

func (c closer) Close() error { return c.err }

func setInput(data []byte, fn func()) {
	tmp, err := ioutil.TempFile("", "")
	scr.Err(err)
	_, err = tmp.Write(data)
	scr.Err(err)
	_, err = tmp.Seek(0, 0)
	scr.Err(err)
	defer func() {
		scr.Err(os.Remove(tmp.Name()))
	}()
	defer scr.Close(tmp)
	orig := os.Stdin
	defer func() { os.Stdin = orig }()
	os.Stdin = tmp
	fn()
}

func captureOutput(f **os.File, fn func()) (output []byte) {
	tmp, err := ioutil.TempFile("", "")
	scr.Err(err)
	defer func() {
		output, err = ioutil.ReadFile(tmp.Name())
		scr.Err(err)
		scr.Err(os.Remove(tmp.Name()))
	}()
	defer scr.Close(tmp)
	orig := *f
	defer func() { *f = orig }()
	*f = tmp
	fn()
	return
}

func TestClose(t *testing.T) {
	defer expectRecover(t, errors.New("err"))
	scr.Close(closer{errors.New("err")})
}

func TestClose_nil(t *testing.T) {
	scr.Close(closer{})
}

func TestErr(t *testing.T) {
	defer expectRecover(t, errors.New("err"))
	scr.Err(errors.New("err"))
}

func TestErr_nil(t *testing.T) {
	scr.Err(nil)
}

func testExec(t *testing.T, stderr bool) {
	var output []byte
	setInput([]byte(`xyz`), func() {
		if stderr {
			output = captureOutput(&os.Stderr, func() {
				switch runtime.GOOS {
				case "windows":
					scr.Exec("powershell", `foreach ($l in $input) { [Console]::Error.WriteLine($l) }`)
				default:
					scr.Exec("sh", "-c", `cat >&2`)
				}
			})
		} else {
			output = captureOutput(&os.Stdout, func() {
				switch runtime.GOOS {
				case "windows":
					scr.Exec("powershell", `$input | write-output`)
				default:
					scr.Exec("cat")
				}
			})
		}
	})
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = "\r\n"
	}
	if !bytes.Equal(output, []byte(`xyz`+suffix)) {
		t.Fatalf("%v\n%s", output, output)
	}
}

func TestExec(t *testing.T) {
	for _, stderr := range []bool{false, true} {
		t.Run(fmt.Sprintf("stderr=%v", stderr), func(t *testing.T) {
			testExec(t, stderr)
		})
	}
}

func TestExec_err(t *testing.T) {
	defer expectRecoverErr(t, "exit status 1")
	switch runtime.GOOS {
	case "windows":
		scr.Exec("powershell", `exit 1`)
	default:
		scr.Exec("sh", "-c", `exit 1`)
	}
}

func TestPanicf(t *testing.T) {
	defer expectRecover(t, errors.New("err: 123"))
	scr.Panicf("%v: %v", "err", "123")
}

func TestWrapf(t *testing.T) {
	defer expectRecover(t, errors.New("err: 123: 456"))
	defer scr.Wrapf("%v: %v", "err", "123")
	scr.Panicf("456")
}
