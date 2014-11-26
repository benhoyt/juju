// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter

import (
	"os"
	"path/filepath"

	"github.com/juju/errors"
	"github.com/juju/utils/symlink"

	"github.com/juju/juju/juju/names"
	"github.com/juju/juju/worker/uniter/jujuc"
)

// EnsureJujucSymlinks creates a symbolic link to jujuc within dir for each
// hook command. If the commands already exist, this operation does nothing.
// If dir is a symbolic link, it will be dereferenced first.
func EnsureJujucSymlinks(dir string) (err error) {
	logger.Infof("ensure jujuc symlinks in %s", dir)
	defer func() {
		if err != nil {
			err = errors.Annotatef(err, "cannot initialize hook commands in %q", dir)
		}
	}()
	st, err := os.Lstat(dir)
	if err != nil {
		return err
	}
	// NOTE: need to work out how to do this on windows.
	if st.Mode()&os.ModeSymlink != 0 {
		link, err := symlink.Read(dir)
		if err != nil {
			return err
		}
		if !filepath.IsAbs(link) {
			logger.Infof("%s is relative", link)
			link = filepath.Join(filepath.Dir(dir), link)
		}
		dir = link
		logger.Infof("was a symlink, now looking at %s", dir)
	}

	jujudPath := filepath.Join(dir, names.Jujud)
	logger.Debugf("jujud path %s", jujudPath)
	for _, name := range jujuc.CommandNames() {
		// The link operation fails when the target already exists,
		// so this is a no-op when the command names already
		// exist.
		err := symlink.New(jujudPath, filepath.Join(dir, name))
		if err != nil && !os.IsExist(err) {
			return err
		}
	}
	return nil
}
