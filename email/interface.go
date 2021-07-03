package email

import (
	"bytes"
)

type Mailer interface {
	Send(msg Messager) error
}

type Messager interface {
	SetFrom(from string)

	SetReceiver(rcpt string)
	GetReceiver() string

	SetType(t string)

	SetSubject(sub string)

	SetBody(msg bytes.Buffer)

	OpenFile(path string) error
	Attach(src bytes.Buffer, fileName string)
	Bytes() []byte

	Validate() error
}
