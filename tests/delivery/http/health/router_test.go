package health_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kadirdeniz/golang-clean-architecture/internal/delivery/http/health"
	healthMocks "github.com/kadirdeniz/golang-clean-architecture/tests/mocks/health"
)

// TestHealthRouter is handled by the main health suite

var _ = Describe("Health Router", func() {
	var (
		app           *fiber.App
		router        health.Router
		mockHandler   *healthMocks.MockHandler
		ctrl          *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockHandler = healthMocks.NewMockHandler(ctrl)
		router = health.NewRouter(mockHandler)
		app = fiber.New()
	})

	AfterEach(func() {
		if ctrl != nil {
			ctrl.Finish()
		}
		if app != nil {
			app.Shutdown()
		}
	})

	Describe("NewRouter", func() {
		It("should create a new router instance", func() {
			Expect(router).ToNot(BeNil())
		})

		It("should implement Router interface", func() {
			var _ health.Router = router
		})
	})

	Describe("RegisterRoutes", func() {
		It("should register /health endpoint", func() {
			router.RegisterRoutes(app)

			stack := app.Stack()
			Expect(stack).ToNot(BeEmpty())

			var healthRouteFound bool
			for _, route := range stack {
				for _, r := range route {
					if r.Path == "/health" && r.Method == "GET" {
						healthRouteFound = true
						break
					}
				}
			}
			Expect(healthRouteFound).To(BeTrue())
		})

		It("should delegate requests to handler", func() {
			mockHandler.EXPECT().
				Health(gomock.Any()).
				Return(nil).
				Times(1)

			router.RegisterRoutes(app)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			_, err := app.Test(req)

			Expect(err).To(BeNil())
		})
	})

	Describe("Invalid Requests", func() {
		BeforeEach(func() {
			router.RegisterRoutes(app)
		})

		It("should return 404 for unknown paths", func() {
			req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
			resp, err := app.Test(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should return 405 for unsupported methods", func() {
			req := httptest.NewRequest(http.MethodPost, "/health", nil)
			resp, err := app.Test(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
		})
	})
}) 