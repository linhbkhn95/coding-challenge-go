package seller

func NewEmailProvider() NotiProvider {
	return &emailProvider{}
}

type (
	emailProvider struct {
		// inject email provider here
	}
)

func (ep *emailProvider) StockChanged(oldStock int, newStock int, product string, sl *Seller) {

}

func (ep *emailProvider) Type() ProviderType {
	return Email
}
