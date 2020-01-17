package filexporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// ExportJSON is quick call to export the given data into json file
func ExportJSON(data interface{}, path, fileName string) error {
	return Export("json", data, path, fileName)
}

// ExportYaml is quick call to export the given data into yaml file
func ExportYaml(data interface{}, path, fileName string) error {
	return Export("yaml", data, path, fileName)
}

// Export will export the given data into specified file format, if not specified then it will
// using default json file format
func Export(format string, data interface{}, path, fileName string) error {
	var (
		fileData []byte
		err      error
	)

	switch format {
	case "json":
		fileData, err = json.MarshalIndent(data, "", " ")
		if err != nil {
			return err
		}

		if !strings.Contains(strings.ToLower(fileName), ".json") {
			fileName = fmt.Sprintf("%s.json", fileName)
		}
	case "yaml":
		fileData, err = yaml.Marshal(data)
		if err != nil {
			return err
		}

		if !strings.Contains(strings.ToLower(fileName), ".yaml") {
			fileName = fmt.Sprintf("%s.yaml", fileName)
		}
	default:
		return Export("json", data, path, fileName)
	}

	if path != "" {
		_, err = os.Open(path)
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		}

		fileName = filepath.Join(path, fileName)
	}

	return ioutil.WriteFile(fileName, fileData, os.ModePerm)
}
