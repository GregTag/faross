package service

import (
	"firewall/internal/entity"
	"fmt"
	"time"

	"github.com/package-url/packageurl-go"
)

type packageView struct {
	Type  string
	Name  string
	State string
	Score float64
	Date  time.Time
}

func toView(pkg *entity.Package) packageView {
	purl, _ := packageurl.FromString(pkg.Purl)

	var name string
	if len(purl.Namespace) != 0 {
		name = fmt.Sprint("{}/{}", purl.Namespace, purl.Name)
	} else {
		name = purl.Name
	}
	return packageView{
		Type:  purl.Type,
		Name:  name,
		State: pkg.State.ToSring(),
		Score: pkg.FinalScore,
		Date:  pkg.UpdatedAt,
	}
}

func (s *Service) PrepareData( /*...*/ ) (any, error) {
	pkgs, err := s.storage.Package.GetAll()
	if err != nil {
		return nil, err
	}

	views := make([]packageView, 0, len(pkgs))
	for _, pkg := range pkgs {
		views = append(views, toView(&pkg))
	}

	return struct {
		Entries []packageView
	}{
		Entries: views,
	}, nil
}
