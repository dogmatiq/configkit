package typename

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

// Declare declares a unit test-suite for functions "Of()" defined in "goreflect"
// and "gotypes" packages.
func Declare(dir string, fn func(string) string) {
	d, err := os.Open(dir)
	if err != nil {
		panic(err)
	}

	infos, err := d.Readdir(-1)
	d.Close()
	if err != nil {
		panic(err)
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}

		name := info.Name()
		p := filepath.Join(dir, name)

		ginkgo.It(
			fmt.Sprintf("matches the type defined in file %s", p),
			func() {
				bb, err := ioutil.ReadFile(p)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				expected := strings.TrimSpace(string(bb))
				actual := fn(name)

				gomega.Expect(actual).To(gomega.Equal(expected))
			},
		)
	}
}
