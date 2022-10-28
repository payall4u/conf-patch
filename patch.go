package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type ConfigType interface {
	Decode(content []byte) (interface{}, error)
	// Indent by default
	Encode(v interface{}) ([]byte, error)
}

type JsonConfigType string

func (j JsonConfigType) Decode(content []byte) (interface{}, error) {
	var result interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		return "", err
	}
	return result, nil
}
func (j JsonConfigType) Encode(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

type TomlConfigType string

func (t TomlConfigType) Decode(content []byte) (interface{}, error) {
	var result interface{}
	if err := toml.Unmarshal(content, &result); err != nil {
		return "", err
	}
	return result, nil
}
func (t TomlConfigType) Encode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)

	err := encoder.Encode(v)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

var ConfigTypes = map[string]ConfigType{
	"json": JsonConfigType("json"),
	"toml": TomlConfigType("toml"),
}

func getFileType(filename string) string {
	return strings.TrimLeft(strings.ToLower(filepath.Ext(filename)), ".")
}

func Load(filename string) (interface{}, error) {
	ext := getFileType(filename)
	configType, ok := ConfigTypes[ext]
	if !ok {
		return nil, fmt.Errorf("Unsupported file type: %s", ext)
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return configType.Decode(content)
}

// @format: json|yaml|toml
func Output(v interface{}, format string) ([]byte, error) {
	configType, ok := ConfigTypes[format]
	if !ok {
		return nil, fmt.Errorf("Unsupported file type: %s", format)
	}
	return configType.Encode(v)
}

func PatchConfigFile(originFile string, patchFile string, format string) ([]byte, error) {
	var origin, patch interface{}

	origin, err := Load(originFile)
	if err != nil {
		return nil, err
	}

	patch, err = Load(patchFile)
	if err != nil {
		return nil, err
	}

	if format == "auto" || format == "" {
		format = getFileType(originFile)
	}

	result := PatchInterface(origin, patch)

	output, err := Output(result, format)
	if err != nil {
		return nil, err
	}

	return output, nil
}
