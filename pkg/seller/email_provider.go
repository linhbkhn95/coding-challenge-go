package seller

func NewEmailProvider() EmailProvider {
	return &emailProvider{}
}

type (
	EmailProvider interface {
		StockChanged(oldStock int, newStock int, product string)
	}
	emailProvider struct {
	}
)

func (ep *emailProvider) StockChanged(oldStock int, newStock int, product string) {

}
