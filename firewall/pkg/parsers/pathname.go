package parsers

import (
	"errors"
	"fmt"
	"strings"
)

func ParsePathname(format, pathname string) (map[string]string, error) {
	result := make(map[string]string)
	result["format"] = format

	switch format {
	case "maven":
		parts := strings.Split(pathname, "/")
		if len(parts) < 4 {
			return nil, errors.New("invalid maven pathname format")
		}
		result["groupId"] = strings.Join(parts[:len(parts)-3], ".")
		result["artifactId"] = parts[len(parts)-3]
		result["version"] = parts[len(parts)-2]
		filename, found := strings.CutPrefix(parts[len(parts)-1], strings.Join(parts[len(parts)-3:len(parts)-1], "-"))
		if !found {
			return nil, errors.New("no artifactId or version in filename")
		}
		index := strings.LastIndex(filename, ".")
		if index == -1 {
			return nil, errors.New("no extension")
		}
		result["extension"] = filename[index+1:]

		if filename[0] == '-' {
			result["classifier"] = filename[1:index]
		}
	case "npm":
		path, found := strings.CutSuffix(pathname, ".tgz")
		if !found {
			return nil, errors.New("wrong extension")
		}
		scopeAndPackage, filename, found := strings.Cut(path, "/-/")
		if !found {
			return nil, errors.New("invalid npm pathname format")
		}
		result["packageId"] = scopeAndPackage
		index := strings.LastIndex(filename, "-")
		if index == -1 {
			return nil, errors.New("invalid filename")
		}
		result["version"] = filename[index+1:]
	case "pypi":
		parts := strings.Split(pathname, "/")
		if len(parts) != 4 || parts[0] != "packages" {
			return nil, errors.New("invalid pypi pathname format")
		}
		result["name"] = parts[1]
		result["version"] = parts[2]

	case "golang":
		name, version, found := strings.Cut(pathname, "/@v/")
		if !found {
			return nil, errors.New("invalid golang pathname format")
		}

		result["name"] = name
		if version[0] == 'v' {
			index := strings.LastIndex(version, ".")
			if index >= 0 {
				result["version"] = version[:index]
			}
		}

	default:
		return nil, fmt.Errorf("unsupported format %s", format)
	}

	return result, nil
}
