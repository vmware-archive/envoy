package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-golang/envoy/internal/middleware"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Credentialer struct{}

func NewCredentialer() *Credentialer {
	return &Credentialer{}
}

func (c Credentialer) Credentials() (string, string) {
	return "username", "password"
}

var _ = Describe("Authenticator", func() {
	Describe("ServeHTTP", func() {
		var wasCalled bool
		var authenticator http.Handler
		var writer *httptest.ResponseRecorder
		var request *http.Request

		BeforeEach(func() {
			var err error
			wasCalled = false
			handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				wasCalled = true
				w.WriteHeader(http.StatusTeapot)
			})
			credentialer := NewCredentialer()
			authenticator = middleware.NewAuthenticator(handler, credentialer)

			writer = httptest.NewRecorder()
			request, err = http.NewRequest("GET", "/foo", nil)
			if err != nil {
				panic(err)
			}
		})

		It("delegates to handler when authentication is valid, but doesn't change the status code", func() {
			request.SetBasicAuth("username", "password")

			authenticator.ServeHTTP(writer, request)

			Expect(wasCalled).To(BeTrue())
			Expect(writer.Code).To(Equal(http.StatusTeapot))
		})

		It("returns a 401 when authentication is not valid", func() {
			authenticator.ServeHTTP(writer, request)

			Expect(wasCalled).To(BeFalse())
			Expect(writer.Code).To(Equal(http.StatusUnauthorized))
		})
	})
})
