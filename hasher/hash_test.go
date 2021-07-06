package hasher

import (
	"reflect"
	"testing"
)

func TestNewSHA256(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test sha256 length",
			args: args{
				data: []byte("test"),
			},
			want: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		},
		{
			name: "",
			args: args{
				data: []byte(""),
			},
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "",
			args: args{
				data: []byte("abc"),
			},
			want: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		},
		{
			name: "",
			args: args{
				data: []byte("hello"),
			},
			want: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSHA256(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}
