package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"gitlab.com/InfoBlogFriends/server/config"
)

type SMSC struct {
	client *http.Client
	cfg    *config.SMSC
}

func NewSMSC(cfg *config.SMSC) *SMSC {
	return &SMSC{
		client: &http.Client{},
		cfg:    cfg,
	}
}

func (s *SMSC) Send(ctx context.Context, phone, msg string) error {
	if s.cfg.Dev {
		return nil
	}
	smscUrl, err := s.buildUrl(phone, msg)
	if err != nil {
		return err
	}

	resp, err := s.call(ctx, smscUrl)
	_ = resp
	if err != nil {
		return err
	}
	return nil
}

func (s *SMSC) buildUrl(phone, msg string) (string, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "smsc.ru",
		Path:   "sys/send.php",
	}
	q := u.Query()
	q.Set("login", s.cfg.Login)
	q.Set("psw", s.cfg.Pwd)
	q.Set("phones", phone)
	q.Set("mes", msg)
	q.Set("cost", s.cfg.Cost)
	q.Set("fmt", s.cfg.Fmt)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (s *SMSC) call(ctx context.Context, url string) (*Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	r, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	response := &Response{}
	err = json.NewDecoder(r.Body).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type Response struct {
	ID      int    `json:"id"`
	Cnt     int    `json:"cnt"`
	Cost    string `json:"cost"`
	Balance string `json:"balance"`
}
