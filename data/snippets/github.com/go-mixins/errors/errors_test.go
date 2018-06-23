package errors_test

import (
	errs "errors"
	"testing"

	"github.com/go-mixins/errors"
)

func Test_NewClass(t *testing.T) {
	c1 := errors.NewClass("root")
	c2 := c1.Sub("leaf")
	if c1.Is(c2) {
		t.Errorf("%q should be subclass of %q", c1, c2)
	}
	if !c2.Is(c1) {
		t.Errorf("%q should not be subclass of %q", c1, c2)
	}
}

func Test_New_Wrap_Wrapf_Errorf(t *testing.T) {
	err := errors.New("simple error")
	if _, ok := err.(error); !ok {
		t.Errorf("%q should implement error interface", err)
	}
	err = errors.Errorf("simple error %d", 1)
	if _, ok := err.(error); !ok {
		t.Errorf("%q should implement error interface", err)
	}
	err = errors.Wrap(errs.New("simple error 2"), "wrapped in root")
	if _, ok := err.(error); !ok {
		t.Errorf("%q should implement error interface", err)
	}
	err = errors.Wrapf(errs.New("simple error 3"), "wrapped in %s", "root")
	if _, ok := err.(error); !ok {
		t.Errorf("%q should implement error interface", err)
	}
	if err = errors.Wrap(nil, "wrapped nil"); err != nil {
		t.Errorf("%q should be nil", err)
	}
	if err = errors.Wrapf(nil, "wrapped nil in %s", "root"); err != nil {
		t.Errorf("%q should be nil", err)
	}
}

func TestClass_Contains(t *testing.T) {
	root := errors.NewClass("root")
	leaf := errors.NewClass("root", "leaf")
	rootError := root.Wrap(errs.New("some error"), "wrapped in root")
	leafError := leaf.Wrap(errs.New("some other error"), "wrapped in leaf")
	outsideError := errs.New("some outside error")
	if !root.Contains(leafError) {
		t.Errorf("%q should belong to %q", leafError, root)
	}
	if !root.Contains(rootError) {
		t.Errorf("%q should belong to %q", rootError, root)
	}
	if leaf.Contains(rootError) {
		t.Errorf("%q should not belong to %q", rootError, leaf)
	}
	if !leaf.Contains(leafError) {
		t.Errorf("%q should belong to %q", leafError, leaf)
	}
	if root.Contains(outsideError) {
		t.Errorf("%q should not belong to %q", outsideError, root)
	}
	if leaf.Contains(outsideError) {
		t.Errorf("%q should not belong to %q", outsideError, leaf)
	}
}

func TestClass_ErrorsCause(t *testing.T) {
	root := errors.NewClass("root")
	err := errs.New("some error")
	rootError := root.Wrap(err, "in root")
	if errors.Cause(rootError) != err {
		t.Errorf("%q should be the cause of %q", err, rootError)
	}
}

func TestClass_New_Errorf_Wrap_Wrapf(t *testing.T) {
	c1 := errors.NewClass("root")
	err := c1.New("root error 0")
	if !c1.Contains(err) {
		t.Errorf("%q should belong to %q", err, c1)
	}
	err = errs.New("some other error")
	if c1.Contains(err) {
		t.Errorf("%q should not belong to %q", err, c1)
	}
	err = c1.Errorf("root error %d", 1)
	if !c1.Contains(err) {
		t.Errorf("%q should belong to %q", err, c1)
	}
	err = c1.Wrap(errs.New("root error 2"), "wrapped")
	if !c1.Contains(err) {
		t.Errorf("%q should belong to %q", err, c1)
	}
	err = c1.Wrapf(errs.New("root error"), "wrapped %d", 3)
	if !c1.Contains(err) {
		t.Errorf("%q should belong to %q", err, c1)
	}
	if err = c1.Wrap(nil, "wrapped nil"); err != nil {
		t.Errorf("%q should be nil", err)
	}
	if err = c1.Wrapf(nil, "wrapped nil in %q", c1); err != nil {
		t.Errorf("%q should be nil", err)
	}
}
