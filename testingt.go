package testutils

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestHelper interface {
	Helper()
}

type MsgAndArgs struct {
	MSG  string
	Args []interface{}
}
type FakeTest struct {
	ErrorFormats  []MsgAndArgs
	ErrorMessages []string
	Failed        bool
	Name          string
}

func (t *FakeTest) Errorf(msg string, args ...interface{}) {
	if !t.Failed {
		t.ErrorFormats = append(t.ErrorFormats, MsgAndArgs{MSG: msg, Args: args})
		t.ErrorMessages = append(t.ErrorMessages, fmt.Sprintf(msg, args...))
	}
}

func (t *FakeTest) FailNow() {
	t.Failed = true
}

func (t *FakeTest) String() string {
	s := "--- PASS: "
	if t.Failed || len(t.ErrorMessages) > 0 {
		s = "--- FAIL: "
	}
	return s + t.Name
}

var _ assert.TestingT = &FakeTest{}
var _ require.TestingT = &FakeTest{}
