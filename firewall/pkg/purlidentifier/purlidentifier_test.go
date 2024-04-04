package purlidentifier_test

import (
	"firewall/pkg/purlidentifier"
	"testing"

	"github.com/package-url/packageurl-go"
	"github.com/stretchr/testify/assert"
)

func TestFromCoordinates(t *testing.T) {
	// Maven test case
	mavenCoordinates := map[string]string{
		"format":     "maven",
		"groupId":    "com.example",
		"artifactId": "my-artifact",
		"version":    "1.0.0",
	}
	expectedMavenPURL := &packageurl.PackageURL{
		Type:      "maven",
		Namespace: "com.example",
		Name:      "my-artifact",
		Version:   "1.0.0",
	}
	purl, err := purlidentifier.FromCoordinates(mavenCoordinates)
	assert.NoError(t, err)
	assert.Equal(t, expectedMavenPURL, purl)

	// NPM test case
	npmCoordinates := map[string]string{
		"format":    "npm",
		"packageId": "my-package",
		"version":   "1.0.0",
	}
	expectedNpmPURL := &packageurl.PackageURL{
		Type:    "npm",
		Name:    "my-package",
		Version: "1.0.0",
	}
	purl, err = purlidentifier.FromCoordinates(npmCoordinates)
	assert.NoError(t, err)
	assert.Equal(t, expectedNpmPURL, purl)

	// PyPI test case
	pypiCoordinates := map[string]string{
		"format":  "pypi",
		"name":    "my-package",
		"version": "1.0.0",
	}
	expectedPypiPURL := &packageurl.PackageURL{
		Type:    "pypi",
		Name:    "my-package",
		Version: "1.0.0",
	}
	purl, err = purlidentifier.FromCoordinates(pypiCoordinates)
	assert.NoError(t, err)
	assert.Equal(t, expectedPypiPURL, purl)

	// Golang test case
	golangCoordinates := map[string]string{
		"format":  "golang",
		"name":    "github.com/my-org/my-repo",
		"version": "v1.0.0",
	}
	expectedGolangPURL := &packageurl.PackageURL{
		Type:      "golang",
		Namespace: "github.com/my-org",
		Name:      "my-repo",
		Version:   "v1.0.0",
	}
	purl, err = purlidentifier.FromCoordinates(golangCoordinates)
	assert.NoError(t, err)
	assert.Equal(t, expectedGolangPURL, purl)
}
