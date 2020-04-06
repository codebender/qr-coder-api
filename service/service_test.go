package service_test

import (
	"context"
	"testing"

	"github.com/codebender/qrcode-api/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	qrcode "github.com/skip2/go-qrcode"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("QR Code Service", func() {
	Describe("Generate", func() {
		var (
			subject service.Service
			ctx     context.Context
		)

		BeforeEach(func() {
			subject = service.NewService()
			ctx = context.Background()
		})

		It("returns a generated qr code png's bytes", func() {
			testData := "testing123"
			qrCode, err := subject.Generate(ctx, testData)
			Expect(err).ToNot(HaveOccurred())

			expectedPNG, err := qrcode.Encode(testData, qrcode.Medium, 256)
			Expect(err).ToNot(HaveOccurred())

			Expect(qrCode).To(Equal(expectedPNG))

		})

		It("returns an error when no data is give to encode", func() {
			_, err := subject.Generate(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Data is required"))
		})
	})
})
