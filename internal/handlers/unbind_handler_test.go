package handlers_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-golang/envoy/domain"
	"github.com/pivotal-golang/envoy/internal/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Unbinder struct {
	WasCalledWith domain.UnbindRequest
	UnbindError   error
}

func NewUnbinder() *Unbinder {
	return &Unbinder{}
}

func (f *Unbinder) Unbind(req domain.UnbindRequest) error {
	f.WasCalledWith = req
	return f.UnbindError
}

var _ = Describe("UnbindHandler", func() {
	var unbinder *Unbinder
	var handler handlers.UnbindHandler

	BeforeEach(func() {
		unbinder = NewUnbinder()
		handler = handlers.NewUnbindHandler(unbinder)
	})

	It("calls the binder Unbind method with the correct values", func() {
		writer := httptest.NewRecorder()

		url := fmt.Sprintf("%s?plan_id=%s&service_id=%s",
			"/v2/service_instances/service-instance-id/service_bindings/service-binding-id",
			"the-1gb-plan-id",
			"the-sshfs-service-id")
		request, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(writer, request)

		Expect(unbinder.WasCalledWith).To(Equal(domain.UnbindRequest{
			BindingID:  "service-binding-id",
			InstanceID: "service-instance-id",
			ServiceID:  "the-sshfs-service-id",
			PlanID:     "the-1gb-plan-id",
		}))
	})

	Context("when the unbinder succeeds", func() {
		It("returns a 200 status code with an empty JSON body", func() {
			writer := httptest.NewRecorder()
			request, err := http.NewRequest("DELETE",
				"/v2/service_instances/service-instance-id/service_bindings/service-binding-id?plan_id=some-plan-id&service_id=some-service-id",
				nil)
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON("{}"))
		})
	})

	Context("when the binding does not exist", func() {
		It("returns 410 Gone with empty JSON body", func() {
			writer := httptest.NewRecorder()
			request, err := http.NewRequest("DELETE",
				"/v2/service_instances/service-instance-id/service_bindings/a-non-existent-service-binding-id?plan_id=some-plan-id&service_id=some-service-id",
				nil)
			if err != nil {
				panic(err)
			}

			unbinder.UnbindError = domain.ServiceBindingNotFoundError(
				("that binding doesn't exist!"))

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusGone))

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON("{}"))

		})
	})

	Context("when the unbinder fails", func() {
		It("returns a 500 error with the message", func() {
			writer := httptest.NewRecorder()
			request, err := http.NewRequest("DELETE", "/v2/service_instances/service-instance-id/service_bindings/a-service-binding-id?plan_id=some-plan-id&service_id=some-service-id",
				nil)
			if err != nil {
				panic(err)
			}

			unbinder.UnbindError = errors.New("my database failed somehow!")

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}

			Expect(body).To(MatchJSON(`{"description": "my database failed somehow!"}`))
		})
	})

})