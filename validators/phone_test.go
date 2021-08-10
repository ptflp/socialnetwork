package validators

import "testing"

func TestCheckPhoneFormat(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "casting phone number from 89644288083 to 79644288083",
			args:    args{phone: "89644288083"},
			want:    "79644288083",
			wantErr: false,
		},
		{
			name:    "casting phone number from +79644288083 to 79644288083",
			args:    args{phone: "+79644288083"},
			want:    "79644288083",
			wantErr: false,
		},
		// TODO: fix short number want error
		{
			name:    "check short number",
			args:    args{phone: "+7964428808"},
			want:    "7964428808",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckPhoneFormat(tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPhoneFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckPhoneFormat() got = %v, want %v", got, tt.want)
			}
		})
	}
}
