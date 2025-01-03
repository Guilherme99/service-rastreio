package internal

import (
	"time"
)

type RouteCreatedEvent struct {
	EventName  string       `json:"event"`
	RouteID    string       `json:"id"`
	Distance   int          `json:"distance"`
	Directions []Directions `json:"directions"`
}

type FreightCalculatedEvent struct {
	EventName string  `json:"event"`
	RouteID   string  `json:"route_id"`
	Amount    float64 `json:"amount"`
}

type DeliveryStartedEvent struct {
	EventName string `json:"event"`
	RouteID   string `json:"route_id"`
}

type DriverMovedEvent struct {
	EventName string  `json:"event"`
	RouteID   string  `json:"route_id"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
}

func NewRouteCreatedEvent(routeID string, distance int, directions []Directions) *RouteCreatedEvent {
	return &RouteCreatedEvent{
		EventName:  "RouteCreated",
		RouteID:    routeID,
		Distance:   distance,
		Directions: directions,
	}
}

func RouteCreatedHandler(event *RouteCreatedEvent, routeService *RouteService) (*FreightCalculatedEvent, error) {
	route := NewRoute(event.RouteID, event.Distance, event.Directions)
	routeCreated, err := routeService.CreateRoute(*route)
	if err != nil {
		return nil, err
	}
	freightCalculatedEvent := NewFreightCalculatedEvent(routeCreated.ID, routeCreated.FreightPrice)
	return freightCalculatedEvent, nil
}

func NewFreightCalculatedEvent(routeID string, amount float64) *FreightCalculatedEvent {
	return &FreightCalculatedEvent{
		EventName: "FreightCalculated",
		RouteID:   routeID,
		Amount:    amount,
	}
}

func NewDeliveryStartedEvent(routeID string) *DeliveryStartedEvent {
	return &DeliveryStartedEvent{
		EventName: "DeliveryStarted",
		RouteID:   routeID,
	}
}

func NewDriverMovedEvent(routeID string, lat float64, lng float64) *DriverMovedEvent {
	return &DriverMovedEvent{
		EventName: "DriverMoved",
		RouteID:   routeID,
		Lat:       lat,
		Lng:       lng,
	}
}

func DeliveryStartedHandler(event *DeliveryStartedEvent, routeService *RouteService, ch chan *DriverMovedEvent) error {
	route, err := routeService.GetRoute(event.RouteID)
	if err != nil {
		return err
	}

	driverMovedEvent := NewDriverMovedEvent(route.ID, 0, 0)
	for _, direction := range route.Directions {
		driverMovedEvent.RouteID = route.ID
		driverMovedEvent.Lat = direction.Lat
		driverMovedEvent.Lng = direction.Lng
		time.Sleep(time.Second)
		ch <- driverMovedEvent
	}

	return nil
}
