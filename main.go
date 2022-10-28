package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
)

func main() {
	var originFile string
	var patchFile string
	var format string
	pflag.StringVarP(&originFile, "origin", "o", "", "The origin json/toml file path")
	pflag.StringVarP(&patchFile, "patch", "p", "", "The patching json/toml file in json-merge-patch format")
	pflag.StringVarP(&format, "format", "f", "auto", "auto|json|yaml|toml, default auto(same as origin)")
	pflag.Parse()

	output, err := PatchConfigFile(originFile, patchFile, format)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))

}
