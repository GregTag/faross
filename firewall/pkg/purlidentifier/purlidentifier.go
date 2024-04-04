package purlidentifier

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/package-url/packageurl-go"
)

var trimSlashesRegex = regexp.MustCompile(`(^/+)|(/+$)`)

func FromCoordinates(coordinates map[string]string) (*packageurl.PackageURL, error) {
	format, ok := coordinates["format"]
	if !ok {
		return nil, fmt.Errorf("no format")
	}
	builder := builder{errs: make([]error, 0)}
	switch format {
	case "maven":
		builder.resolveMaven(coordinates)
	case "npm":
		builder.resolveNpm(coordinates)
	case "pypi":
		builder.resolvePypi(coordinates)
	case "golang":
		builder.resolveGolang(coordinates)
	default:
		builder.errs = append(builder.errs, fmt.Errorf("unsupported format"))
	}
	if len(builder.errs) != 0 {
		return nil, errors.Join(builder.errs...)
	}
	return &builder.purl, nil
}

type builder struct {
	purl packageurl.PackageURL
	errs []error
}

func (b *builder) get(m *map[string]string, k string) string {
	v, ok := (*m)[k]
	if !ok {
		b.errs = append(b.errs, fmt.Errorf("no value for %s", k))
	}
	return v
}

func (b *builder) buildQualifier(m *map[string]string, coordinateName, purlCoordinateName string) {
	value, ok := (*m)[coordinateName]
	if ok {
		b.purl.Qualifiers = append(b.purl.Qualifiers, packageurl.Qualifier{Key: purlCoordinateName, Value: value})
	}
}

func (b *builder) resolveNameAndNamespace(packageId string) {
	packageId = trimSlashesRegex.ReplaceAllString(packageId, "")
	var namespace string
	if strings.Contains(packageId, "/") {
		lastSlashIndex := strings.LastIndex(packageId, "/")
		namespace = packageId[:lastSlashIndex]
		packageId = packageId[lastSlashIndex+1:]
	}
	b.purl.Name = packageId
	if len(namespace) != 0 {
		b.purl.Namespace = namespace
	}
}

func (b *builder) resolveMaven(coordinates map[string]string) {
	b.purl.Type = "maven"
	b.purl.Name = b.get(&coordinates, "artifactId")
	b.purl.Namespace = b.get(&coordinates, "groupId")
	b.purl.Version = b.get(&coordinates, "version")
	b.buildQualifier(&coordinates, "extension", "type")
	b.buildQualifier(&coordinates, "classifier", "classifier")
}

func (b *builder) resolveNpm(coordinates map[string]string) {
	b.purl.Type = "npm"
	b.purl.Version = b.get(&coordinates, "version")
	b.resolveNameAndNamespace(b.get(&coordinates, "packageId"))
}

func (b *builder) resolvePypi(coordinates map[string]string) {
	b.purl.Type = "pypi"
	b.resolveNameAndNamespace(b.get(&coordinates, "name"))
	b.purl.Version = b.get(&coordinates, "version")
	b.buildQualifier(&coordinates, "qualifier", "qualifier")
	b.buildQualifier(&coordinates, "extension", "extension")
}

func (b *builder) resolveGolang(coordinates map[string]string) {
	b.purl.Type = "golang"
	b.purl.Version = b.get(&coordinates, "version")
	b.resolveNameAndNamespace(b.get(&coordinates, "name"))
}
