package fileparser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v2"
)

// ParseJSONFile will parse the json file and store it to the struct
func ParseJSONFile(filename string, target interface{}) error {
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonFile, target)
}

// ParseYamlFile will parse the yaml file and store it to the struct
func ParseYamlFile(filename string, target interface{}) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, target)
}

// ParseCSVFile
func ParseCSVFile(filename string, target interface{}) error {
	csvFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return gocsv.Unmarshal(bytes.NewReader(csvFile), target)
}
