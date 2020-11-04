package fileparser

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseJSONFile(t *testing.T) {
	type args struct {
		filename string
		target   interface{}
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(a args)
		wantErr bool
	}{
		{
			name: "Test Success",
			args: args{
				filename: "test.json",
				target: struct {
					Test string `json:"test"`
				}{},
			},
			mockFn: func(a args) {
				_ = ioutil.WriteFile(a.filename, []byte(`{"test":"aku test"}`), os.ModePerm)
			},
			wantErr: false,
		}, {
			name: "Test Failed - Not Found",
			args: args{
				filename: "test.json",
				target: struct {
					Test string `json:"test"`
				}{},
			},
			mockFn:  func(a args) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			if err := ParseJSONFile(tt.args.filename, &tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("ParseJSONFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			_ = os.RemoveAll(tt.args.filename)
		})
	}
}

func TestParseYamlFile(t *testing.T) {
	type args struct {
		filename string
		target   interface{}
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(a args)
		wantErr bool
	}{
		{
			name: "Test Success",
			args: args{
				filename: "test.yaml",
				target: struct {
					Test string `yaml:"test"`
				}{},
			},
			mockFn: func(a args) {
				_ = ioutil.WriteFile(a.filename, []byte(`test: "aku test"`), os.ModePerm)
			},
			wantErr: false,
		}, {
			name: "Test Failed - Not Found",
			args: args{
				filename: "test.yaml",
				target: struct {
					Test string `yaml:"test"`
				}{},
			},
			mockFn:  func(a args) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			if err := ParseYamlFile(tt.args.filename, &tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("ParseYamlFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			_ = os.RemoveAll(tt.args.filename)
		})
	}
}

func TestParseCSVFile(t *testing.T) {
	type csv struct {
		Test string `csv:"test"`
	}

	type args struct {
		filename string
		target   interface{}
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(a args)
		wantErr bool
	}{
		{
			name: "Test Success",
			args: args{
				filename: "test.csv",
				target:   &[]csv{},
			},
			mockFn: func(a args) {
				_ = ioutil.WriteFile(a.filename, []byte("test\ntesss"), os.ModePerm)
			},
			wantErr: false,
		}, {
			name: "Test Failed - Not Found",
			args: args{
				filename: "test.csv",
				target:   &[]csv{},
			},
			mockFn:  func(a args) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			if err := ParseCSVFile(tt.args.filename, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("ParseCSVFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			_ = os.RemoveAll(tt.args.filename)
		})
	}
}
