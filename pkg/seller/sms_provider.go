package seller

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func NewSMSProvider() NotiProvider {
	return &smsProvider{}
}

type (
	smsProvider struct {
		// inject SMS provider here
	}
)

func (ep *smsProvider) StockChanged(oldStock int, newStock int, product string, sl *Seller) {
	log.Info().Msg(fmt.Sprintf("%s Warning sent to %s (Phone: %s): %s Product stock changed", "SMS", sl.UUID, sl.Phone, product))
}

func (ep *smsProvider) Type() ProviderType {
	return SMS
}
