package providers

import "context"

type SMS interface {
	Send(ctx context.Context, phone, msg string) error
}
