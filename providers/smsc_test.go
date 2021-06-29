package providers

import (
	"net/http"
	"testing"

	"gitlab.com/InfoBlogFriends/server/config"
)

func TestSMSC_buildUrl(t *testing.T) {
	type fields struct {
		client *http.Client
		cfg    *config.SMSC
	}
	type args struct {
		phone string
		msg   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "common build url test",
			fields: fields{
				client: nil,
				cfg: &config.SMSC{
					Pwd:   "parol",
					Login: "login",
					Cost:  "3",
					Fmt:   "3",
				},
			},
			args: args{
				phone: "79644288083",
				msg:   "message",
			},
			wantErr: false,
			want:    "https://smsc.ru/sys/send.php?cost=3&fmt=3&login=login&mes=message&phones=79644288083&psw=parol",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SMSC{
				client: tt.fields.client,
				cfg:    tt.fields.cfg,
			}
			got, err := s.buildUrl(tt.args.phone, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("buildQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}
