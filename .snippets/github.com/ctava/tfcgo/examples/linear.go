package main

import (
	"fmt"
	"log"

	"github.com/ctava/tfcgo"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

func main() {

	// linear model = ((W * x) + b) - y

	var (
		s = op.NewScope()

		initWValue             = op.Const(s.SubScope("W"), float32(0.3))
		initW, handleW, valueW = tfcgo.Variable(s, initWValue)

		x     = op.Placeholder(s.SubScope("x"), tf.Float)
		WxMul = op.AssignVariableOp(s, handleW, op.Mul(s, valueW, x))

		initbValue             = op.Const(s.SubScope("b"), float32(-0.3))
		initb, handleb, valueb = tfcgo.Variable(s, initbValue)

		WxPlusb      = op.AssignAddVariableOp(s, handleW, valueb)
		y            = op.Placeholder(s.SubScope("y"), tf.Float)
		WxPlusMinusY = op.AssignSubVariableOp(s, handleb, y)

		initLossValue                   = op.Const(s.SubScope("loss"), float32(1.0))
		initLoss, handleLoss, valueLoss = tfcgo.Variable(s, initLossValue)
		initDeltaValue                  = op.Const(s.SubScope("Delta"), float32(0.01))
		//TODO reduce_sum and square
		GradientDescent = op.ResourceApplyGradientDescent(s, handleLoss, valueLoss, initDeltaValue)
	)

	g, err := s.Finalize()
	if err != nil {
		log.Fatal(err)
	}
	sess, err := tf.NewSession(g, nil)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := sess.Run(nil, nil, []*tf.Operation{initW}); err != nil {
		log.Fatal(err)
	}

	if _, err := sess.Run(nil, nil, []*tf.Operation{initb}); err != nil {
		log.Fatal(err)
	}

	if _, err := sess.Run(nil, nil, []*tf.Operation{initLoss}); err != nil {
		log.Fatal(err)
	}

	xS := []float32{1.0, 2.0, 3.0, 4.0}
	xTrain, err := tf.NewTensor(xS)
	if err != nil {
		log.Fatal(err)
	}

	yS := []float32{-0.0, -1.0, -2.0, -3.0}
	yTrain, err := tf.NewTensor(yS)
	if err != nil {
		log.Fatal(err)
	}

	err = g.WriteGraphAsText()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i <= 10; i++ {
		_, err := sess.Run(map[tf.Output]*tf.Tensor{x: xTrain, y: yTrain},
			nil, []*tf.Operation{GradientDescent})
		if err != nil {
			log.Fatal(s.Err())
			log.Fatal(err)
		}
	}

	result, err := sess.Run(map[tf.Output]*tf.Tensor{x: xTrain, y: yTrain},
		[]tf.Output{g.Operation("Variable/Read/ReadVariableOp").Output(0), g.Operation("Variable_1/Read/ReadVariableOp").Output(0), g.Operation("Variable_2/Read/ReadVariableOp").Output(0)}, []*tf.Operation{WxMul, WxPlusb, WxPlusMinusY, GradientDescent})
	if err != nil {
		log.Fatal(s.Err())
		log.Fatal(err)
	}
	fmt.Printf("W: %v b: %v loss: %v \n", result[0].Value(), result[1].Value(), result[2].Value())

}
