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
		func(_ *testing.T, in aureus.Input, out aureus.Output) error {
			// Make a temporary directory to write the Go source code.
			//
			// Creating the directory within the repository (as opposed to in
			// the system temp directory) allows the test code to use this
			// repo's go.mod file, ensuring the statically analyzed code uses
			// the same versions of Dogma, etc.
			dir, err := os.MkdirTemp("testdata/aureus", "")
			if err != nil {
				return err
			}
			defer os.RemoveAll(dir)

			main, err := os.Create(filepath.Join(dir, "main.go"))
			if err != nil {
				return err
			}
			defer main.Close()

			if _, err := io.Copy(main, in); err != nil {
				return err
			}

			apps := FromDir(dir)

			if len(apps) == 0 {
				_, err := io.WriteString(out, "(no applications found)\n")
				return err
			}

			noise := []string{
				"github.com/dogmatiq/configkit/static/testdata/aureus/" + filepath.Base(dir) + ".",
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
				if _, err := io.WriteString(out, s); err != nil {
					return err
				}
			}

			return nil
		},
	)
}
