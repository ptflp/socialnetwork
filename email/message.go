package email

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"gitlab.com/InfoBlogFriends/server/validators"
)

const (
	Delimiter = "**=myohmy689407924327"

	MessageHtml  = "text/html"
	MessagePlain = "text/plain"
)

type Message struct {
	typeRaw     string
	fromRaw     string
	toRaw       string
	from        string
	to          string
	cc          string
	subject     string
	mimeVersion string
	contentType string
	delimiter   string
	body        body
}

func NewMessage() *Message {
	return &Message{
		typeRaw:     MessagePlain,
		mimeVersion: "MIME-Version: 1.0\r\n",
		contentType: fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", Delimiter),
		delimiter:   fmt.Sprintf("\r\n--%s\r\n", Delimiter),
		body: body{
			contentType:             fmt.Sprintf("Content-Type: %s; charset=\"utf-8\"\r\n", MessagePlain),
			contentTransferEncoding: "Content-Transfer-Encoding: 7bit\r\n",
		},
	}
}

type body struct {
	contentType             string
	contentTransferEncoding string
	body                    string
	files                   files
}

func (b *body) String() string {
	return concat(b.contentType, b.contentTransferEncoding, b.body, "\r\n", b.files.String())
}

func concat(args ...string) string {
	d := strings.Join(args, "")
	return d
}

type files []file

func (f *files) String() string {
	s := make([]string, 0, len(*f))
	for _, fl := range *f {
		s = append(s, fl.String())
	}

	return strings.Join(s, "")
}

type file struct {
	fileName string
	src      bytes.Buffer
}

func newFile(fileName string) *file {
	return &file{
		fileName: fileName,
	}
}

func (f *file) String() string {
	return strings.Join([]string{
		fmt.Sprintf("\r\n--%s\r\n", Delimiter),
		"Content-Type: text/plain; charset=\"utf-8\"\r\n",
		"Content-Transfer-Encoding: base64\r\n",
		"Content-Disposition: attachment;filename=\"" + f.fileName + "\"\r\n",
		"\r\n" + base64.StdEncoding.EncodeToString(f.src.Bytes()),
	}, "")
}

func (m *Message) SetFrom(from string) {
	m.fromRaw = from
	m.from = fmt.Sprintf("From: %s\r\n", from)
}

func (m *Message) SetReceiver(rcpt string) {
	m.toRaw = rcpt
	m.to = fmt.Sprintf("To: %s\r\n", rcpt)
	m.cc = fmt.Sprintf("Cc: %s\r\n", rcpt)
}

func (m *Message) GetReceiver() string {
	return m.toRaw
}

func (m *Message) SetType(t string) {
	m.typeRaw = t
	m.body.contentType = fmt.Sprintf("Content-Type: %s; charset=\"utf-8\"\r\n", t)
}

func (m *Message) SetSubject(sub string) {
	m.subject = fmt.Sprintf("Subject: %s\r\n", sub)
}

func (m *Message) AttachFile(src bytes.Buffer, fileName string) {
	f := newFile(fileName)
	f.src.Write(src.Bytes())
	m.body.files = append(m.body.files, *f)
}

func (m *Message) Validate() error {
	err := validators.CheckEmailFormat(m.fromRaw)
	if err != nil {
		return err
	}
	err = validators.CheckEmailFormat(m.toRaw)
	if err != nil {
		return err
	}
	if m.body.body == "" {
		return errors.New("body not set")
	}
	if m.typeRaw != MessagePlain && m.typeRaw != MessageHtml {
		return fmt.Errorf("wrong email body content-type %s", m.typeRaw)
	}

	return nil
}

func (m *Message) String() string {
	return concat(
		m.from,
		m.to,
		m.cc,
		m.subject,
		m.mimeVersion,
		m.contentType,
		m.delimiter,
		m.body.String())
}

func (m *Message) Bytes() []byte {
	return []byte(m.String())
}

func (m *Message) SetBody(msg bytes.Buffer) {
	m.body.body = msg.String()
}

func (m *Message) OpenFile(path string) error {
	fl, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("path: %s error: %s", func() string { s, _ := os.Getwd(); return s }(), err)
	}
	defer fl.Close()
	fileInfo, err := fl.Stat()
	if err != nil {
		return fmt.Errorf("path: %s error: %s", func() string { s, _ := os.Getwd(); return s }(), err)
	}

	fb := make([]byte, fileInfo.Size())
	_, err = fl.Read(fb)
	b := bytes.NewBuffer(fb)

	m.AttachFile(*b, fl.Name())

	return err
}

func (m *Message) ImplementsMessager() {
	panic("implement me")
}
