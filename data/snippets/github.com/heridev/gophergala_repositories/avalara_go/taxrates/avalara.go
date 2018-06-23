package taxrates

import (
	"fmt"
	"net/http"

	"github.com/jmcvetta/napping"
)

const AvalaraEndpoint = "https://taxrates.api.avalara.com:443/postal"

type Avalara struct {
	Endpoint string
	APIKey   string
	Session  napping.Session
}

func NewAvalaraClient(APIKey string) *Avalara {
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("AvalaraApiKey %s", APIKey))

	session := napping.Session{
		Header: &header,
	}

	return &Avalara{
		Endpoint: AvalaraEndpoint,
		APIKey:   APIKey,
		Session:  session,
	}
}

func (a *Avalara) GetRates(country, postal string) (TaxResponse, error) {
	response := TaxResponse{}

	params := napping.Params{
		"country": country,
		"postal":  postal,
	}

	_, err := a.Session.Get(a.Endpoint, &params, &response, nil)

	return response, err
}
