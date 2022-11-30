package cmd

import (
	"errors"
	"fmt"
	"os"
)

// returns an error if the file already exists (or if something wromg happened while check it)
func checkOutputFile(filename string) error {
	_, err := os.Stat(filename)
	switch {
	case errors.Is(err, os.ErrNotExist) || force:
		// it's what we want
		return nil
	case err != nil:
		// something unexpected happened
		return err
	default:
		// it's not what we want (we don't want to overwrite an existing file!)
		return fmt.Errorf("%s already exists", filename)
	}

}
