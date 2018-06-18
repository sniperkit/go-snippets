package mitigation

import (
	"errors"
	"fmt"

	"github.com/go-zero-boilerplate/facades/alerting/sentry"
	"github.com/go-zero-boilerplate/facades/logging"
)

//TryAndAlert will (if the mitigateCondition is TRUE) call the mitigateFunc and log the answer getMitigationMessage()
func TryAndAlert(mitigateCondition bool, mitigateFunc func() error, getMitigationMessage func() string) error {
	if !mitigateCondition {
		return nil
	}

	logger := logging.Get()
	mitigationMessage := getMitigationMessage()
	sentry.GetClient().CaptureError(errors.New(mitigationMessage), nil)
	logger.WithField("is-before-mitigate", true).Error("Attempting to mitigate: " + mitigationMessage)

	if err := mitigateFunc(); err != nil {
		userMessage := "Unexpected server error, attempt and failed to mitigate"
		logger.WithError(err).WithField("failed-mitigate-attempt", true).Error(userMessage)
		sentry.CaptureErrorPacketAndWait(fmt.Sprintf("Attempted and FAILED to mitigate '%s'", mitigationMessage), err, nil, nil)
		return errors.New(userMessage)
	}

	return nil
}
