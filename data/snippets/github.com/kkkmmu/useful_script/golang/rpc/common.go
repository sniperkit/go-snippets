package common

/*Go RPC的函数只有符合下面的条件才能被远程访问，不然会被忽略，详细的要求如下：
函数必须是导出的(首字母大写)
必须有两个导出类型的参数，
第一个参数是接收的参数，第二个参数是返回给客户端的参数，第二个参数必须是指针类型的
函数还要有一个返回值error
举个例子，正确的RPC函数格式如下：
func (t *T) MethodName(argType T1, replyType *T2) error
*/

import (
	"errors"
)

type Args struct {
	A, B float32
}
type Result struct {
	Value float32
}
type MathService struct {
}

func (s *MathService) Add(args *Args, result *Result) error {
	result.Value = args.A + args.B
	return nil
}

func (s *MathService) Divide(args *Args, result *Result) error {
	if args.B == 0 {
		return errors.New("除数不能为零！")
	}
	result.Value = args.A / args.B
	return nil
}
