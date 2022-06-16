package fileparser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
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

// ParseCSVFile will parse the csv file and store it to the struct
func ParseCSVFile(filename string, target interface{}) error {
	csvFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return gocsv.Unmarshal(bytes.NewReader(csvFile), target)
}

// ParseCSVFromGin parse uploaded files and store it to the struct
func ParseCSVFromGin(c *gin.Context, fileName string, target interface{}) error {
	file, err := c.FormFile(fileName)
	if err != nil {
		return err
	}

	open, _ := file.Open()
	return gocsv.Unmarshal(open, target)
}
