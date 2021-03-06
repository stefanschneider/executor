package bytefmt_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-golang/bytefmt"
)

var _ = Describe("bytefmt", func() {

	Context("ByteSize", func() {
		It("Prints in the largest possible unit", func() {
			Expect(ByteSize(100 * MEGABYTE)).To(Equal("100M"))
			Expect(ByteSize(uint64(100.5 * MEGABYTE))).To(Equal("100.5M"))
		})
	})

	Context("ToMegabytes", func() {
		It("parses byte amounts with short units (e.g. M, G)", func() {
			var (
				megabytes uint64
				err       error
			)

			megabytes, err = ToMegabytes("5M")
			Expect(megabytes).To(Equal(uint64(5)))
			Expect(err).NotTo(HaveOccurred())

			megabytes, err = ToMegabytes("5m")
			Expect(megabytes).To(Equal(uint64(5)))
			Expect(err).NotTo(HaveOccurred())

			megabytes, err = ToMegabytes("2G")
			Expect(megabytes).To(Equal(uint64(2 * 1024)))
			Expect(err).NotTo(HaveOccurred())

			megabytes, err = ToMegabytes("3T")
			Expect(megabytes).To(Equal(uint64(3 * 1024 * 1024)))
			Expect(err).NotTo(HaveOccurred())
		})

		It("parses byte amounts with long units (e.g MB, GB)", func() {
			var (
				megabytes uint64
				err       error
			)

			megabytes, err = ToMegabytes("5MB")
			Expect(megabytes).To(Equal(uint64(5)))
			Expect(err).NotTo(HaveOccurred())

			megabytes, err = ToMegabytes("5mb")
			Expect(megabytes).To(Equal(uint64(5)))
			Expect(err).NotTo(HaveOccurred())

			megabytes, err = ToMegabytes("2GB")
			Expect(megabytes).To(Equal(uint64(2 * 1024)))
			Expect(err).NotTo(HaveOccurred())

			megabytes, err = ToMegabytes("3TB")
			Expect(megabytes).To(Equal(uint64(3 * 1024 * 1024)))
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error when the unit is missing", func() {
			_, err := ToMegabytes("5")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unit of measurement"))
		})

		It("returns an error when the unit is unrecognized", func() {
			_, err := ToMegabytes("5MBB")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unit of measurement"))
		})

		It("allows whitespace before and after the value", func() {
			megabytes, err := ToMegabytes("\t\n\r 5MB ")
			Expect(megabytes).To(Equal(uint64(5)))
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error for negative values", func() {
			_, err := ToMegabytes("-5MB")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unit of measurement"))
		})

		It("returns an error for zero values", func() {
			_, err := ToMegabytes("0TB")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unit of measurement"))
		})
	})
})
