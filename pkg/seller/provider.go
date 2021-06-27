package seller

type (
	NotiProvider interface {
		StockChanged(oldStock int, newStock int, product string)
		Type() ProviderType
	}
)
