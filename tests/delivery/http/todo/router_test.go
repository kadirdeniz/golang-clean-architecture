package todo_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kadirdeniz/golang-clean-architecture/internal/delivery/http/todo"
	todoMocks "github.com/kadirdeniz/golang-clean-architecture/tests/mocks/todo"
)

var _ = Describe("Todo Router", func() {
	var (
		app           *fiber.App
		router        todo.Router
		mockHandler   *todoMocks.MockHandler
		ctrl          *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockHandler = todoMocks.NewMockHandler(ctrl)
		router = todo.NewRouter(mockHandler)
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
			var _ todo.Router = router
		})
	})

	Describe("RegisterRoutes", func() {
		It("should register /todos group with POST endpoint", func() {
			router.RegisterRoutes(app)

			stack := app.Stack()
			Expect(stack).ToNot(BeEmpty())

			var todoRouteFound bool
			for _, route := range stack {
				for _, r := range route {
					if r.Path == "/todos/" && r.Method == "POST" {
						todoRouteFound = true
						break
					}
				}
			}
			Expect(todoRouteFound).To(BeTrue())
		})

		It("should register /todos group with PUT endpoint", func() {
			router.RegisterRoutes(app)

			stack := app.Stack()
			Expect(stack).ToNot(BeEmpty())

			var todoUpdateRouteFound bool
			for _, route := range stack {
				for _, r := range route {
					if r.Path == "/todos/:id" && r.Method == "PUT" {
						todoUpdateRouteFound = true
						break
					}
				}
			}
			Expect(todoUpdateRouteFound).To(BeTrue())
		})

		It("should delegate POST requests to handler", func() {
			mockHandler.EXPECT().
				Create(gomock.Any()).
				Return(nil).
				Times(1)

			router.RegisterRoutes(app)

			req := httptest.NewRequest(http.MethodPost, "/todos/", nil)
			_, err := app.Test(req)

			Expect(err).To(BeNil())
		})

		It("should delegate PUT requests to handler", func() {
			mockHandler.EXPECT().
				Update(gomock.Any()).
				Return(nil).
				Times(1)

			router.RegisterRoutes(app)

			req := httptest.NewRequest(http.MethodPut, "/todos/1", nil)
			_, err := app.Test(req)

			Expect(err).To(BeNil())
		})

		It("should register /todos group with DELETE endpoint", func() {
			router.RegisterRoutes(app)

			stack := app.Stack()
			Expect(stack).ToNot(BeEmpty())

			var todoDeleteRouteFound bool
			for _, route := range stack {
				for _, r := range route {
					if r.Path == "/todos/:id" && r.Method == "DELETE" {
						todoDeleteRouteFound = true
						break
					}
				}
			}
			Expect(todoDeleteRouteFound).To(BeTrue())
		})

		It("should register /todos group with GET endpoint", func() {
			router.RegisterRoutes(app)

			stack := app.Stack()
			Expect(stack).ToNot(BeEmpty())

			var todoGetRouteFound bool
			for _, route := range stack {
				for _, r := range route {
					if r.Path == "/todos/:id" && r.Method == "GET" {
						todoGetRouteFound = true
						break
					}
				}
			}
			Expect(todoGetRouteFound).To(BeTrue())
		})

		It("should register /todos group with GET all endpoint", func() {
			router.RegisterRoutes(app)

			stack := app.Stack()
			Expect(stack).ToNot(BeEmpty())

			var todoGetAllRouteFound bool
			for _, route := range stack {
				for _, r := range route {
					if r.Path == "/todos/" && r.Method == "GET" {
						todoGetAllRouteFound = true
						break
					}
				}
			}
			Expect(todoGetAllRouteFound).To(BeTrue())
		})

		It("should delegate GET all requests to handler", func() {
			mockHandler.EXPECT().
				GetAll(gomock.Any()).
				Return(nil).
				Times(1)

			router.RegisterRoutes(app)

			req := httptest.NewRequest(http.MethodGet, "/todos/", nil)
			_, err := app.Test(req)

			Expect(err).To(BeNil())
		})

		It("should delegate GET by ID requests to handler", func() {
			mockHandler.EXPECT().
				GetByID(gomock.Any()).
				Return(nil).
				Times(1)

			router.RegisterRoutes(app)

			req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
			_, err := app.Test(req)

			Expect(err).To(BeNil())
		})

		It("should delegate DELETE requests to handler", func() {
			mockHandler.EXPECT().
				Delete(gomock.Any()).
				Return(nil).
				Times(1)

			router.RegisterRoutes(app)

			req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
			_, err := app.Test(req)

			Expect(err).To(BeNil())
		})
	})

	Describe("Invalid Requests", func() {
		BeforeEach(func() {
			router.RegisterRoutes(app)
		})

		It("should return 404 for unknown paths", func() {
			req := httptest.NewRequest(http.MethodPost, "/unknown", nil)
			resp, err := app.Test(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should return 405 for unsupported methods on /todos", func() {
			req := httptest.NewRequest(http.MethodPatch, "/todos", nil)
			resp, err := app.Test(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
		})

		It("should return 405 for PUT method on /todos", func() {
			req := httptest.NewRequest(http.MethodPut, "/todos", nil)
			resp, err := app.Test(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
		})

		It("should return 405 for DELETE method on /todos", func() {
			req := httptest.NewRequest(http.MethodDelete, "/todos", nil)
			resp, err := app.Test(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
		})
	})
}) 