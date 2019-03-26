package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
)


const ValidateInput = 1
const ValidateOutput = 0

type Args struct {
	input  string
	output string
	sep    string
}

func (a *Args) Parse() {
	var err error = nil

	a.input, err = a.expandHomeDir(a.input)
	if err != nil {
		log.Fatal("Unable to expand home directory in input path")
	}

	a.output, err = a.expandHomeDir(a.output)
	if err != nil {
		log.Fatal("Unable to expand home directory in output path")
	}
}

func (a *Args) Validate() {
	if err := a.validatePath(a.input, ValidateInput); err != nil {
		log.Fatal(err)
	}

	if err := a.validatePath(a.output, ValidateOutput); err != nil {
		log.Fatal(err)
	}
}

func (a Args) expandHomeDir(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", home, 1)
	}

	return path, nil
}

func (a Args) validatePath(path string, mode int) error {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		if mode == ValidateInput {
			return errors.New(fmt.Sprintf("Input: %s doesn't exist", path) )
		} else {
			return nil
		}
	}

	if info.IsDir() {
		return errors.New(fmt.Sprintf("Input: %s is a directory", path))
	}

	return nil
}
