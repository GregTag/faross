package parsers_test

import (
	"firewall/pkg/parsers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePathname(t *testing.T) {
	// Maven test case
	coordinates, err := parsers.ParsePathname("maven", "com/example/my-library/1.0.0/my-library-1.0.0.jar")
	expected := map[string]string{
		"format":     "maven",
		"groupId":    "com.example",
		"artifactId": "my-library",
		"version":    "1.0.0",
		"extension":  "jar",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, coordinates)

	// NPM test case
	coordinates, err = parsers.ParsePathname("npm", "kind-of/-/kind-of-3.2.2.tgz")
	expected = map[string]string{
		"format":    "npm",
		"packageId": "kind-of",
		"version":   "3.2.2",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, coordinates)

	coordinates, err = parsers.ParsePathname("npm", "is-integer/-/is-integer-1.0.7.tgz")
	expected = map[string]string{
		"format":    "npm",
		"packageId": "is-integer",
		"version":   "1.0.7",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, coordinates)

	coordinates, err = parsers.ParsePathname("npm", "@reduxjs/toolkit/-/toolkit-2.2.3.tgz")
	expected = map[string]string{
		"format":    "npm",
		"packageId": "@reduxjs/toolkit",
		"version":   "2.2.3",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, coordinates)

	coordinates, err = parsers.ParsePathname("npm", "@oriano-dev/is-even/-/is-even-1.0.1.tgz")
	expected = map[string]string{
		"format":    "npm",
		"packageId": "@oriano-dev/is-even",
		"version":   "1.0.1",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, coordinates)

	// PyPI test case
	coordinates, err = parsers.ParsePathname("pypi", "packages/is-even/1.0.7/is_even-1.0.7-py2.py3-none-any.whl")
	expected = map[string]string{
		"format":  "pypi",
		"name":    "is-even",
		"version": "1.0.7",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, coordinates)

	coordinates, err = parsers.ParsePathname("pypi", "packages/numpy/1.26.4/numpy-1.26.4-cp310-cp310-manylinux_2_17_x86_64.manylinux2014_x86_64.whl")
	expected = map[string]string{
		"format":  "pypi",
		"name":    "numpy",
		"version": "1.26.4",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, coordinates)

	// Golang test case
	coordinates, err = parsers.ParsePathname("golang", "github.com/package-url/packageurl-go/@v/v0.1.0.info")
	expected = map[string]string{
		"format":  "golang",
		"name":    "github.com/package-url/packageurl-go",
		"version": "v0.1.0",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, coordinates)

	coordinates, err = parsers.ParsePathname("golang", "github.com/package-url/packageurl-go/@v/list")
	expected = map[string]string{
		"format": "golang",
		"name":   "github.com/package-url/packageurl-go",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, coordinates)

	coordinates, err = parsers.ParsePathname("golang", "gorm.io/driver/sqlite/@v/list")
	expected = map[string]string{
		"format": "golang",
		"name":   "gorm.io/driver/sqlite",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, coordinates)
}
