package mmf_test

import (
	"os"

	"github.com/draganm/mmf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Write", func() {
	var e error
	var f *mmf.File
	JustBeforeEach(func() {
		f, e = mmf.Open(fileName)
	})

	AfterEach(func() {
		Expect(f.Close()).To(Succeed())
	})

	JustBeforeEach(func() {
		e = f.Append([]byte{4, 5, 6})
	})

	Context("When the file did not exist before opening", func() {

		It("should not return error", func() {
			Expect(e).ToNot(HaveOccurred())
		})

		It("Should extend mmap to contain the appended data", func() {
			Expect([]byte(f.MMap)).To(Equal([]byte{4, 5, 6}))
		})

	})

	Context("When the file exists and is empty", func() {

		BeforeEach(func() {
			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0700)
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Close()).To(Succeed())
		})

		It("should not return error", func() {
			Expect(e).ToNot(HaveOccurred())
		})

		It("Should extend mmap to contain the appended data", func() {
			Expect([]byte(f.MMap)).To(Equal([]byte{4, 5, 6}))
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

		It("should not return error", func() {
			Expect(e).ToNot(HaveOccurred())
		})

		It("Should extend mmap to contain the appended data", func() {
			Expect([]byte(f.MMap)).To(Equal([]byte{1, 2, 3, 4, 5, 6}))
		})

	})

	Context("When the file exists and has 4k of data", func() {

		BeforeEach(func() {
			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0700)
			Expect(err).ToNot(HaveOccurred())
			data := make([]byte, 4096)
			for i := 0; i < 4096; i++ {
				data[i] = byte(i % 255)
			}
			_, err = f.Write(data)
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Close()).To(Succeed())
		})

		It("should not return error", func() {
			Expect(e).ToNot(HaveOccurred())
		})

		It("Should extend mmap to contain the appended data", func() {
			Expect([]byte(f.MMap[4096:])).To(Equal([]byte{4, 5, 6}))
		})

	})

	Context("When the file exists and has 128k of data", func() {

		BeforeEach(func() {
			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0700)
			Expect(err).ToNot(HaveOccurred())
			data := make([]byte, 128*1024)
			for i := 0; i < 128*1024; i++ {
				data[i] = byte(i % 255)
			}
			_, err = f.Write(data)
			Expect(err).ToNot(HaveOccurred())
			Expect(f.Close()).To(Succeed())
		})

		It("should not return error", func() {
			Expect(e).ToNot(HaveOccurred())
		})

		It("Should extend mmap to contain the appended data", func() {
			Expect([]byte(f.MMap[128*1024:])).To(Equal([]byte{4, 5, 6}))
		})

	})

})
