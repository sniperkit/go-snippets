package fedex

import (
	"context"
	"encoding/xml"
	"time"

	"github.com/lensrentals/lensrentals-shipping/common"
	"github.com/lensrentals/lensrentals-shipping/services"
)

type (
	serviceAvailabilityRequest struct {
		XMLName                 xml.Name `xml:"http://fedex.com/ws/vacs/v8 ServiceAvailabilityRequest"`
		WebAuthenticationDetail webAuthenticationDetail
		ClientDetail            clientDetail
		Version                 version

		Origin      address
		Destination address
		ShipDate    string
		Packaging   string
	}

	serviceAvailabilityResponse struct {
		HighestSeverity string
		Notifications   []Notification

		ServiceAvailabilityOptions []serviceAvailabilityOption
	}

	serviceAvailabilityOption struct {
		Service      string
		DeliveryDate string // YYYY-MM-DD
		DeliveryDay  string
		TransitTime  string
	}
)

func newTimeInTransitRequest(c *Client, request services.TimeInTransitRequest) serviceAvailabilityRequest {
	var r serviceAvailabilityRequest

	r.WebAuthenticationDetail, r.ClientDetail = c.apiCredentials()
	r.Version = newRequestVersion("vacs", "8", "0", "0")

	r.ShipDate = request.ShipDate.Format("2006-01-02")
	r.Packaging = "YOUR_PACKAGING"
	r.Origin = address{
		PostalCode:  "38018",
		CountryCode: "US",
	}
	r.Destination = address{
		PostalCode:  request.PostalCode,
		CountryCode: "US",
	}

	return r
}

func (r serviceAvailabilityResponse) errorStatus() error {
	if r.HighestSeverity == "ERROR" || r.HighestSeverity == "FAILURE" {
		for _, notification := range r.Notifications {
			if err := notification.ToError(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (o serviceAvailabilityOption) timeInTransit() int {
	return timeInTransitFromString(o.TransitTime)
}

func (o serviceAvailabilityOption) deliveryDate() (time.Time, error) {
	return time.Parse("2006-01-02", o.DeliveryDate)
}

// TimeInTransit takes the given request and returns transit data in the
// response.
func (c *Client) TimeInTransit(ctx context.Context, request services.TimeInTransitRequest) (services.TimeInTransitResponse, error) {
	req := newTimeInTransitRequest(c, request)

	resp, err := c.serviceAvailabilityRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := resp.errorStatus(); err != nil {
		return nil, err
	}

	var response services.TimeInTransitResponse
	for _, option := range resp.ServiceAvailabilityOptions {

		serviceType := serviceFromXMLString(option.Service)
		var transitTime int
		var deliveryTime string

		if option.TransitTime != "" {
			transitTime = option.timeInTransit()
		} else if ts, err := option.deliveryDate(); !ts.IsZero() && err == nil {
			deliveryTime = ts.Format(time.RFC3339)
			transitTime = common.BusinessDaysBetween(request.ShipDate, ts)
		} else {
			transitTime = serviceType.expectedTimeInTransit()
		}

		service := services.TimeInTransitService{
			Service:       serviceType,
			TimeInTransit: transitTime,
			DeliveryTime:  deliveryTime,
		}
		response = append(response, service)
	}

	return response, nil
}

func (c *Client) serviceAvailabilityRequest(ctx context.Context, request serviceAvailabilityRequest) (*serviceAvailabilityResponse, error) {
	var response serviceAvailabilityResponse
	if err := c.apiRequest(ctx, &request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func timeInTransitFromString(s string) int {
	switch s {
	case "ONE_DAY":
		return 1
	case "TWO_DAYS":
		return 2
	case "THREE_DAYS":
		return 3
	case "FOUR_DAYS":
		return 4
	case "FIVE_DAYS":
		return 5
	case "SIX_DAYS":
		return 6
	case "SEVEN_DAYS":
		return 7
	case "EIGHT_DAYS":
		return 8
	case "NINE_DAYS":
		return 9
	case "TEN_DAYS":
		return 10
	case "ELEVEN_DAYS":
		return 11
	case "TWELVE_DAYS":
		return 12
	case "THIRTEEN_DAYS":
		return 13
	case "FOURTEEN_DAYS":
		return 14
	case "FIFTEEN_DAYS":
		return 15
	case "SIXTEEN_DAYS":
		return 16
	case "SEVENTEEN_DAYS":
		return 17
	case "EIGHTEEN_DAYS":
		return 18
	case "NINETEEN_DAYS":
		return 19
	case "TWENTY_DAYS":
		return 20
	case "UNKNOWN":
		return 0
	default:
		return -1
	}
}
