package mmf_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/draganm/mmf"
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

var _ = Describe("Open", func() {
	var err error
	var f *mmf.File
	JustBeforeEach(func() {
		f, err = mmf.Open(fileName)
	})

	Context("When the file does not exist", func() {

		It("should not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("Should return a non nil file", func() {
			Expect(f).ToNot(BeNil())
		})

		It("Should create the file on the filesystem", func() {
			fs, err := os.Stat(fileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(fs.IsDir()).To(BeFalse())
			Expect(fs.Size()).To(Equal(int64(0)))
		})

	})
})
