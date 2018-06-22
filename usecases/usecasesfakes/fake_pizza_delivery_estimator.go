// This file was generated by counterfeiter
package usecasesfakes

import (
	"sync"
	"time"

	"github.com/tjarratt/go-best-practices/domain"
	"github.com/tjarratt/go-best-practices/usecases"
)

type FakePizzaDeliveryEstimator struct {
	EstimatedDeliveryTimeStub        func(domain.Pizza) time.Duration
	estimatedDeliveryTimeMutex       sync.RWMutex
	estimatedDeliveryTimeArgsForCall []struct {
		arg1 domain.Pizza
	}
	estimatedDeliveryTimeReturns struct {
		result1 time.Duration
	}
}

func (fake *FakePizzaDeliveryEstimator) EstimatedDeliveryTime(arg1 domain.Pizza) time.Duration {
	fake.estimatedDeliveryTimeMutex.Lock()
	fake.estimatedDeliveryTimeArgsForCall = append(fake.estimatedDeliveryTimeArgsForCall, struct {
		arg1 domain.Pizza
	}{arg1})
	fake.estimatedDeliveryTimeMutex.Unlock()
	if fake.EstimatedDeliveryTimeStub != nil {
		return fake.EstimatedDeliveryTimeStub(arg1)
	} else {
		return fake.estimatedDeliveryTimeReturns.result1
	}
}

func (fake *FakePizzaDeliveryEstimator) EstimatedDeliveryTimeCallCount() int {
	fake.estimatedDeliveryTimeMutex.RLock()
	defer fake.estimatedDeliveryTimeMutex.RUnlock()
	return len(fake.estimatedDeliveryTimeArgsForCall)
}

func (fake *FakePizzaDeliveryEstimator) EstimatedDeliveryTimeArgsForCall(i int) domain.Pizza {
	fake.estimatedDeliveryTimeMutex.RLock()
	defer fake.estimatedDeliveryTimeMutex.RUnlock()
	return fake.estimatedDeliveryTimeArgsForCall[i].arg1
}

func (fake *FakePizzaDeliveryEstimator) EstimatedDeliveryTimeReturns(result1 time.Duration) {
	fake.EstimatedDeliveryTimeStub = nil
	fake.estimatedDeliveryTimeReturns = struct {
		result1 time.Duration
	}{result1}
}

var _ usecases.PizzaDeliveryEstimator = new(FakePizzaDeliveryEstimator)
