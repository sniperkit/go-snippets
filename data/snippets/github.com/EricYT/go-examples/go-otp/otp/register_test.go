package otp

import "testing"

func TestRegisterName(t *testing.T) {
	name := Name("foo")
	if err := registerName(name, nil); err != nil {
		t.Errorf("register name should not error: %s", err)
	}
	defer unregisterName(name)
	if err := registerName(name, nil); err == nil {
		t.Errorf("register name should return error already exists")
	}
}

func TestUnregisterName(t *testing.T) {
	name := Name("foo")
	if err := registerName(name, nil); err != nil {
		t.Errorf("register name should not error: %s", err)
	}
	if err := unregisterName(name); err != nil {
		t.Errorf("unregister name error: %s", err)
	}
	if err := unregisterName(name); err == nil {
		t.Errorf("unregister name already not exists should error")
	}
}

func TestGetGenByName(t *testing.T) {
	name := Name("foo")
	if err := registerName(name, nil); err != nil {
		t.Errorf("register name should not error: %s", err)
	}
	if _, err := getGenByName(name); err != nil {
		t.Errorf("register name already put in, but we got a error: %s", err)
	}
	unregisterName(name)
	if _, err := getGenByName(name); err == nil {
		t.Errorf("register name already delete, but we got a entry")
	}
}
