package mmf_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMmf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mmf Suite")
}

var dir string
var fileName string
var _ = BeforeEach(func() {
	var err error
	dir, err = ioutil.TempDir("", "")
	Expect(err).ToNot(HaveOccurred())
	fileName = filepath.Join(dir, "some-file")
})

var _ = AfterEach(func() {
	if dir != "" {
		Expect(os.RemoveAll(dir)).To(Succeed())
	}
})
