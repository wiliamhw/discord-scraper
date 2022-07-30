package util

import (
	"os"

	"github.com/karrick/godirwalk"
)

func PruneEmptyDirectories(osDirname string) (int, error) {
	var count int

	err := godirwalk.Walk(osDirname, &godirwalk.Options{
		Unsorted: true,
		Callback: func(_ string, _ *godirwalk.Dirent) error {
			// no-op while diving in; all the fun happens in PostChildrenCallback
			return nil
		},
		PostChildrenCallback: func(osPathname string, _ *godirwalk.Dirent) error {
			s, err := godirwalk.NewScanner(osPathname)
			if err != nil {
				return err
			}

			// Attempt to read only the first directory entry. Remember that
			// Scan skips both "." and ".." entries.
			hasAtLeastOneChild := s.Scan()

			// If error reading from directory, wrap up and return.
			if err := s.Err(); err != nil {
				return err
			}

			if hasAtLeastOneChild {
				return nil // do not remove directory with at least one child
			}
			if osPathname == osDirname {
				return nil // do not remove directory that was provided top-level directory
			}

			err = os.Remove(osPathname)
			if err == nil {
				count++
			}
			return err
		},
	})

	return count, err
}
