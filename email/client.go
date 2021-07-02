package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"gitlab.com/InfoBlogFriends/server/config"

	"go.uber.org/zap"
)

type Client struct {
	smtpc  *smtp.Client
	logger *zap.Logger
	cfg    *config.Email
}

func NewClient(cfg *config.Email, logger *zap.Logger) (*Client, error) {
	c := &Client{
		logger: logger,
		cfg:    cfg,
	}
	err := c.connect()

	return c, err
}

func (c *Client) Send(msg Messager) error {
	msg.SetFrom(c.cfg.From)
	err := msg.Validate()
	if err != nil {
		return err
	}
	err = c.connect()
	if err != nil {
		return err
	}
	rclient := c.smtpc

	if err = rclient.Rcpt(msg.GetReceiver()); err != nil {
		c.logger.Error("set email receiver", zap.String("To:", msg.GetReceiver()), zap.Error(err))
		return err
	}
	writer, err := rclient.Data()
	if err != nil {
		c.logger.Error("smpt client Data()", zap.Error(err))
		return err
	}

	//write into email client stream writter
	if _, err = writer.Write(msg.Bytes()); err != nil {
		c.logger.Error("write content into client writter I/O", zap.Error(err))
		return err
	}

	if err = writer.Close(); err != nil {
		c.logger.Error("smtp writer close", zap.Error(err))
	}

	return err
}

func (c *Client) connect() error {
	tlsConfig := tls.Config{
		ServerName:         c.cfg.ServerAddress,
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", c.cfg.ServerAddress, c.cfg.Port), &tlsConfig)
	if err != nil {
		c.logger.Error("TLS connection", zap.Error(err))
		return err
	}

	rclient, err := smtp.NewClient(conn, c.cfg.ServerAddress)

	if err != nil {
		c.logger.Error("smtp client creation", zap.Error(err))
		return err
	}

	auth := smtp.PlainAuth("", c.cfg.Login, c.cfg.Password, c.cfg.ServerAddress)

	if err = rclient.Auth(auth); err != nil {
		c.logger.Error("smtp auth", zap.Error(err))
		return err
	}

	if err = rclient.Mail(c.cfg.Login); err != nil {
		c.logger.Error("start mail transaction", zap.Error(err))
		return err
	}

	c.smtpc = rclient

	return nil
}
