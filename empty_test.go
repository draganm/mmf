package mmf_test

import (
	"os"

	"github.com/draganm/mmf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Empty", func() {

	var f *mmf.File
	var e error

	BeforeEach(func() {
		var err error
		f, err = mmf.Open(fileName)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(f.Close()).To(Succeed())
	})

	JustBeforeEach(func() {
		e = f.Empty()
	})

	Context("When the file is empty", func() {

		It("should not return an error", func() {
			Expect(e).ToNot(HaveOccurred())
		})

	})

	Context("When the file exists and has some data", func() {

		BeforeEach(func() {
			Expect(f.Append([]byte{1, 2, 3})).To(Succeed())
		})

		It("should not return an error", func() {
			Expect(e).ToNot(HaveOccurred())
		})

		It("Should remove mmap", func() {
			Expect(len(f.MMap)).To(Equal(0))
		})

		It("Should truncate file to zero size", func() {
			s, err := os.Stat(fileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(s.Size()).To(Equal(int64(0)))
		})

	})

})
