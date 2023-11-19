package retryable

import (
	"context"
	"errors"
	"github.com/wureny/webook/webook/Internal/service/sms"
)

type Service struct {
	svc      sms.Service
	retryMax int
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)
	cnt := 1
	for err != nil && cnt < s.retryMax {
		err = s.svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("重试都失败了")
}
