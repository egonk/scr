//go:build go1.18
// +build go1.18

package scr_test

import (
	"errors"
	"testing"

	"github.com/egonk/scr"
)

func mustFunc(s string, err error) (string, error) {
	return s, err
}

func TestMust(t *testing.T) {
	defer expectRecover(t, errors.New("err"))
	scr.Must(mustFunc("str", errors.New("err")))
}

func TestMust_nil(t *testing.T) {
	if m := scr.Must(mustFunc("str", nil)); m != "str" {
		t.Fatal(m)
	}
}
