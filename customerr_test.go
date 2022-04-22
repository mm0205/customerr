package customerr

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestNew1(t *testing.T) {
	var testError1 = errors.New("testError1")
	err := New(testError1, []Tag{}, "wrap message %d", 10)
	if !errors.Is(err, testError1) {
		t.Fatalf("errors.Is(err, testError1)")
	}
	if "wrap message 10: [testError1]" != err.Error() {
		t.Fatalf("invalid error message %s", err.Error())
	}
}

func TestNew2(t *testing.T) {
	type testStruct struct{}
	x := testStruct{}
	testError1 := json.Unmarshal([]byte("{,}"), &x)
	err := New(testError1, []Tag{}, "wrap message %d", 10)
	if !errors.Is(err, testError1) {
		t.Fatalf("errors.Is(err, testError1)")
	}
	syntaxError := &json.SyntaxError{}
	if !errors.As(err, &syntaxError) {
		t.Fatalf("invalid customError")
	}

	if !strings.HasPrefix(err.Error(), "wrap message 10: [") {
		t.Fatalf("invalid error message %s", err.Error())
	}
}

func TestIsCustomErr(t *testing.T) {
	if IsCustomErr(errors.New("testError")) {
		t.Fatalf("IsCustomErr returns true for testError")
	}

	target := New(errors.New("inner"), []Tag{}, "custom")
	if !IsCustomErr(target) {
		t.Fatalf("IsCustomErr returns false for customError")
	}
}

func TestHasTag(t *testing.T) {
	const testTag1 = "abc"
	target1 := New(errors.New("inner"), []Tag{testTag1}, "custom")
	if !HasTag(target1, testTag1) {
		t.Fatalf("Tag %s not found", testTag1)
	}

	const testTag2 = "efg"
	target2 := New(target1, []Tag{testTag2}, "custom2")
	if !HasTag(target2, testTag2) {
		t.Fatalf("Tag %s not found", testTag2)
	}

	if !HasTag(target2, testTag1) {
		t.Fatalf("Tag %s not found", testTag1)
	}

	if HasTag(target1, testTag2) {
		t.Fatalf("Tag %s found", testTag1)
	}
}
