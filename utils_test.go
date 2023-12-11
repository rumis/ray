package ray

import (
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				v: struct {
					Name  string `qs:"name"`
					Email string `qs:"email"`
				}{
					Name:  "John Doe",
					Email: "johndoe@example.com",
				},
			},
			want:    "email=johndoe%40example.com&name=John+Doe",
			wantErr: false,
		},
		{
			name: "Test case 2",
			args: args{
				v: struct {
					ID    int    `qs:"id"`
					Token string `qs:"token"`
				}{
					ID:    123,
					Token: "abc123",
				},
			},
			want:    "id=123&token=abc123",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
