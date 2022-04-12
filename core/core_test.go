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

func Test_GetBooleanInput(t *testing.T) {
	setEnvVars(t, []env{
		{key: "INPUT_BOOLEAN_INPUT", val: "true"},
		{key: "INPUT_BOOLEAN_INPUT_TRUE1", val: "true"},
		{key: "INPUT_BOOLEAN_INPUT_TRUE2", val: "True"},
		{key: "INPUT_BOOLEAN_INPUT_TRUE3", val: "TRUE"},
		{key: "INPUT_BOOLEAN_INPUT_FALSE1", val: "false"},
		{key: "INPUT_BOOLEAN_INPUT_FALSE2", val: "False"},
		{key: "INPUT_BOOLEAN_INPUT_FALSE3", val: "FALSE"},
	})

	table := []struct {
		name     string
		options  *InputOptions
		expected bool
	}{
		{name: "boolean input", options: nil, expected: true},
		{name: "boolean input", options: &InputOptions{Required: ptr.Bool(true)}, expected: true},
		{name: "boolean input true1", options: nil, expected: true},
		{name: "boolean input true2", options: nil, expected: true},
		{name: "boolean input true3", options: nil, expected: true},
		{name: "boolean input false1", options: nil, expected: false},
		{name: "boolean input false2", options: nil, expected: false},
		{name: "boolean input false2", options: nil, expected: false},
	}

	for _, tt := range table {
		val, err := GetBooleanInput(tt.name, tt.options)
		if err != nil {
			t.Fatal("errors are not expected but found:", err)
		}
		if val != tt.expected {
			t.Fatalf("expected value is %v, but result was %v.\ntest case: %v", tt.expected, val, table)
		}
	}
}

/* TODO: delete cases below when implemented
// MIT License
// Copyright 2019 GitHub
describe('@actions/core', () => {
  beforeAll(() => {
    const filePath = path.join(__dirname, `test`)
    if (!fs.existsSync(filePath)) {
      fs.mkdirSync(filePath)
    }
  })

  it('legacy exportVariable produces the correct command and sets the env', () => {
    core.exportVariable('my var', 'var val')
    assertWriteCalls([`::set-env name=my var::var val${os.EOL}`])
  })

  it('legacy exportVariable escapes variable names', () => {
    core.exportVariable('special char var \r\n,:', 'special val')
    expect(process.env['special char var \r\n,:']).toBe('special val')
    assertWriteCalls([
      `::set-env name=special char var %0D%0A%2C%3A::special val${os.EOL}`
    ])
  })

  it('legacy exportVariable escapes variable values', () => {
    core.exportVariable('my var2', 'var val\r\n')
    expect(process.env['my var2']).toBe('var val\r\n')
    assertWriteCalls([`::set-env name=my var2::var val%0D%0A${os.EOL}`])
  })

  it('legacy exportVariable handles boolean inputs', () => {
    core.exportVariable('my var', true)
    assertWriteCalls([`::set-env name=my var::true${os.EOL}`])
  })

  it('legacy exportVariable handles number inputs', () => {
    core.exportVariable('my var', 5)
    assertWriteCalls([`::set-env name=my var::5${os.EOL}`])
  })

  it('exportVariable produces the correct command and sets the env', () => {
    const command = 'ENV'
    createFileCommandFile(command)
    core.exportVariable('my var', 'var val')
    verifyFileCommand(
      command,
      `my var<<_GitHubActionsFileCommandDelimeter_${os.EOL}var val${os.EOL}_GitHubActionsFileCommandDelimeter_${os.EOL}`
    )
  })

  it('exportVariable handles boolean inputs', () => {
    const command = 'ENV'
    createFileCommandFile(command)
    core.exportVariable('my var', true)
    verifyFileCommand(
      command,
      `my var<<_GitHubActionsFileCommandDelimeter_${os.EOL}true${os.EOL}_GitHubActionsFileCommandDelimeter_${os.EOL}`
    )
  })

  it('exportVariable handles number inputs', () => {
    const command = 'ENV'
    createFileCommandFile(command)
    core.exportVariable('my var', 5)
    verifyFileCommand(
      command,
      `my var<<_GitHubActionsFileCommandDelimeter_${os.EOL}5${os.EOL}_GitHubActionsFileCommandDelimeter_${os.EOL}`
    )
  })

  it('setSecret produces the correct command', () => {
    core.setSecret('secret val')
    assertWriteCalls([`::add-mask::secret val${os.EOL}`])
  })

  it('prependPath produces the correct commands and sets the env', () => {
    const command = 'PATH'
    createFileCommandFile(command)
    core.addPath('myPath')
    expect(process.env['PATH']).toBe(
      `myPath${path.delimiter}path1${path.delimiter}path2`
    )
    verifyFileCommand(command, `myPath${os.EOL}`)
  })

  it('legacy prependPath produces the correct commands and sets the env', () => {
    core.addPath('myPath')
    expect(process.env['PATH']).toBe(
      `myPath${path.delimiter}path1${path.delimiter}path2`
    )
    assertWriteCalls([`::add-path::myPath${os.EOL}`])
  })

  it('getBooleanInput handles wrong boolean input', () => {
    expect(() => core.getBooleanInput('wrong boolean input')).toThrow(
      'Input does not meet YAML 1.2 "Core Schema" specification: wrong boolean input\n' +
        `Support boolean input list: \`true | True | TRUE | false | False | FALSE\``
    )
  })

  it('setOutput produces the correct command', () => {
    core.setOutput('some output', 'some value')
    assertWriteCalls([
      os.EOL,
      `::set-output name=some output::some value${os.EOL}`
    ])
  })

  it('setOutput handles bools', () => {
    core.setOutput('some output', false)
    assertWriteCalls([os.EOL, `::set-output name=some output::false${os.EOL}`])
  })

  it('setOutput handles numbers', () => {
    core.setOutput('some output', 1.01)
    assertWriteCalls([os.EOL, `::set-output name=some output::1.01${os.EOL}`])
  })

  it('setFailed sets the correct exit code and failure message', () => {
    core.setFailed('Failure message')
    expect(process.exitCode).toBe(core.ExitCode.Failure)
    assertWriteCalls([`::error::Failure message${os.EOL}`])
  })

  it('setFailed escapes the failure message', () => {
    core.setFailed('Failure \r\n\nmessage\r')
    expect(process.exitCode).toBe(core.ExitCode.Failure)
    assertWriteCalls([`::error::Failure %0D%0A%0Amessage%0D${os.EOL}`])
  })

  it('setFailed handles Error', () => {
    const message = 'this is my error message'
    core.setFailed(new Error(message))
    expect(process.exitCode).toBe(core.ExitCode.Failure)
    assertWriteCalls([`::error::Error: ${message}${os.EOL}`])
  })

  it('error sets the correct error message', () => {
    core.error('Error message')
    assertWriteCalls([`::error::Error message${os.EOL}`])
  })

  it('error escapes the error message', () => {
    core.error('Error message\r\n\n')
    assertWriteCalls([`::error::Error message%0D%0A%0A${os.EOL}`])
  })

  it('error handles an error object', () => {
    const message = 'this is my error message'
    core.error(new Error(message))
    assertWriteCalls([`::error::Error: ${message}${os.EOL}`])
  })

  it('error handles parameters correctly', () => {
    const message = 'this is my error message'
    core.error(new Error(message), {
      title: 'A title',
      file: 'root/test.txt',
      startColumn: 1,
      endColumn: 2,
      startLine: 5,
      endLine: 5
    })
    assertWriteCalls([
      `::error title=A title,file=root/test.txt,line=5,endLine=5,col=1,endColumn=2::Error: ${message}${os.EOL}`
    ])
  })

  it('warning sets the correct message', () => {
    core.warning('Warning')
    assertWriteCalls([`::warning::Warning${os.EOL}`])
  })

  it('warning escapes the message', () => {
    core.warning('\r\nwarning\n')
    assertWriteCalls([`::warning::%0D%0Awarning%0A${os.EOL}`])
  })

  it('warning handles an error object', () => {
    const message = 'this is my error message'
    core.warning(new Error(message))
    assertWriteCalls([`::warning::Error: ${message}${os.EOL}`])
  })

  it('warning handles parameters correctly', () => {
    const message = 'this is my error message'
    core.warning(new Error(message), {
      title: 'A title',
      file: 'root/test.txt',
      startColumn: 1,
      endColumn: 2,
      startLine: 5,
      endLine: 5
    })
    assertWriteCalls([
      `::warning title=A title,file=root/test.txt,line=5,endLine=5,col=1,endColumn=2::Error: ${message}${os.EOL}`
    ])
  })

  it('notice sets the correct message', () => {
    core.notice('Notice')
    assertWriteCalls([`::notice::Notice${os.EOL}`])
  })

  it('notice escapes the message', () => {
    core.notice('\r\nnotice\n')
    assertWriteCalls([`::notice::%0D%0Anotice%0A${os.EOL}`])
  })

  it('notice handles an error object', () => {
    const message = 'this is my error message'
    core.notice(new Error(message))
    assertWriteCalls([`::notice::Error: ${message}${os.EOL}`])
  })

  it('notice handles parameters correctly', () => {
    const message = 'this is my error message'
    core.notice(new Error(message), {
      title: 'A title',
      file: 'root/test.txt',
      startColumn: 1,
      endColumn: 2,
      startLine: 5,
      endLine: 5
    })
    assertWriteCalls([
      `::notice title=A title,file=root/test.txt,line=5,endLine=5,col=1,endColumn=2::Error: ${message}${os.EOL}`
    ])
  })

  it('annotations map field names correctly', () => {
    const commandProperties = toCommandProperties({
      title: 'A title',
      file: 'root/test.txt',
      startColumn: 1,
      endColumn: 2,
      startLine: 5,
      endLine: 5
    })
    expect(commandProperties.title).toBe('A title')
    expect(commandProperties.file).toBe('root/test.txt')
    expect(commandProperties.col).toBe(1)
    expect(commandProperties.endColumn).toBe(2)
    expect(commandProperties.line).toBe(5)
    expect(commandProperties.endLine).toBe(5)

    expect(commandProperties.startColumn).toBeUndefined()
    expect(commandProperties.startLine).toBeUndefined()
  })

  it('startGroup starts a new group', () => {
    core.startGroup('my-group')
    assertWriteCalls([`::group::my-group${os.EOL}`])
  })

  it('endGroup ends new group', () => {
    core.endGroup()
    assertWriteCalls([`::endgroup::${os.EOL}`])
  })

  it('group wraps an async call in a group', async () => {
    const result = await core.group('mygroup', async () => {
      process.stdout.write('in my group\n')
      return true
    })
    expect(result).toBe(true)
    assertWriteCalls([
      `::group::mygroup${os.EOL}`,
      'in my group\n',
      `::endgroup::${os.EOL}`
    ])
  })

  it('debug sets the correct message', () => {
    core.debug('Debug')
    assertWriteCalls([`::debug::Debug${os.EOL}`])
  })

  it('debug escapes the message', () => {
    core.debug('\r\ndebug\n')
    assertWriteCalls([`::debug::%0D%0Adebug%0A${os.EOL}`])
  })

  it('saveState produces the correct command', () => {
    core.saveState('state_1', 'some value')
    assertWriteCalls([`::save-state name=state_1::some value${os.EOL}`])
  })

  it('saveState handles numbers', () => {
    core.saveState('state_1', 1)
    assertWriteCalls([`::save-state name=state_1::1${os.EOL}`])
  })

  it('saveState handles bools', () => {
    core.saveState('state_1', true)
    assertWriteCalls([`::save-state name=state_1::true${os.EOL}`])
  })

  it('getState gets wrapper action state', () => {
    expect(core.getState('TEST_1')).toBe('state_val')
  })

  it('isDebug check debug state', () => {
    const current = process.env['RUNNER_DEBUG']
    try {
      delete process.env.RUNNER_DEBUG
      expect(core.isDebug()).toBe(false)

      process.env['RUNNER_DEBUG'] = '1'
      expect(core.isDebug()).toBe(true)
    } finally {
      process.env['RUNNER_DEBUG'] = current
    }
  })

  it('setCommandEcho can enable echoing', () => {
    core.setCommandEcho(true)
    assertWriteCalls([`::echo::on${os.EOL}`])
  })

  it('setCommandEcho can disable echoing', () => {
    core.setCommandEcho(false)
    assertWriteCalls([`::echo::off${os.EOL}`])
  })
})

// Assert that process.stdout.write calls called only with the given arguments.
function assertWriteCalls(calls: string[]): void {
  expect(process.stdout.write).toHaveBeenCalledTimes(calls.length)

  for (let i = 0; i < calls.length; i++) {
    expect(process.stdout.write).toHaveBeenNthCalledWith(i + 1, calls[i])
  }
}

function createFileCommandFile(command: string): void {
  const filePath = path.join(__dirname, `test/${command}`)
  process.env[`GITHUB_${command}`] = filePath
  fs.appendFileSync(filePath, '', {
    encoding: 'utf8'
  })
}

function verifyFileCommand(command: string, expectedContents: string): void {
  const filePath = path.join(__dirname, `test/${command}`)
  const contents = fs.readFileSync(filePath, 'utf8')
  try {
    expect(contents).toEqual(expectedContents)
  } finally {
    fs.unlinkSync(filePath)
  }
}

function getTokenEndPoint(): string {
  return 'https://vstoken.actions.githubusercontent.com/.well-known/openid-configuration'
}

describe('oidc-client-tests', () => {
  it('Get Http Client', async () => {
    const http = new HttpClient('actions/oidc-client')
    expect(http).toBeDefined()
  })

  it('HTTP get request to get token endpoint', async () => {
    const http = new HttpClient('actions/oidc-client')
    const res = await http.get(getTokenEndPoint())
    expect(res.message.statusCode).toBe(200)
  })
})
*/
