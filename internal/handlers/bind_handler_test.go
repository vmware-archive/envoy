package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/pivotal-cf-experimental/envoy/domain"
	"github.com/pivotal-cf-experimental/envoy/internal/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Binder struct {
	WasCalledWith  domain.BindRequest
	Credentials    domain.BindingCredentials
	Error          error
	SyslogDrainURL string
}

func NewBinder() *Binder {
	return &Binder{}
}

func (b *Binder) Bind(binding domain.BindRequest) (domain.BindResponse, error) {
	b.WasCalledWith = binding

	return domain.BindResponse{
		Credentials:    b.Credentials,
		SyslogDrainURL: b.SyslogDrainURL,
	}, b.Error
}

var _ = Describe("BindHandler", func() {
	var handler handlers.BindHandler
	var binder *Binder

	BeforeEach(func() {
		binder = NewBinder()
		handler = handlers.NewBindHandler(binder)
	})

	It("calls the binder Bind method with the correct values", func() {
		writer := httptest.NewRecorder()
		reqBody, err := json.Marshal(map[string]string{
			"service_id": "service-id",
			"plan_id":    "plan-id",
			"app_guid":   "app-guid",
		})
		if err != nil {
			panic(err)
		}

		request, err := http.NewRequest("PUT", "/v2/service_instances/service-instance-id/service_bindings/service-binding-id", bytes.NewBuffer(reqBody))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(writer, request)

		Expect(binder.WasCalledWith).To(Equal(domain.BindRequest{
			BindingID:  "service-binding-id",
			InstanceID: "service-instance-id",
			ServiceID:  "service-id",
			PlanID:     "plan-id",
			AppGUID:    "app-guid",
		}))
	})

	It("returns a 201 status code with an empty JSON body", func() {
		writer := httptest.NewRecorder()
		reqBody, err := json.Marshal(map[string]string{
			"service_id": "service-id",
			"plan_id":    "plan-id",
			"app_guid":   "app-guid",
		})
		if err != nil {
			panic(err)
		}

		request, err := http.NewRequest("PUT", "/v2/service_instances/service-instance-id/service_bindings/service-binding-id", bytes.NewBuffer(reqBody))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(writer, request)

		Expect(writer.Code).To(Equal(http.StatusCreated))
		body, err := ioutil.ReadAll(writer.Body)
		if err != nil {
			panic(err)
		}
		Expect(body).To(MatchJSON("{}"))
	})

	Context("when binding credentials are provided", func() {
		BeforeEach(func() {
			binder.Credentials = domain.BindingCredentials{
				"uri":      "mysql://mysqluser:pass@mysqlhost:3306/dbname",
				"username": "mysqluser",
				"password": "pass",
				"host":     "mysqlhost",
				"port":     3306,
				"database": "dbname",
			}
		})

		It("returns the credentials in the response body", func() {
			writer := httptest.NewRecorder()
			reqBody, err := json.Marshal(map[string]string{
				"service_id": "service-id",
				"plan_id":    "plan-id",
				"app_guid":   "app-guid",
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", "/v2/service_instances/service-instance-id/service_bindings/service-binding-id", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusCreated))
			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{
				"credentials": {
					"uri":      "mysql://mysqluser:pass@mysqlhost:3306/dbname",
					"username": "mysqluser",
					"password": "pass",
					"host":     "mysqlhost",
					"port":     3306,
					"database": "dbname"
				}
			}`))
		})
	})

	Context("when binding syslog drain URL is provided", func() {
		BeforeEach(func() {
			binder.SyslogDrainURL = "syslog://something"
		})

		It("returns the syslog drain URL in the response body", func() {
			writer := httptest.NewRecorder()
			reqBody, err := json.Marshal(map[string]string{
				"service_id": "service-id",
				"plan_id":    "plan-id",
				"app_guid":   "app-guid",
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", "/v2/service_instances/service-instance-id/service_bindings/service-binding-id", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusCreated))
			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{
				"syslog_drain_url": "syslog://something"
			}`))
		})
	})

	Context("when there is a binding failure", func() {
		BeforeEach(func() {
			binder.Error = errors.New("BANG!")
		})

		It("returns a 500 and the error as the body", func() {
			writer := httptest.NewRecorder()
			reqBody, err := json.Marshal(map[string]string{
				"service_id": "my-service-id",
				"plan_id":    "my-plan-id",
				"app_guid":   "my-app-guid",
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", "/v2/service_instances/instance-guid/service_bindings/binding-guid", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{"description":"BANG!"}`))
		})
	})

	Context("when the service binding has already been provisioned", func() {
		BeforeEach(func() {
			binder.Error = domain.ServiceBindingAlreadyExistsError
		})

		It("returns a 409 and the error message", func() {
			writer := httptest.NewRecorder()
			reqBody, err := json.Marshal(map[string]string{
				"service_id": "my-service-id",
				"plan_id":    "my-plan-id",
				"app_guid":   "my-app-guid",
			})
			if err != nil {
				panic(err)
			}

			request, err := http.NewRequest("PUT", "/v2/service_instances/instance-guid/service_bindings/binding-guid", bytes.NewBuffer(reqBody))
			if err != nil {
				panic(err)
			}

			handler.ServeHTTP(writer, request)

			Expect(writer.Code).To(Equal(http.StatusConflict))

			body, err := ioutil.ReadAll(writer.Body)
			if err != nil {
				panic(err)
			}
			Expect(body).To(MatchJSON(`{}`))
		})
	})
})
