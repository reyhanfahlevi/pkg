package token

import (
	"testing"
)

func TestGenerateBytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Success",
			args: args{
				n: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateBytes(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGenerateString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Success",
			args: args{
				n: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateString(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
