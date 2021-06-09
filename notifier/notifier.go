package notifier

import (
	"time"

	"github.com/google/uuid"
)

const (
	Received Status = iota
	FailedIvalidCart
	Fulfilled
)

type (
	Status int

	OrderEvent struct {
		OrderID      uuid.UUID
		Status       Status
		Deliverytime time.Time
	}

	Observer interface {
		OnNotify(OrderEvent)
	}

	Notifier interface {
		Register(Observer)
		Deregister(Observer)
		Notify(OrderEvent)
	}
)

func (e *OrderEvent) EventStatus() string {
	orderStatusMap := map[Status]string{
		Received:         "Received",
		FailedIvalidCart: "Failed Invalid Cart",
		Fulfilled:        "Fulfilled",
	}
	return orderStatusMap[e.Status]
}

func (e *OrderEvent) EstDeliveryTime() string {
	return e.Deliverytime.Format(time.ANSIC)
}

type (
	EventNotifier struct {
		observers map[Observer]struct{}
	}
)

func NewEventNotifier() *EventNotifier {
	enotifier := EventNotifier{}
	enotifier.observers = make(map[Observer]struct{})
	return &enotifier
}

func (en *EventNotifier) Register(l Observer) {
	en.observers[l] = struct{}{}
}

func (en *EventNotifier) Deregister(l Observer) {
	delete(en.observers, l)
}

func (en *EventNotifier) Notify(e OrderEvent) {
	for o := range en.observers {
		o.OnNotify(e)
	}
}
