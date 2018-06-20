package tfcgo

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

// Variable - this function creats a variable for training purposes.
//
// Returns:
// the init operation - very important - this needs to be run in a session to be initialized
// a handle to the variable - for use in assignment operations
// an Output that produces the current value of the variable
func Variable(scope *op.Scope, initialValue tf.Output) (init *tf.Operation, handle, value tf.Output) {
	scope = scope.SubScope("Variable")
	dtype := initialValue.DataType()
	handle = op.VarHandleOp(scope, dtype, initialValue.Shape())
	init = op.AssignVariableOp(scope.SubScope("Assign"), handle, initialValue)
	value = op.ReadVariableOp(scope.SubScope("Read"), handle, dtype)
	return init, handle, value
}

// Variable - this function creats a variable for training purposes.
//
// Returns:
// the init operation - very important - this needs to be run in a session to be initialized
// a handle to the variable - for use in assignment operations
// an Output that produces the current value of the variable
func VariableAndTensor(scope *op.Scope, initialValue tf.Output, v interface{}) (t *tf.Tensor, init *tf.Operation, handle, value tf.Output) {

	t, err := tf.NewTensor(v)
	if err != nil {
		return nil, nil, tf.Output{}, tf.Output{}
	}

	scope = scope.SubScope("Variable")
	dtype := initialValue.DataType()

	handle = op.VarHandleOp(scope, dtype, initialValue.Shape())
	init = op.AssignVariableOp(scope.SubScope("Assign"), handle, initialValue)
	value = op.ReadVariableOp(scope.SubScope("Read"), handle, dtype)

	return t, init, handle, value
}
