package sms

import "context"

type Service interface {
	//number 为...string,是因为可能会发送多个手机号
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
}

type NamedArg struct {
	Val  string
	Name string
}
