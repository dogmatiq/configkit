package static

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dogmatiq/aureus"
	"github.com/dogmatiq/configkit"
)

func TestFromPackages(t *testing.T) {
	aureus.Run(
		t,
		func(w io.Writer, in aureus.Content, out aureus.ContentMetaData) error {
			pkg := strings.TrimSuffix(
				filepath.Base(in.File),
				filepath.Ext(in.File),
			)

			// Make a temporary directory to write the Go source code.
			//
			// The name is based on the input file name rather than using a
			// random temporary directory, otherwise the test output would be
			// non-deterministic.
			//
			// Additionally, creating the directory within the repository allows
			// the test code to use this repo's go.mod file, ensuring the
			// statically analyzed code uses the same versions of Dogma, etc.
			dir := filepath.Join(
				filepath.Dir(in.File),
				pkg,
			)
			if err := os.MkdirAll(dir, 0700); err != nil {
				return err
			}
			defer os.RemoveAll(dir)

			if err := os.WriteFile(
				filepath.Join(dir, "main.go"),
				[]byte(in.Data),
				0600,
			); err != nil {
				return err
			}

			apps := FromDir(dir)

			if len(apps) == 0 {
				_, err := io.WriteString(w, "(no applications found)\n")
				return err
			}

			noise := []string{
				"github.com/dogmatiq/configkit/static/testdata/aureus/" + pkg + ".",
				"github.com/dogmatiq/enginekit/enginetest/stubs.",
			}

			for i, app := range apps {
				s := configkit.ToString(app)
				for _, p := range noise {
					s = strings.ReplaceAll(s, p, "")
				}

				if i > 0 {
					s = "\n" + s
				}
				if _, err := io.WriteString(w, s); err != nil {
					return err
				}
			}

			return nil
		},
	)
}
