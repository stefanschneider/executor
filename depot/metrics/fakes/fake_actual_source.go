// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/executor/depot/metrics"
	garden_api "github.com/cloudfoundry-incubator/garden/api"
)

type FakeActualSource struct {
	ContainersStub        func() ([]garden_api.Container, error)
	containersMutex       sync.RWMutex
	containersArgsForCall []struct{}
	containersReturns     struct {
		result1 []garden_api.Container
		result2 error
	}
}

func (fake *FakeActualSource) Containers() ([]garden_api.Container, error) {
	fake.containersMutex.Lock()
	defer fake.containersMutex.Unlock()
	fake.containersArgsForCall = append(fake.containersArgsForCall, struct{}{})
	if fake.ContainersStub != nil {
		return fake.ContainersStub()
	} else {
		return fake.containersReturns.result1, fake.containersReturns.result2
	}
}

func (fake *FakeActualSource) ContainersCallCount() int {
	fake.containersMutex.RLock()
	defer fake.containersMutex.RUnlock()
	return len(fake.containersArgsForCall)
}

func (fake *FakeActualSource) ContainersReturns(result1 []garden_api.Container, result2 error) {
	fake.containersReturns = struct {
		result1 []garden_api.Container
		result2 error
	}{result1, result2}
}

var _ metrics.ActualSource = new(FakeActualSource)