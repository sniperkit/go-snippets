package optional

import "testing"

func TestOptional(t *testing.T) {
	f := func(v interface{}) Optional {
		return &Just{v}
	}
	g := func(v interface{}) Optional {
		return &None{}
	}
	h := func(v interface{}) Optional {
		return &Just{nil}
	}
	o := Wrap(1)

	if !o.HasValue() || o.Get() == nil {
		t.Error("must have value")
	}

	if value := o.Apply(f).Or(2); value != 1 {
		t.Errorf("want 1, got %v", value)
	}

	o = g(2)
	if o.HasValue() || o.Get() != nil {
		t.Error("must not have value")
	}

	if value := o.Apply(g).Or(2); value != 2 {
		t.Errorf("want 2, got %v", value)
	}

	o = h(3)
	if value := o.Or(3); value != 3 {
		t.Errorf("want 3, got %v", value)
	}
}
