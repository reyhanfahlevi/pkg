package filexporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v2"
)

const (
	FormatYAML = "yaml"
	FormatJSON = "json"
	FormatCSV  = "csv"
)

// ExportJSON is quick call to export the given data into json file
func ExportJSON(data interface{}, path, fileName string) error {
	return Export(FormatJSON, data, path, fileName)
}

// ExportYaml is quick call to export the given data into yaml file
func ExportYaml(data interface{}, path, fileName string) error {
	return Export(FormatYAML, data, path, fileName)
}

// ExportCSV is quick call to export the given data into csv file
// data must be a slice of struct otherwise will error
func ExportCSV(data interface{}, path, fileName string) error {
	return Export(FormatCSV, data, path, fileName)
}

// Export will export the given data into specified file format, if not specified then it will
// using default json file format
func Export(format string, data interface{}, path, fileName string) error {
	var (
		fileData []byte
		err      error
	)

	switch format {
	case FormatJSON:
		fileData, err = json.MarshalIndent(data, "", " ")
		if err != nil {
			return err
		}

		if !strings.Contains(strings.ToLower(fileName), ".json") {
			fileName = fmt.Sprintf("%s.json", fileName)
		}
	case FormatYAML:
		fileData, err = yaml.Marshal(data)
		if err != nil {
			return err
		}

		if !strings.Contains(strings.ToLower(fileName), ".yaml") {
			fileName = fmt.Sprintf("%s.yaml", fileName)
		}
	case FormatCSV:
		fileData, err = gocsv.MarshalBytes(data)
		if err != nil {
			return err
		}

		if !strings.Contains(strings.ToLower(fileName), ".csv") {
			fileName = fmt.Sprintf("%s.csv", fileName)
		}
	default:
		return Export(FormatJSON, data, path, fileName)
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
