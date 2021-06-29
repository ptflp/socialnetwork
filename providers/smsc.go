package providers

import "context"

type SMSC struct {
	login string
	pswd  string
}

func (s *SMSC) Send(ctx context.Context, phone, msg string) error {
	return nil
}
