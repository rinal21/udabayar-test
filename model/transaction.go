package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	Code                    string    `form:"code" json:"code"`
	Status                  int       `form:"status" json:"status"`
	Shop                    Shop      `form:"shop" json:"shop" binding:"-"`
	ShopID                  int       `form:"shop_id" json:"shop_id" binding:"required"`
	Amount                  int       `form:"amount" json:"amount"`
	UserBank                UserBank  `form:"userbank" json:"userbank" binding:"-"`
	UserBankID              int       `form:"userbank_id" json:"userbank_id"`
	PaymentExpiredAt        time.Time `form:"payment_expired_at" json:"payment_expired_at"`
	Courier                 string    `form:"courier" json:"courier"`
	CourierType             string    `form:"courier_type" json:"courier_type"`
	PaymentCode             int       `form:"payment_code" json:"payment_code"`
	PriceShipment           int       `form:"price_shipment" json:"price_shipment"`
	EstimatedShipping       string    `form:"estimated_shipping" json:"estimated_shipping"`
	TrackingNumber          string    `form:"tracking_number" json:"tracking_number"`
	Note                    string    `form:"note" json:"note"`
	CustomerName            string    `form:"customer_name" json:"customer_name" binding:"required"`
	CustomerEmail           string    `form:"customer_email" json:"customer_email" binding:"required"`
	CustomerPhoneNumber     string    `form:"customer_phone" json:"customer_phone" binding:"required"`
	CustomerAddress         string    `form:"customer_address" json:"customer_address"`
	CustomerProvince        string    `form:"customer_province" json:"customer_province"`
	CustomerCityOrDistrict  string    `form:"customer_city_or_district" json:"customer_city_or_district"`
	CustomerSubdistrict     string    `form:"customer_subdistrict" json:"customer_subdistrict"`
	CustomerDestinationCode int       `form:"customer_destination_code" json:"customer_destination_code"`
	PostalCode              int       `form:"postal_code" json:"postal_code"`
	Order                   Order     `json:"order" binding:"required"`
}

type Order struct {
	gorm.Model
	TransactionID   int           `form:"transaction_id" json:"transaction_id"`
	Product         Product       `form:"product" json:"product" binding:"-"`
	ProductID       int           `form:"product_id" json:"product_id"`
	ProductDetail   ProductDetail `form:"product_detail" json:"product_detail" binding:"-"`
	ProductDetailID int           `form:"product_detail_id" json:"product_detail_id" binding:"required"`
	Qty             int           `form:"qty" json:"qty" binding:"required"`
	Price           int           `form:"price" json:"price"`
}
