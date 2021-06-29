package seller

type ProviderType int

//go:generate enumer -type=ProviderType -transform=kebab -text -yaml -output=z_provider_type_enumer.go
const (
	Email ProviderType = iota
	SMS
)

type (
	NotiProvider interface {
		StockChanged(oldStock int, newStock int, product string, sl *Seller)
		Type() ProviderType
	}
)
