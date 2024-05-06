package util

// TODO: add other elements
const OutputTemplate = `{
	"static analysis": { {{range $index, $element := .}}{{ if $index }},{{ end }}
		"{{ $element.ToolName }}": {
			"result": {{ $element.Output }},
			"exit_code": {{ $element.ExitCode }}
		}{{ end }}
	},
	"dynamic analysis": {
		"packj-trace": {
			"result": "",
			"exit-code": 0
		}
	}
}
`
