package model

type ProductType string

const (
	ProductElectronics ProductType = "электроника"
	ProductClothes     ProductType = "одежда"
	ProductShoes       ProductType = "обувь"
)

type Product struct {
	ID          string
	DateTime    string
	ReceptionID string
	Type        ProductType
}
