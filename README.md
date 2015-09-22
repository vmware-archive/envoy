# A Cautionary Tale of Two Broker API Packages

We've been told by the powers that be that you shouldn't be using this repo and instead should be using an alternative implementation here:

[CF Service Broker](https://github.com/pivotal-cf/brokerapi)

# Envoy
[![GoDoc](https://godoc.org/github.com/pivotal-cf-experimental/envoy?status.svg)](https://godoc.org/github.com/pivotal-cf-experimental/envoy)

A CloudFoundry Service Broker front-end

Provides a RESTful [Service Broker](http://docs.cloudfoundry.org/services/api.html) for Cloud Controller to consume.

Here is a rather fanciful example of what one might implement for a service broker. It shows how to integrate your existing service with the service broker API using `envoy`.

```go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/pivotal-cf-experimental/envoy"
	"github.com/pivotal-cf-experimental/envoy/domain"
)

func main() {
	broker := Broker{}
	handler := envoy.NewBrokerHandler(broker)
	log.Fatalln(http.ListenAndServe(":0", handler))
}

type Broker struct {
	randomNumbers map[string]int
}

func (b Broker) Catalog() domain.Catalog {
	return domain.Catalog{
		Services: []domain.Service{
			{
				ID:          "64927e6f-b4ae-4ecb-b183-fe73f08339b8",
				Name:        "numbers",
				Description: "Picks a random large number and makes it available to all apps bound to the service.",
				Bindable:    true,
				Plans: []domain.Plan{
					{
						ID:          "c95e94de-b241-44be-96d1-dfef64a50e6a",
						Name:        "free",
						Description: "Provides a number at no cost.",
						Free:        domain.FreeTrue,
					},
				},
			},
		},
	}
}

func (b Broker) Credentials() (string, string) {
	return "username", "password"
}

func (b Broker) Provision(request domain.ProvisionRequest) (domain.ProvisionResponse, error) {
	if _, ok := b.randomNumbers[request.InstanceID]; ok {
		return domain.ProvisionResponse{}, domain.ServiceInstanceAlreadyExistsError("this instance already exists")
	}

	b.randomNumbers[request.InstanceID] = rand.Intn(10000)

	return domain.ProvisionResponse{
		DashboardURL: fmt.Sprintf("https://numbers.example.com/dashboard/%s", request.InstanceID),
	}, nil
}

func (b Broker) Bind(request domain.BindRequest) (domain.BindResponse, error) {
	number, ok := b.randomNumbers[request.InstanceID]

	if !ok {
		return domain.BindResponse{}, domain.ServiceInstanceNotFoundError("could not find this service instance")
	}

	return domain.BindResponse{
		Credentials: domain.BindingCredentials{
			"number": number,
		},
	}, nil
}

func (b Broker) Unbind(request domain.UnbindRequest) error {
	_, ok := b.randomNumbers[request.InstanceID]

	if !ok {
		return domain.ServiceInstanceNotFoundError("could not find this service instance")
	}

	return nil
}

func (b Broker) Deprovision(request domain.DeprovisionRequest) error {
	_, ok := b.randomNumbers[request.InstanceID]

	if !ok {
		return domain.ServiceInstanceNotFoundError("could not find this service instance")
	}

	delete(b.randomNumbers, request.InstanceID)

	return nil
}
```
