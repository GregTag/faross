package storage_test

import (
	"firewall/internal/entity"
	"firewall/internal/storage"
	"firewall/internal/storage/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPackageStore(t *testing.T) {
	db, err := driver.NewSQLiteDB("file::memory:")
	assert.NoError(t, err)

	s := storage.NewStorage(db)

	purl1 := "pkg:pypi/is-even@1.0.7"
	purl2 := "pkg:pypi/numpy@1.26.4"

	pkg_ptr, err := s.Package.GetByPurl(purl1)
	assert.ErrorIs(t, err, entity.ErrPackageNotFound)
	assert.Nil(t, pkg_ptr)

	pkg1 := entity.Package{
		Purl:       purl1,
		FinalScore: 6.0,
		State:      entity.Quarantined,
	}

	err = s.Package.Save(&pkg1)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, pkg1.ID)

	pkg_ptr, err = s.Package.GetByPurl(purl1)
	assert.NoError(t, err)
	assert.Equal(t, pkg1.ID, pkg_ptr.ID)

	pathname := "packages/numpy/1.26.4/numpy-1.26.4-cp310-cp310-manylinux_2_17_x86_64.manylinux2014_x86_64.whl"

	pkg2 := entity.Package{
		Pathname: pathname,
		Purl:     purl1,
		State:    entity.Healthy,
	}
	err = s.Package.Save(&pkg2)
	assert.NoError(t, err)
	assert.Equal(t, pkg1.ID, pkg2.ID)
	assert.Equal(t, pathname, pkg2.Pathname)
	assert.Equal(t, entity.Healthy, pkg2.State)

	pkg3 := entity.Package{
		Purl:  purl2,
		State: entity.Quarantined,
	}
	assert.NoError(t, s.Package.Save(&pkg3))

	now := time.Now().Add(-time.Minute)

	rep1 := entity.Repository{InstanceID: "aaa", RepositoryDTO: entity.RepositoryDTO{Name: "a"}}
	rep2 := entity.Repository{InstanceID: "bbb", RepositoryDTO: entity.RepositoryDTO{Name: "b"}}
	assert.NoError(t, s.Repository.Save(&rep1))
	assert.NoError(t, s.Repository.Save(&rep2))

	pkgs := []entity.Package{pkg2, pkg3}
	assert.NoError(t, s.Repository.AppendPackages(rep1.ID, pkgs))
	assert.NoError(t, s.Repository.AppendPackages(rep2.ID, pkgs))

	rep3, err := s.Repository.Load("aaa", "a")
	assert.NoError(t, err)
	assert.Len(t, rep3.Packages, 2)

	unq, err := s.Repository.GetUnquarantined(rep1.ID, now)
	assert.NoError(t, err)
	assert.Empty(t, unq)

	comment := "Some interesting comment"
	assert.Error(t, s.Package.Unquarantine(purl1, comment))
	assert.NoError(t, s.Package.Unquarantine(purl2, comment))

	unq, err = s.Repository.GetUnquarantined(rep1.ID, now)
	assert.NoError(t, err)
	assert.Len(t, unq, 1)
	assert.Equal(t, pkg3.ID, unq[0].ID)
	assert.Equal(t, entity.Unquarantined, unq[0].State)
	assert.Equal(t, comment, unq[0].Comment)

	packages, err := s.Package.GetAll()
	assert.NoError(t, err)
	assert.Len(t, packages, 2)

	err = s.Package.UpdateComment(purl1, comment)
	assert.NoError(t, err)

	pkg4, err := s.Package.GetByPurl(purl1)
	assert.NoError(t, err)
	assert.Equal(t, comment, pkg4.Comment)
}
