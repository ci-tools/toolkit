package core

import (
	"reflect"
	"testing"

	"github.com/ci-tools/toolkit/ptr"
)

type env struct {
	key string
	val string
}

func setEnvVars(t *testing.T, vars []env) {
	for _, kv := range vars {
		t.Setenv(kv.key, kv.val)
	}
}

/* TODO: remove the comment below after creating tests w/ testing.Setenv
func init() {
	testEnvVars := []env{
		{key: "INPUT_MY_INPUT", val: "val"},
		{key: "INPUT_MISSING", val: ""},
		{key: "INPUT_SPECIAL_CHARS_'\t\"\\", val: "'\t\"\\ response "},
		{key: "INPUT_MULTIPLE_SPACES_VARIABLE", val: "I have multiple spaces"},
		{key: "INPUT_BOOLEAN_INPUT", val: "true"},
		{key: "INPUT_BOOLEAN_INPUT_TRUE2", val: "true"},
		{key: "INPUT_BOOLEAN_INPUT_TRUE2", val: "True"},
		{key: "INPUT_BOOLEAN_INPUT_TRUE3", val: "TRUE"},
		{key: "INPUT_BOOLEAN_INPUT_FALSE1", val: "false"},
		{key: "INPUT_BOOLEAN_INPUT_FALSE2", val: "False"},
		{key: "INPUT_BOOLEAN_INPUT_FALSE3", val: "FALSE"},
		{key: "INPUT_WRONG_BOOLEAN_INPUT", val: "wrong"},
		{key: "INPUT_WITH_TRAILING_WHITESPACE", val: "  some val  "},
		{key: "INPUT_MY_INPUT_LIST", val: "val1\nval2\nval3"},
	}

	for _, kv := range testEnvVars {
		if err := os.Setenv(kv.key, kv.val); err != nil {
			log.Println("failed to set an environment variable:", err)
			os.Exit(1)
		}
	}
}
*/

func Test_GetInput(t *testing.T) {
	setEnvVars(t, []env{
		{key: "INPUT_MY_INPUT", val: "val"},
		{key: "INPUT_MISSING", val: ""},
		{key: "INPUT_SPECIAL_CHARS_'\t\"\\", val: "'\t\"\\ response "},
		{key: "INPUT_MULTIPLE_SPACES_VARIABLE", val: "I have multiple spaces"},
		{key: "INPUT_WITH_TRAILING_WHITESPACE", val: "  some val  "},
	})

	table := []struct {
		name     string
		options  *InputOptions
		expected string
	}{
		{name: "my input", options: nil, expected: "val"},
		{name: "my input", options: &InputOptions{}, expected: "val"},
		{name: "my input", options: &InputOptions{Required: ptr.Bool(true)}, expected: "val"},
		{name: "missing", options: &InputOptions{Required: ptr.Bool(false)}, expected: ""},
		{name: "My InPuT", options: nil, expected: "val"},
		{name: "special chars_'\t\"\\", options: nil, expected: "'\t\"\\ response"},
		{name: "multiple spaces variable", options: nil, expected: "I have multiple spaces"},
		{name: "with trailing whitespace", options: nil, expected: "some val"},
		{name: "with trailing whitespace", options: &InputOptions{TrimWhitespace: ptr.Bool(true)}, expected: "some val"},
		{name: "with trailing whitespace", options: &InputOptions{TrimWhitespace: ptr.Bool(false)}, expected: "  some val  "},
	}

	for _, tt := range table {
		val, err := GetInput(tt.name, tt.options)
		if err != nil {
			t.Fatal("errors are not expected but found:", err)
		}
		if val != tt.expected {
			t.Fatalf("expected value is %s, but result was %s.\ntest case: %v", tt.expected, val, table)
		}
	}
}

func Test_GetMultilineInput(t *testing.T) {
	setEnvVars(t, []env{
		{key: "INPUT_MY_INPUT_LIST", val: "val1\nval2\nval3"},
	})

	table := []struct {
		name     string
		options  *InputOptions
		expected []string
	}{
		{name: "my input list", options: nil, expected: []string{"val1", "val2", "val3"}},
	}

	for _, tt := range table {
		val, err := GetMultilineInput(tt.name, tt.options)
		if err != nil {
			t.Fatal("errors are not expected but found:", err)
		}
		if !reflect.DeepEqual(val, tt.expected) {
			t.Fatalf("expected value is %s, but result was %s.\ntest case: %v", tt.expected, val, table)
		}
	}
}
