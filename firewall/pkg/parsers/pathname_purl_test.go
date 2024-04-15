package parsers_test

import (
	"firewall/pkg/parsers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathnameToPurl(t *testing.T) {
	package_url, err := parsers.PathnameToPurl("maven", "com/example/my-library/1.0.0/my-library-1.0.0.jar")
	assert.NoError(t, err)
	assert.Equal(t, "pkg:maven/com.example/my-library@1.0.0?type=jar", package_url.ToString())
	package_url, err = parsers.PathnameToPurl("npm", "kind-of/-/kind-of-3.2.2.tgz")
	assert.NoError(t, err)
	assert.Equal(t, "pkg:npm/kind-of@3.2.2", package_url.ToString())
	package_url, err = parsers.PathnameToPurl("npm", "is-integer/-/is-integer-1.0.7.tgz")
	assert.NoError(t, err)
	assert.Equal(t, "pkg:npm/is-integer@1.0.7", package_url.ToString())
	package_url, err = parsers.PathnameToPurl("npm", "@reduxjs/toolkit/-/toolkit-2.2.3.tgz")
	assert.NoError(t, err)
	assert.Equal(t, "pkg:npm/%40reduxjs/toolkit@2.2.3", package_url.ToString())
	package_url, err = parsers.PathnameToPurl("npm", "@oriano-dev/is-even/-/is-even-1.0.1.tgz")
	assert.NoError(t, err)
	assert.Equal(t, "pkg:npm/%40oriano-dev/is-even@1.0.1", package_url.ToString())
	package_url, err = parsers.PathnameToPurl("pypi", "packages/is-even/1.0.7/is_even-1.0.7-py2.py3-none-any.whl")
	assert.NoError(t, err)
	assert.Equal(t, "pkg:pypi/is-even@1.0.7", package_url.ToString())
	package_url, err = parsers.PathnameToPurl("pypi", "packages/numpy/1.26.4/numpy-1.26.4-cp310-cp310-manylinux_2_17_x86_64.manylinux2014_x86_64.whl")
	assert.NoError(t, err)
	assert.Equal(t, "pkg:pypi/numpy@1.26.4", package_url.ToString())
	package_url, err = parsers.PathnameToPurl("golang", "github.com/package-url/packageurl-go/@v/v0.1.0.info")
	assert.NoError(t, err)
	assert.Equal(t, "pkg:golang/github.com/package-url/packageurl-go@v0.1.0", package_url.ToString())
	package_url, err = parsers.PathnameToPurl("golang", "github.com/package-url/packageurl-go/@v/list")
	assert.NoError(t, err)
	assert.Equal(t, "pkg:golang/github.com/package-url/packageurl-go", package_url.ToString())
	package_url, err = parsers.PathnameToPurl("golang", "gorm.io/driver/sqlite/@v/list")
	assert.NoError(t, err)
	assert.Equal(t, "pkg:golang/gorm.io/driver/sqlite", package_url.ToString())
}
