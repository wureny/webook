package sms

import (
	"context"
	"github.com/wureny/webook/webook/Internal/service/sms/memory"
)

type Service interface {
	//number 为...string,是因为可能会发送多个手机号
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
}

func Newseservice() Service {
	return memory.NewService()
}

type NamedArg struct {
	Val  string
	Name string
}
