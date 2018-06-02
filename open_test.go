package mmf_test

import (
	"os"

	"github.com/draganm/mmf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Open", func() {
	var e error
	var f *mmf.File
	JustBeforeEach(func() {
		f, e = mmf.Open(fileName)
	})

	Context("When the file does not exist", func() {

		It("should not return an error", func() {
			Expect(e).ToNot(HaveOccurred())
		})

		It("Should return a non-nil file", func() {
			Expect(f).ToNot(BeNil())
		})

		It("Should create the file on the filesystem", func() {
			fs, err := os.Stat(fileName)
			Expect(err).ToNot(HaveOccurred())
			Expect(fs.IsDir()).To(BeFalse())
			Expect(fs.Size()).To(Equal(int64(0)))
		})

	})

	Context("When the file exists and is empty", func() {

		BeforeEach(func() {
			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0700)
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Close()).To(Succeed())
		})

		It("should not return an error", func() {
			Expect(e).ToNot(HaveOccurred())
		})

		It("Should return a non-nil file", func() {
			Expect(f).ToNot(BeNil())
		})

		It("Should have size 0 mmap", func() {
			Expect(len(f.MMap)).To(Equal(0))
		})

	})

	Context("When the file exists and has some data", func() {

		BeforeEach(func() {
			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0700)
			Expect(err).ToNot(HaveOccurred())
			_, err = f.Write([]byte{1, 2, 3})
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Close()).To(Succeed())
		})

		It("should not return an error", func() {
			Expect(e).ToNot(HaveOccurred())
		})

		It("Should return a non-nil file", func() {
			Expect(f).ToNot(BeNil())
		})

		It("Should mmap the file's data", func() {
			Expect([]byte(f.MMap)).To(Equal([]byte{1, 2, 3}))
		})

	})

})
