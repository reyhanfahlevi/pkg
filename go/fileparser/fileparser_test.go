package fileparser

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/stretchr/testify/assert"
)

func setupRouter(middleware gin.HandlerFunc, routers ...func(r gin.IRoutes) gin.IRoutes) *gin.Engine {
	r := gin.Default()

	g := r.Group("/test")
	if middleware != nil {
		g.Use(middleware)
	}
	for _, router := range routers {
		router(g)
	}

	return r
}

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

func TestParseCSVFromGin(t *testing.T) {
	type csv struct {
		Test string `csv:"test"`
	}

	type args struct {
		c        func() *gin.Context
		fileName string
		target   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Test Success",
			args: args{
				c: func() *gin.Context {
					te := []csv{
						{Test: "test"},
					}
					b, _ := gocsv.MarshalBytes(te)
					buf := new(bytes.Buffer)
					mw := multipart.NewWriter(buf)
					w, err := mw.CreateFormFile("file", "test")
					if assert.NoError(t, err) {
						_, err = w.Write(b)
						assert.NoError(t, err)
					}
					mw.Close()
					c, _ := gin.CreateTestContext(httptest.NewRecorder())
					c.Request, _ = http.NewRequest("POST", "/", buf)
					c.Request.Header.Set("Content-Type", mw.FormDataContentType())
					return c
				},
				fileName: "file",
				target:   &[]csv{},
			},
			wantErr: assert.NoError,
		}, {
			name: "Test Failed",
			args: args{
				c: func() *gin.Context {
					te := []csv{
						{Test: "test"},
					}
					b, _ := gocsv.MarshalBytes(te)
					buf := new(bytes.Buffer)
					mw := multipart.NewWriter(buf)
					w, err := mw.CreateFormFile("file", "test")
					if assert.NoError(t, err) {
						_, err = w.Write(b)
						assert.NoError(t, err)
					}
					mw.Close()
					c, _ := gin.CreateTestContext(httptest.NewRecorder())
					c.Request, _ = http.NewRequest("POST", "/", buf)
					c.Request.Header.Set("Content-Type", mw.FormDataContentType())
					return c
				},
				fileName: "files",
				target:   &[]csv{},
			},
			wantErr: assert.Error,
		}, {
			name: "Test Failed 2",
			args: args{
				c: func() *gin.Context {
					c, _ := gin.CreateTestContext(httptest.NewRecorder())
					c.Request, _ = http.NewRequest("POST", "/", nil)
					return c
				},
				fileName: "files",
				target:   &[]csv{},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, ParseCSVFromGin(tt.args.c(), tt.args.fileName, tt.args.target))
		})
	}
}
