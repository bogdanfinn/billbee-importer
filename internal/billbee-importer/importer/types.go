package importer

import (
	"github.com/justtrackio/gosoline/pkg/cfg"
	"github.com/justtrackio/gosoline/pkg/log"
)

type ImporterFactory func(name string, config cfg.Config, logger log.Logger) (Importer, error)

type Importer interface {
	GetName() string
	GetEmailEndpoint() string
	GetBillbeeSales() []BillbeeSale
}

var Sources = map[string]ImporterFactory{
	"alias":  NewAliasImporter,
	"stockx": NewStockxImporter,
}

type AliasSale struct {
	Rows []AliasSaleRow
}

type AliasSaleRow struct {
	Site                                string  `csv:"Site"`
	OrderID                             string  `csv:"OrderID"`
	Amount                              float64 `csv:"Amount"`
	ItemNameWithSize                    string  `csv:"ItemNameWithSize"`
	InvoiceDate                         string  `csv:"InvoiceDate"`
	DeliveryDate                        string  `csv:"DeliveryDate"`
	Currency                            string  `csv:"Currency"`
	TrackingNumber                      string  `csv:"TrackingNumber"`
	IsThisAnInvoiceWithCurrencyExchange string  `csv:"IsThisAnInvoiceWithCurrencyExchange"`
	ConvertedCurrencySign               string  `csv:"ConvertedCurrencySign"`
	ExchangeRatePorcentage              float64 `csv:"ExchangeRatePorcentage"`
	FinalAmountExchanged                float64 `csv:"FinalAmountExchanged"`
	StorePrice                          float64 `csv:"StorePrice"`
	ProcessingFee                       float64 `csv:"ProcessingFee"`
	ShippingFee                         float64 `csv:"ShippingFee"`
	TransactionFee                      float64 `csv:"TransactionFee"`
	BulkShipDiscount                    string  `csv:"BulkShipDiscount"`
	TaxExemptNotice                     string  `csv:"TaxExemptNotice"`
	TotalAmountNet                      float64 `csv:"TotalAmountNet"`
	VAT                                 int     `csv:"VAT"`
	TotalVATIncluded                    float64 `csv:"TotalVATIncluded"`
}

type StockxSale struct {
	SellerName                  string  `csv:"Seller Name"`
	SellerAddressBilling        string  `csv:"Seller Address (Billing)"`
	SellerVatNumber             string  `csv:"Seller VAT Number"`
	SoldToName                  string  `csv:"Sold To Name"`
	SoldToAddress               string  `csv:"Sold To address"`
	SoldToVatNumber             string  `csv:"Sold to VAT Number"`
	FiscalRepresentativeName    string  `csv:"Fiscal Representative Name"`
	FiscalRepresentativeAddress string  `csv:"Fiscal Representative Address"`
	StockxDestinationAddress    string  `csv:"StockX Destination Address"`
	SaleDate                    string  `csv:"Sale Date"`
	PayoutDate                  string  `csv:"Payout date"`
	InvoiceDate                 string  `csv:"Invoice Date"`
	InvoiceNumber               string  `csv:"Invoice Number"`
	OrderNumber                 string  `csv:"Order Number"`
	Item                        string  `csv:"Item"`
	SkuName                     string  `csv:"Sku Name"`
	SkuSize                     string  `csv:"Sku Size"`
	Style                       string  `csv:"Style"`
	Quantity                    float64 `csv:"Quantity"`
	ShipFromAddress             string  `csv:"Ship From Address"`
	ShipToAddress               string  `csv:"Ship To Address"`
	ListingPrice                float64 `csv:"Listing Price"`
	ListingPriceCurrency        string  `csv:"Listing Price Currency"`
	SellerFee                   float64 `csv:"Seller Fee"`
	SellerFeeCurrency           string  `csv:"Seller Fee Currency"`
	PaymentProcessing           float64 `csv:"Payment Processing"`
	PaymentProcessingCurrency   string  `csv:"Payment Processing Currency"`
	ShippingFee                 string  `csv:"Shipping Fee"`
	ShippingFeeCurrency         string  `csv:"Shipping Fee Currency"`
	VatRate                     string  `csv:"VAT Rate"`
	TotalNetAmount              string  `csv:"Total Net Amount (Payout Excluding VAT)"`
	TotalNetAmountCurrency      string  `csv:"Total Net Amount (Payout Excluding VAT) Currency"`
	TotalVatAmount              string  `csv:"Total VAT Amount"`
	TotalVatAmountCurrency      string  `csv:"Total VAT Amount Currency"`
	TotalGrossAmount            string  `csv:"Total Gross Amount (Total Payout)"`
	TotalGrossAmountCurrency    string  `csv:"Total Gross Amount (Total Payout) Currency"`
	SpecialReferences           string  `csv:"Special references"`
}

type BillbeeSale struct {
	OrderID       string             `json:"order_id"`
	CurrencyCode  string             `json:"currency_code"`
	PaymentMethod string             `json:"payment_method"`
	OrderDate     string             `json:"order_date"`
	UstID         string             `json:"ust_id"`
	Email         string             `json:"email"`
	Telephone     string             `json:"telephone"`
	Invoice       BillbeeInvoice     `json:"invoice"`
	Shipping      BillbeeShipping    `json:"shipping"`
	OrderStatusID int                `json:"order_status_id"`
	InternalNote  string             `json:"internal_note"`
	Comment       string             `json:"comment"`
	Shipcost      float64            `json:"shipcost"`
	ShippingDate  string             `json:"shipping_date"`
	OrderItems    []BillbeeOrderItem `json:"order_items"`
	Tags          []string           `json:"tags"`
}

type BillbeeShipping struct {
	Salutation  int    `json:"salutation"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Housenumber string `json:"housenumber"`
	Address1    string `json:"address_1"`
	Address2    string `json:"address_2"`
	Zip         string `json:"zip"`
	State       string `json:"state"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Company     string `json:"company"`
}

type BillbeeInvoice struct {
	Salutation  int    `json:"salutation"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Housenumber string `json:"housenumber"`
	Address1    string `json:"address_1"`
	Address2    string `json:"address_2"`
	Zip         string `json:"zip"`
	State       string `json:"state"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Company     string `json:"company"`
}

type BillbeeOrderItem struct {
	Quantity   float64     `json:"quantity"`
	Totalprice float64     `json:"totalprice"`
	ProductID  interface{} `json:"product_id"`
	Name       string      `json:"name"`
	Sku        string      `json:"sku"`
	Imageurl   string      `json:"imageurl"`
	Taxrate    int         `json:"taxrate"`
}
