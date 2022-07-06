package importer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"github.com/justtrackio/gosoline/pkg/cfg"
	"github.com/justtrackio/gosoline/pkg/log"
)

type stockxImporter struct {
	name          string
	emailEndpoint string
	sales         []StockxSale
	payoutMethod  string
	currency      string
	orderStatusId int
	vatRate       int
	internalNote  string
	importComment string
}

func NewStockxImporter(name string, config cfg.Config, logger log.Logger) (Importer, error) {
	ex, _ := os.Executable()
	stockxCsvFilePath := filepath.Join(filepath.Dir(ex), "stockx.csv")

	stockxSales, err := ParseCsvToStockxStruct(stockxCsvFilePath)

	if err != nil {
		println("Error occurred while parsing stockx.csv file. Please check it")

		return nil, err
	}

	stockxBillbeeEmailEndpoint := config.GetString("stockx_billbee_email_endpoint", "FILL_HERE")

	if stockxBillbeeEmailEndpoint == "SHOP_ID-json-orderimport@inbound.billbee.de" {
		println("Please provide your stockx Billbee Email Endpoint in the config.dist.yml file")

		return nil, fmt.Errorf("billbee email endpoint is missing")
	}

	stockxPayoutMethod := config.GetString("stockx_payout_method", "paypal")
	stockxCurrency := config.GetString("stockx_currency", "EUR")
	stockxOrderStatusId := config.GetInt("stockx_order_status_id", 7)
	stockxImportComment := config.GetString("stockx_optional_comment", "0")
	stockxInternalNote := config.GetString("stockx_internal_note", "Generated from Importer Application")
	stockxVatRate := config.GetInt("stockx_vat_rate", 19)

	return &stockxImporter{
		name:          name,
		emailEndpoint: stockxBillbeeEmailEndpoint,
		sales:         stockxSales,
		payoutMethod:  stockxPayoutMethod,
		currency:      stockxCurrency,
		orderStatusId: stockxOrderStatusId,
		vatRate:       stockxVatRate,
		internalNote:  stockxInternalNote,
		importComment: stockxImportComment,
	}, nil
}

func (s stockxImporter) GetName() string {
	return s.name
}

func (s stockxImporter) GetEmailEndpoint() string {
	return s.emailEndpoint
}

func (s stockxImporter) GetBillbeeSales() []BillbeeSale {
	var billbeeSales []BillbeeSale

	for _, stockxSale := range s.sales {
		orderItems := []BillbeeOrderItem{
			{
				Quantity:   stockxSale.Quantity,
				Totalprice: stockxSale.ListingPrice,
				ProductID:  nil,
				Name:       stockxSale.SkuName,
				Sku:        "",
				Imageurl:   "",
				Taxrate:    s.vatRate,
			},
			{
				Quantity:   1,
				Totalprice: negate(stockxSale.SellerFee),
				ProductID:  nil,
				Name:       "./. Transaktionsgebühr",
				Sku:        "",
				Imageurl:   "",
				Taxrate:    s.vatRate,
			},
			{
				Quantity:   1,
				Totalprice: negate(stockxSale.PaymentProcessing),
				ProductID:  nil,
				Name:       "./. Zahlungsabwicklungsgebühr",
				Sku:        "",
				Imageurl:   "",
				Taxrate:    s.vatRate,
			},
		}

		billbeeSale := BillbeeSale{
			OrderID:       stockxSale.OrderNumber,
			CurrencyCode:  s.currency,
			PaymentMethod: s.payoutMethod,
			OrderDate:     stockxSale.SaleDate,
			UstID:         stockxSale.SoldToVatNumber,
			Email:         "",
			Telephone:     "",
			Invoice: BillbeeInvoice{
				Salutation:  0,
				Firstname:   "StockX LLC",
				Lastname:    "",
				Housenumber: "1046",
				Address1:    "Woodward Avenue",
				Address2:    "",
				Zip:         "48226",
				State:       "MI",
				City:        "Detroit",
				Country:     "UNITED STATES OF AMERICA",
				Company:     "StockX LLC",
			},
			Shipping: BillbeeShipping{
				Salutation:  0,
				Firstname:   "StockX LLC",
				Lastname:    "",
				Housenumber: "",
				Address1:    "DE RUN 4312",
				Address2:    "SMOT",
				Zip:         "5503LN",
				State:       "",
				City:        "VELDHOVEN",
				Country:     "NETHERLANDS",
				Company:     "",
			},
			OrderStatusID: s.orderStatusId,
			InternalNote:  s.internalNote,
			Comment:       s.importComment,
			Shipcost:      0,
			ShippingDate:  stockxSale.PayoutDate,
			OrderItems:    orderItems,
			Tags:          nil,
		}

		billbeeSales = append(billbeeSales, billbeeSale)
	}

	return billbeeSales
}

func ParseCsvToStockxStruct(filePath string) ([]StockxSale, error) {
	salesFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer salesFile.Close()

	var stockxSales []StockxSale
	err = gocsv.UnmarshalFile(salesFile, &stockxSales)

	if err != nil {
		return nil, err
	}

	return stockxSales, nil
}
