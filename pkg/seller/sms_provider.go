package seller

func NewSMSProvider() NotiProvider {
	return &emailProvider{}
}

type (
	smsProvider struct {
		// inject SMS provider here
	}
)

func (ep *smsProvider) StockChanged(oldStock int, newStock int, product string) {

}

func (ep *smsProvider) Type() ProviderType {
	return SMS
}
