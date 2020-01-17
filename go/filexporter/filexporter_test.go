package filexporter

import (
	"os"
	"testing"
)

func TestExportJson(t *testing.T) {
	type args struct {
		data     interface{}
		path     string
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Success with sub and not with file type",
			args: args{
				data: struct {
					Test string `json:"test"`
					Sub  struct {
						Test string `json:"test"`
					} `json:"sub"`
				}{
					Test: "ini test aja",
					Sub: struct {
						Test string `json:"test"`
					}{"ini test sub"},
				},
				path:     "example/subfolder/",
				fileName: "test_tanpa_ext",
			},
			wantErr: false,
		}, {
			name: "Test Success 2",
			args: args{
				data: struct {
					Test string `json:"test"`
				}{
					Test: "ini test aja",
				},
				path:     "example",
				fileName: "test.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExportJSON(tt.args.data, tt.args.path, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("ExportJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			_ = os.RemoveAll(tt.args.path)
		})
	}
}

func TestExportYaml(t *testing.T) {
	type args struct {
		data     interface{}
		path     string
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Success 1",
			args: args{
				data: struct {
					Test string `yaml:"test"`
					Sub  struct {
						Test string `yaml:"test"`
					} `yaml:"sub"`
				}{
					Test: "ini test aja",
					Sub: struct {
						Test string `yaml:"test"`
					}{"ini test sub"},
				},
				path:     "example",
				fileName: "test_tanpa_ext",
			},
			wantErr: false,
		}, {
			name: "Test Success 2",
			args: args{
				data: struct {
					Test string `yaml:"test"`
					Sub  struct {
						Test string `yaml:"test"`
					} `yaml:"sub"`
				}{
					Test: "ini test aja",
					Sub: struct {
						Test string `yaml:"test"`
					}{"ini test sub"},
				},
				path:     "example",
				fileName: "test.yaml",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExportYaml(tt.args.data, tt.args.path, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("ExportYaml() error = %v, wantErr %v", err, tt.wantErr)
			}

			_ = os.RemoveAll(tt.args.path)
		})
	}
}
