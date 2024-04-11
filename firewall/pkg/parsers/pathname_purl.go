package parsers

import "github.com/package-url/packageurl-go"

func PathnameToPurl(format, pathname string) (*packageurl.PackageURL, error) {
	coordinates, err := ParsePathname(format, pathname)
	if err != nil {
		return nil, err
	}
	return IdentifyPurl(coordinates)
}
