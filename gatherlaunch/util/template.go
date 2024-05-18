package util

// TODO: add other elements
const OutputTemplate = `{
	"static_analysis": { {{range $index, $element := .}}{{ if ne $element.ToolName "packj-trace" }}{{ if $index }},{{ end }}
		"{{ $element.ToolName }}": {
			"result": {{ $element.Output }},
			"exit_code": {{ $element.ExitCode }}
		}{{ end }}{{ end }}
	},
	"dynamic_analysis": { {{range $index, $element := .}}{{ if eq $element.ToolName "packj-trace" }}
		"{{ $element.ToolName }}": {
			"result": {{ $element.Output }},
			"exit_code": {{ $element.ExitCode }}
		}{{ end }}{{ end }}
	}
}
`
