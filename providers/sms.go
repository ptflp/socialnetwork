package providers

import "context"

type SMSProvider interface {
	Send(ctx context.Context, phone, message string) error
}
