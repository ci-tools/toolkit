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
