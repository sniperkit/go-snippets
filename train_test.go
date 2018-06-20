package tfcgo

import (
	"testing"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

func TestVariable(t *testing.T) {
	var (
		s                   = op.NewScope()
		initValue           = op.Const(s.SubScope("init"), int32(1))
		increment           = op.Const(s.SubScope("inc"), int32(3))
		init, handle, value = Variable(s, initValue)
		update              = op.AssignVariableOp(s, handle, op.Add(s, value, increment))
	)
	graph, err := s.Finalize()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := tf.NewSession(graph, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sess.Run(nil, nil, []*tf.Operation{init}); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		if _, err := sess.Run(nil, nil, []*tf.Operation{update}); err != nil {
			t.Fatal(err)
		}
	}
	result, err := sess.Run(nil, []tf.Output{value}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := result[0].Value().(int32), int32(16); got != want {
		t.Errorf("Got %v, want %v", got, want)
	}
}
