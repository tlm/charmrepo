// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package testing // import "github.com/juju/charmrepo/v7/testing"

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/juju/charm/v9"
	"github.com/juju/utils/v3/fs"
)

const (
	// BundlePathSegment is the postfix segment added to file paths to get the
	// directory used for storing bundles in a local repo.
	BundlePathSegment = "bundles"

	// CharmsPathSegment is the postfix segment added to file paths to get the
	// directory used for storing charms in a local repo.
	CharmsPathSegment = "charms"
)

// Repo represents a charm repository on disk that is used for testing.
// Charm repositories are expected to conform to the following directory
// structure on disk.
// ./
//    charms/
//      charm1/
//      charm2/
//    bundles/
//      bundle1/
//      bundle2/
type Repo struct {
	path string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// NewRepo returns a new testing charm repository rooted at the given
// path, relative to the package directory of the calling package.
func NewRepo(path string) *Repo {
	// Find the repo directory. This is only OK to do
	// because this is running in a test context
	// so we know the source is available.
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("cannot get caller to determine repo path")
	}
	r := &Repo{
		path: filepath.Join(filepath.Dir(file), path),
	}
	if _, err := os.Stat(r.path); err != nil {
		panic(fmt.Errorf("cannot read repository found at %q: %v", r.path, err))
	}
	return r
}

func (r *Repo) Path() string {
	return r.path
}

func clone(dst, src string) string {
	dst = filepath.Join(dst, filepath.Base(src))
	check(fs.Copy(src, dst))
	return dst
}

// BundleDirPath returns the path to a bundle directory with the given name in the
// default series
func (r *Repo) BundleDirPath(name string) string {
	return filepath.Join(r.Path(), BundlePathSegment, name)
}

// BundleDir returns the actual charm.BundleDir named name.
func (r *Repo) BundleDir(name string) *charm.BundleDir {
	b, err := charm.ReadBundleDir(r.BundleDirPath(name))
	check(err)
	return b
}

// CharmDirPath returns the path to a charm directory with the given name in the
// default series
func (r *Repo) CharmDirPath(name string) string {
	return filepath.Join(r.Path(), CharmsPathSegment, name)
}

// CharmDir returns the actual charm.CharmDir named name.
func (r *Repo) CharmDir(name string) *charm.CharmDir {
	ch, err := charm.ReadCharmDir(r.CharmDirPath(name))
	check(err)
	return ch
}

// ClonedDirPath returns the path to a new copy of the default charm directory
// named name.
func (r *Repo) ClonedDirPath(dst, name string) string {
	return clone(dst, r.CharmDirPath(name))
}

// ClonedDirPath returns the path to a new copy of the default bundle directory
// named name.
func (r *Repo) ClonedBundleDirPath(dst, name string) string {
	return clone(dst, r.BundleDirPath(name))
}

// RenamedClonedDirPath returns the path to a new copy of the default
// charm directory named name, renamed to newName.
func (r *Repo) RenamedClonedDirPath(dst, name, newName string) string {
	dstPath := filepath.Join(dst, newName)
	err := fs.Copy(r.CharmDirPath(name), dstPath)
	check(err)
	return dstPath
}

// ClonedDir returns an actual charm.CharmDir based on a new copy of the charm directory
// named name, in the directory dst.
func (r *Repo) ClonedDir(dst, name string) *charm.CharmDir {
	ch, err := charm.ReadCharmDir(r.ClonedDirPath(dst, name))
	check(err)
	return ch
}

// ClonedURL makes a copy of the charm directory into the new location specified
// by dst. The return value is a URL pointing at the local charm.
func (r *Repo) ClonedURL(dst, series, name string) *charm.URL {
	clone(dst, r.CharmDirPath(name))
	return &charm.URL{
		Schema:   "local",
		Name:     name,
		Revision: -1,
		Series:   series,
	}
}

// CharmArchivePath returns the path to a new charm archive file
// in the directory dst, created from the charm directory named name.
func (r *Repo) CharmArchivePath(dst, name string) string {
	dir := r.CharmDir(name)
	path := filepath.Join(dst, "archive.charm")
	file, err := os.Create(path)
	check(err)
	defer file.Close()
	check(dir.ArchiveTo(file))
	return path
}

// BundleArchivePath returns the path to a new bundle archive file
// in the directory dst, created from the bundle directory named name.
func (r *Repo) BundleArchivePath(dst, name string) string {
	dir := r.BundleDir(name)
	path := filepath.Join(dst, "archive.bundle")
	file, err := os.Create(path)
	check(err)
	defer file.Close()
	check(dir.ArchiveTo(file))
	return path
}

// CharmArchive returns an actual charm.CharmArchive created from a new
// charm archive file created from the charm directory named name, in
// the directory dst.
func (r *Repo) CharmArchive(dst, name string) *charm.CharmArchive {
	ch, err := charm.ReadCharmArchive(r.CharmArchivePath(dst, name))
	check(err)
	return ch
}

// BundleArchive returns an actual charm.BundleArchive created from a new
// bundle archive file created from the bundle directory named name, in
// the directory dst.
func (r *Repo) BundleArchive(dst, name string) *charm.BundleArchive {
	b, err := charm.ReadBundleArchive(r.BundleArchivePath(dst, name))
	check(err)
	return b
}
