// scr is a convenience package to support writing short automation scripts that
// simply panic on errors. It is similar to shell scripting with "set -e"
// option. scr is the long name variant of https://github.com/egonk/sc
//  package main
//
//  import (
//      "bufio"
//      "fmt"
//      "os"
//      "regexp"
//
//      "github.com/egonk/scr"
//  )
//
//  func main() {
//      // set -e
//      // grep "abc" example
//      re := regexp.MustCompile(`abc`)
//      f := scr.Must(os.Open("example"))
//      defer scr.Close(f)
//      s := bufio.NewScanner(f)
//      for s.Scan() {
//          if re.Match(s.Bytes()) {
//              scr.Must(fmt.Printf("%s\n", s.Bytes()))
//          }
//      }
//      scr.Err(s.Err())
//  }
package scr

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Close c and panic on error.
//  f := scr.Must(os.Create("example"))
//  defer scr.Close(f)
//  // write to f
func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}

// Err panics if err is not nil.
//  scr.Err(os.Chdir("example"))
func Err(err error) {
	if err != nil {
		panic(err)
	}
}

// Exec calls ExecCmd(name, arg...).Run() and panics on error.
//  scr.Exec("git", "status")
func Exec(name string, arg ...string) {
	if err := ExecCmd(name, arg...).Run(); err != nil {
		panic(err)
	}
}

// ExecCmd calls exec.Command(name, arg...) and sets up os.Stdin, os.Stdout and
// os.Stderr in the returned exec.Cmd.
//  c := scr.ExecCmd("go", "build", "-v", ".")
//  c.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
//  scr.Err(c.Run())
func ExecCmd(name string, arg ...string) *exec.Cmd {
	c := exec.Command(name, arg...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

// Panicf panics with fmt.Errorf(format, a...).
//  scr.Panicf("invalid argument: %v", arg)
func Panicf(format string, a ...interface{}) {
	panic(fmt.Errorf(format, a...))
}

// Wrapf panics with fmt.Errorf(format+": %v", append(a, recover())...) if
// recover() is not nil.
//  defer scr.Wrapf("file: %v", fn) // panic: file: <fn>: invalid argument: <arg>
//  scr.Panicf("invalid argument: %v", arg)
func Wrapf(format string, a ...interface{}) {
	if r := recover(); r != nil {
		panic(fmt.Errorf(format+": %v", append(a, r)...))
	}
}
