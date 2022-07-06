package importer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"github.com/justtrackio/gosoline/pkg/cfg"
	"github.com/justtrackio/gosoline/pkg/log"
)

type aliasImporter struct {
	name          string
	emailEndpoint string
	sales         []AliasSale
	payoutMethod  string
	currency      string
	orderStatusId int
	vatRate       int
	internalNote  string
	importComment string
	vatId         string
}

func NewAliasImporter(name string, config cfg.Config, logger log.Logger) (Importer, error) {
	ex, _ := os.Executable()
	aliasCsvFilePath := filepath.Join(filepath.Dir(ex), "alias.csv")

	aliasSales, err := ParseCsvToAliasStruct(aliasCsvFilePath)

	if err != nil {
		println("Error occurred while parsing alias.csv file. Please check it")

		return nil, err
	}

	aliasBillbeeEmailEndpoint := config.GetString("alias_billbee_email_endpoint", "FILL_HERE")

	if aliasBillbeeEmailEndpoint == "SHOP_ID-json-orderimport@inbound.billbee.de" {
		println("Please provide your Billbee Email Endpoint for alias in the config.dist.yml file")

		return nil, fmt.Errorf("alias billbee email endpoint is missing")
	}

	aliasPayoutMethod := config.GetString("alias_payout_method", "paypal")
	aliasCurrency := config.GetString("alias_currency", "EUR")
	aliasOrderStatusId := config.GetInt("alias_order_status_id", 7)
	aliasImportComment := config.GetString("alias_optional_comment", "0")
	aliasInternalNote := config.GetString("alias_internal_note", "Generated from Importer Application")
	aliasVatRate := config.GetInt("alias_vat_rate", 19)
	aliasVatId := config.GetString("alias_vat_id", "NL 826259716B01")

	return &aliasImporter{
		name:          name,
		emailEndpoint: aliasBillbeeEmailEndpoint,
		sales:         aliasSales,
		payoutMethod:  aliasPayoutMethod,
		currency:      aliasCurrency,
		orderStatusId: aliasOrderStatusId,
		vatRate:       aliasVatRate,
		internalNote:  aliasInternalNote,
		importComment: aliasImportComment,
		vatId:         aliasVatId,
	}, nil
}

func (a aliasImporter) GetName() string {
	return a.name
}

func (a aliasImporter) GetEmailEndpoint() string {
	return a.emailEndpoint
}

func (a aliasImporter) GetBillbeeSales() []BillbeeSale {
	var billbeeSales []BillbeeSale

	for _, aliasSale := range a.sales {
		firsSale := aliasSale.Rows[0]
		var orderItems []BillbeeOrderItem

		transactionFee := 0.0
		processingFee := 0.0
		for _, aliasSaleItem := range aliasSale.Rows {
			orderItems = append(orderItems, BillbeeOrderItem{
				Quantity:   aliasSaleItem.Amount,
				Totalprice: aliasSaleItem.StorePrice * aliasSaleItem.ExchangeRatePorcentage,
				ProductID:  nil,
				Name:       aliasSaleItem.ItemNameWithSize,
				Sku:        "",
				Imageurl:   "",
				Taxrate:    aliasSaleItem.VAT,
			})

			transactionFee += aliasSaleItem.TransactionFee * aliasSaleItem.ExchangeRatePorcentage
			processingFee += aliasSaleItem.ProcessingFee * aliasSaleItem.ExchangeRatePorcentage
		}

		orderItems = append(orderItems, BillbeeOrderItem{
			Quantity:   1,
			Totalprice: processingFee,
			ProductID:  nil,
			Name:       "./. Verkaufsgebühr",
			Sku:        "",
			Imageurl:   "",
			Taxrate:    firsSale.VAT,
		})

		orderItems = append(orderItems, BillbeeOrderItem{
			Quantity:   1,
			Totalprice: transactionFee,
			ProductID:  nil,
			Name:       "./. Zahlungsabwicklungsgebühr",
			Sku:        "",
			Imageurl:   "",
			Taxrate:    firsSale.VAT,
		},
		)

		billbeeSale := BillbeeSale{
			OrderID:       firsSale.OrderID,
			CurrencyCode:  a.currency,
			PaymentMethod: a.payoutMethod,
			OrderDate:     convertAliasDate(firsSale.InvoiceDate),
			UstID:         a.vatId,
			Email:         "",
			Telephone:     "",
			Invoice: BillbeeInvoice{
				Salutation:  0,
				Firstname:   "",
				Lastname:    "",
				Housenumber: "3433",
				Address1:    "W Exposition Place",
				Address2:    "",
				Zip:         "90018",
				State:       "CA",
				City:        "Los Angeles",
				Country:     "UNITED STATES OF AMERICA",
				Company:     "GOAT",
			},
			Shipping: BillbeeShipping{
				Salutation:  0,
				Firstname:   "",
				Lastname:    "",
				Housenumber: "25",
				Address1:    "Columbusstraat",
				Address2:    "",
				Zip:         "3165AC",
				State:       "",
				City:        "Rotterdam-Albrandswaard",
				Country:     "NETHERLANDS",
				Company:     "1661 Inc",
			},
			OrderStatusID: a.orderStatusId,
			InternalNote:  a.internalNote,
			Comment:       a.importComment,
			Shipcost:      0,
			ShippingDate:  convertAliasDate(firsSale.DeliveryDate),
			OrderItems:    orderItems,
			Tags:          nil,
		}

		billbeeSales = append(billbeeSales, billbeeSale)
	}
	return billbeeSales
}

func ParseCsvToAliasStruct(filePath string) ([]AliasSale, error) {
	salesFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer salesFile.Close()

	var aliasSaleRows []AliasSaleRow
	err = gocsv.UnmarshalFile(salesFile, &aliasSaleRows)

	if err != nil {
		return nil, err
	}

	aliasSales := combineRowsToAliasSales(aliasSaleRows)

	return aliasSales, nil
}

func combineRowsToAliasSales(rows []AliasSaleRow) []AliasSale {
	var sales []AliasSale

	sale := AliasSale{Rows: []AliasSaleRow{}}
	lastExchangeRate := 0.0
	for i, row := range rows {
		if row.ExchangeRatePorcentage != lastExchangeRate {
			if len(sale.Rows) > 0 {
				sales = append(sales, sale)
			}

			sale = AliasSale{Rows: []AliasSaleRow{}}
		}

		lastExchangeRate = row.ExchangeRatePorcentage
		sale.Rows = append(sale.Rows, row)

		if i == len(rows)-1 {
			sales = append(sales, sale)
		}
	}

	return sales
}
