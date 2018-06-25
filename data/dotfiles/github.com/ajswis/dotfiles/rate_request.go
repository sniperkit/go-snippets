package fedex

import (
	"context"
	"encoding/xml"
	"time"

	"github.com/lensrentals/lensrentals-shipping/services"
)

type (
	// A rateRequest outlines the structure for a Rate Request that FedEx will
	// reply to.
	rateRequest struct {
		XMLName                 xml.Name `xml:"http://fedex.com/ws/rate/v20 RateRequest"`
		WebAuthenticationDetail webAuthenticationDetail
		ClientDetail            clientDetail
		Version                 version

		ReturnTransitAndCommit bool
		RequestedShipment      requestedShipment
	}

	// A rateReply outlines the structure for a Rate Reply that FedEx will
	// respond with.
	rateReply struct {
		HighestSeverity   string
		Notifications     []Notification
		TransactionDetail transactionDetail
		Version           version
		RateReplyDetails  []rateReplyDetail
	}

	rateReplyDetail struct {
		ServiceType       string
		PackagingType     string
		DeliveryDayOfWeek string
		DeliveryTimestamp string
		TransitTime       string

		CommitDetails        commitDetails
		RatedShipmentDetails []ratedShipmentDetail
	}

	commitDetails []commitDetail

	ratedShipmentDetail struct {
		ShipmentRateDetail shipmentRateDetail
	}

	shipmentRateDetail struct {
		RateType       string
		TotalNetCharge money
	}

	commitDetail struct {
		CommitTimestamp string
	}
)

func newRateRequest(c *Client, request services.ShipmentRequest) (*rateRequest, error) {
	var r rateRequest
	r.WebAuthenticationDetail, r.ClientDetail = c.apiCredentials()
	r.Version = newRequestVersion("crs", "20", "0", "0")
	r.ReturnTransitAndCommit = true
	r.RequestedShipment = requestedShipment{RateRequestTypes: "LIST"}

	returning := request.Service.Return
	if err := r.RequestedShipment.setService(request.Service.Name, returning); err != nil {
		return nil, err
	}
	if err := r.RequestedShipment.setShipper(request.Shipper); err != nil {
		return nil, err
	}
	if err := r.RequestedShipment.setOrigin(request.ShipFrom, returning); err != nil {
		return nil, err
	}
	if err := r.RequestedShipment.setRecipient(request.ShipTo); err != nil {
		return nil, err
	}
	if err := r.RequestedShipment.setPackages(request.Packages, request.PONumber); err != nil {
		return nil, err
	}

	r.RequestedShipment.setShipDate(request.ShipDate.Time)

	return &r, nil
}

// errorStatus returns any errors from an otherwise successful rateReply.
func (r *rateReply) errorStatus() error {
	if r.HighestSeverity == "ERROR" || r.HighestSeverity == "FAILURE" {
		for _, notification := range r.Notifications {
			if err := notification.ToError(); err != nil {
				return err
			}
		}
	}

	for _, notification := range r.Notifications {
		if notification.Code == "556" {
			return ErrorNotification{notification}
		}
	}

	return nil
}

func (r rateReplyDetail) price() float32 {
	var amount float32
	for _, detail := range r.RatedShipmentDetails {
		if detail.ShipmentRateDetail.RateType == "PAYOR_ACCOUNT_PACKAGE" {
			// PAYOR_ACCOUNT_PACKAGE contains total amounts, so default to that
			return detail.ShipmentRateDetail.TotalNetCharge.Amount
		} else if detail.ShipmentRateDetail.TotalNetCharge.Amount != 0 {
			// Pick one of the other ShipmentRateDetails and use that for the amount
			amount = detail.ShipmentRateDetail.TotalNetCharge.Amount
		}
	}

	return amount
}

// timeInTransit converts FedEx integer strings to integers.. because
// FedEx doesn't know what an integer is.. or numbers past 20.
func (r rateReplyDetail) timeInTransit() int {
	return timeInTransitFromString(r.TransitTime)
}

func (r rateReplyDetail) deliveryBy() time.Time {
	t, _ := r.CommitDetails.timestamp()
	return t
}

func (d commitDetails) timestamp() (time.Time, error) {
	var ts string
	for _, detail := range d {
		if detail.CommitTimestamp != "" {
			ts = detail.CommitTimestamp
			break
		}
	}

	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05", ts)
	}
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

// RatedService contains pertinent rate information for a service.
type RatedService struct {
	Service
	price         float32
	timeInTransit int
	deliveryBy    time.Time
}

var _ services.RatedService = (*RatedService)(nil)

// Price returns the rated price for the service.
func (rs RatedService) Price() float32 {
	return rs.price
}

func (rs RatedService) TimeInTransit() int {
	return rs.timeInTransit
}

func (rs RatedService) DeliveryBy() time.Time {
	return rs.deliveryBy
}

// ListOptions returns a slice of potential Services for the given request,
// without accounting for transit time, price, etc; only the selected
// Carrier/Service Name is inspected to create the slice.
func (c *Client) ListOptions(request services.ShipmentRequest) []services.Service {
	possibleServices := []Service{
		Ground,
		ExpressSaver,
		TwoDay,
		StandardOvernight,
		PriorityOvernight,
	}
	var shippableServices []services.Service

	// Pre-filter available services
	for _, shipService := range possibleServices {
		carrier := request.Service.Carrier
		service := request.Service.Name

		shippable := shipService.isShippable(request.Service.Return)
		validCarrier := carrier == "" || carrier == c.Name()
		validService := service == "" || service == shipService.ServiceName()

		// After checking constraints, if the service is still valid, append it
		if shippable && validCarrier && validService {
			shippableServices = append(shippableServices, shipService)
		}
	}

	return shippableServices
}

func (c *Client) Rate(ctx context.Context, request services.ShipmentRequest) (services.RatedServices, services.ShipError) {
	// Don't request a specific service for these ratings
	request.Service.Name = ""

	response, err := c.rate(ctx, request)
	if err != nil {
		return nil, newShipError(err)
	}

	proposal := services.RatedServices{}
	for _, detail := range response.RateReplyDetails {
		proposal = append(proposal, RatedService{
			Service:       serviceFromXMLString(detail.ServiceType),
			price:         detail.price(),
			deliveryBy:    detail.deliveryBy(),
			timeInTransit: detail.timeInTransit(),
		})
	}

	return proposal, nil
}

// rate returns raw rate data for the shipment request
func (c *Client) rate(ctx context.Context, request services.ShipmentRequest) (*rateReply, error) {
	// Prepare the request
	r, err := newRateRequest(c, request)
	if err != nil {
		return nil, err
	}

	resp, err := c.rateRequest(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// rateRequest() performs the API action by calling the Client's apiRequest()
// method. If the apiRequest() fails or the reply contained a Notification
// indicating an error, rateRequest() returns an error.  Otherwise, the FedEx
// rateReply is returned.
func (c *Client) rateRequest(ctx context.Context, request *rateRequest) (*rateReply, error) {
	var response rateReply

	if err := c.apiRequest(ctx, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
