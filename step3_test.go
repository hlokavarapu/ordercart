package main

import (
	"context"
	"ordercart/notifier"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testObserver struct {
	mock.Mock
}

func (to *testObserver) OnNotify(e notifier.OrderEvent) {
	to.Called()
}

func TestObserver(t *testing.T) {
	s := newTestOrderCartServer(t, []byte(testCatalogStep1))

	obs := testObserver{}
	s.eNotifier.Register(&obs)
	obs.On("OnNotify").Return()

	request := createGetOrderCostRequest([]string{"Apple", "Apple", "Orange", "Apple"})
	_, err := s.GetOrderCost(context.Background(), request)
	assert.NoError(t, err)
	obs.AssertNumberOfCalls(t, "OnNotify", 2)
}

func TestObserverInvalidItem(t *testing.T) {
	s := newTestOrderCartServer(t, []byte(testCatalogStep1))

	obs := testObserver{}
	s.eNotifier.Register(&obs)
	obs.On("OnNotify").Return()

	request := createGetOrderCostRequest([]string{"InvalidItem"})
	_, err := s.GetOrderCost(context.Background(), request)
	assert.Error(t, err)
	obs.AssertNumberOfCalls(t, "OnNotify", 2)
}
