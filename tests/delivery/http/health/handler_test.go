package health_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kadirdeniz/golang-clean-architecture/internal/delivery/http/health"
)

var _ = Describe("Health Handler", func() {
	var (
		handler health.Handler
		app     *fiber.App
	)

	BeforeEach(func() {
		handler = health.NewHandler()
		app = fiber.New()
		app.Get("/health", handler.Health)
	})

	AfterEach(func() {
		if app != nil {
			app.Shutdown()
		}
	})

	Describe("NewHandler", func() {
		It("should create a new handler instance", func() {
			Expect(handler).ToNot(BeNil())
		})

		It("should implement Handler interface", func() {
			var _ health.Handler = handler
		})
	})

	Describe("Health Method", func() {
		It("should return 200 status code", func() {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			resp, err := app.Test(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should return OK response body", func() {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			resp, err := app.Test(req)

			Expect(err).To(BeNil())
			
			body, err := io.ReadAll(resp.Body)
			Expect(err).To(BeNil())
			Expect(string(body)).To(Equal("OK"))
		})

		It("should set correct content type", func() {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			resp, err := app.Test(req)

			Expect(err).To(BeNil())
			Expect(resp.Header.Get("Content-Type")).To(ContainSubstring("text/plain"))
		})
	})
}) 