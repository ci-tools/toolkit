package core

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ci-tools/toolkit/ptr"
)

type InputOptions struct {
	/** Optional. Whether the input is Required. If Required and not present, will return errors. Defaults to false */
	Required *bool

	/** Optional. Whether leading/trailing whitespace will be trimmed for the input. Defaults to true */
	TrimWhitespace *bool

	initOnce sync.Once
}

func initializeInputOptions(options *InputOptions) *InputOptions {
	if options == nil {
		options = &InputOptions{}
	}
	options.initOnce.Do(options.init)
	return options
}

// DO NOT CALL THIS OUTSIDE initializeInputOptions
func (o *InputOptions) init() {
	if o.Required == nil {
		o.Required = ptr.Bool(false)
	}
	if o.TrimWhitespace == nil {
		o.TrimWhitespace = ptr.Bool(true)
	}
}

// GetMultilineInput obtains the values of a multiline input.  Each value is also trimmed.
func GetMultilineInput(name string, options *InputOptions) ([]string, error) {
	withNewLine, err := GetInput(name, options)
	if err != nil {
		return nil, err
	}
	inputs := strings.Split(withNewLine, "\n")
	res := make([]string, 0, len(inputs))
	for _, line := range inputs {
		if line != "" {
			res = append(res, line)
		}
	}
	return res, nil
}

// GetInput obtains the value of an input.
// Unless TrimWhitespace is set to false in InputOptions, the value is also trimmed.
// Returns an empty string if the value is not defined.
func GetInput(name string, options *InputOptions) (string, error) {
	options = initializeInputOptions(options)
	envKey := fmt.Sprintf("INPUT_%s", strings.ToTitle(strings.ReplaceAll(name, " ", "_")))
	val := os.Getenv(envKey)
	if *options.Required && val == "" {
		return val, fmt.Errorf("Input required and not supplied: %s", name)
	}
	if !*options.TrimWhitespace {
		return val, nil
	}
	return strings.TrimSpace(val), nil
}

// GetBooleanInput obtains the input value of the boolean type in the YAML 1.2 "core schema" specification.
// Support boolean input list: `true | True | TRUE | false | False | FALSE` .
// The return value is also in boolean type.
// ref: https://yaml.org/spec/1.2/spec.html#id2804923
func GetBooleanInput(name string, options *InputOptions) (bool, error) {
	val, err := GetInput(name, options)
	if err != nil {
		return false, err
	}
	switch val {
	case "true", "True", "TRUE":
		return true, nil
	case "false", "False", "FALSE":
		return false, nil
	default:
		return false, fmt.Errorf(getBooleanErrFmt, val)
	}
}

const getBooleanErrFmt = "Input does not meet YAML 1.2 \"Core Schema\" specification: %s\nSupport boolean input list: \\`true | True | TRUE | false | False | FALSE\\`"
