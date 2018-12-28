package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"udabayar-go-api-di/config"
	"udabayar-go-api-di/helper"
	"udabayar-go-api-di/model"

	"github.com/biezhi/gorm-paginator/pagination"
	humanize "github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	"github.com/thanhpk/randstr"
)

type TransactionHandler struct {
	DB                  *gorm.DB
	Config              *config.Config
	TxConfig            *config.TransactionConfig
	Email               *helper.Email
	ShippingHandler     *ShippingHandler
	SocketServerHandler *SocketServerHandler
	TimeHelper          *helper.TimeHelper
	UnitHelper          *helper.UnitHelper
}

func (h *TransactionHandler) AllTransaction(c *gin.Context) {
	var transactions []model.Transaction
	db := h.DB

	UserID, _ := c.Get("UserID")
	ShopID := c.Param("shop_id")

	transactionStatus := c.DefaultQuery("status", "")

	if transactionStatus != "" {
		db = db.Where("status = ?", transactionStatus)
	}

	db = db.Preload("Shop").Preload("Order.Product.Shop").Where("shop_id = ?", ShopID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	order := c.DefaultQuery("order", "asc")

	pagination.Pagging(&pagination.Param{
		DB:      db,
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id " + order},
	}, &transactions)

	if transactions[0].Shop.OwnerID != UserID {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Toko bukan milik kamu.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   transactions,
	})
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	var transaction model.Transaction
	var shop model.Shop

	transaction_code := c.Param("transaction_code")
	transaction_id := c.Param("id")

	if transaction_id != "" {
		if check := h.DB.Preload("Shop").Preload("Order.ProductDetail.ProductVariant").Preload("Order.Product.ProductReview").Preload("Order.Product.ProductImage").Preload("Order.Product.Shop").First(&transaction, transaction_id).RecordNotFound(); check == true {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "Transaksi tidak ditemukan.",
			})

			return
		}
	} else {
		if check := h.DB.Preload("UserBank.Bank.PaymentMethod").Preload("Shop").Preload("Order.ProductDetail.ProductVariant").Preload("Order.Product.ProductReview").Preload("Order.Product.ProductImage").Preload("Order.Product.Shop").First(&transaction, "code = ?", transaction_code).RecordNotFound(); check == true {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "Transaksi tidak ditemukan.",
			})

			return
		}
	}

	UserID, _ := c.Get("UserID")

	if UserID != nil {
		if err := h.DB.Where("owner_id = ?", UserID).First(&shop, transaction.ShopID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "Bukan milik toko kamu.",
			})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   transaction,
	})
}

type Order struct {
	ProductID int `form:"product_id" json:"product_id"`
}

type SendTransaction struct {
	UserBankID    int `json:"user_bank_id"`
	TransactionID int `json:"transaction_id"`
	Amount        int `json:"amount"`
}

type GenerateRandomJSON struct {
	Status  int                `json:"status"`
	Data    GenerateRandomData `json:"data"`
	Message string             `json:"message"`
}

type GenerateRandomData struct {
	Amount       int `json:"amount"`
	RandomNumber int `json:"random_number"`
	RawAmount    int `json:"raw_amount"`
}

// generate random
func (h *TransactionHandler) getGenerateRandom(userbankID int, totalAmount int) (*GenerateRandomJSON, error) {
	var response GenerateRandomJSON

	generateInput := &SendTransaction{
		UserBankID: userbankID,
		Amount:     totalAmount,
	}

	jsonGenerateInput, _ := json.Marshal(generateInput)
	stringJSONGenerateInput := []byte(string(jsonGenerateInput))

	reqGenerateRandom, _ := http.NewRequest("POST", h.Config.BankAPIUrl+"/randomNumber", bytes.NewBuffer(stringJSONGenerateInput))
	reqGenerateRandom.Header.Set("Content-Type", "application/json")

	clientGenerateRandom := &http.Client{}
	respGenerateRandom, err := clientGenerateRandom.Do(reqGenerateRandom)

	if err != nil {
		return &response, errors.New("Gagal menghubungkan ke server 2")
	}

	resGenerateRandom, _ := ioutil.ReadAll(respGenerateRandom.Body)

	json.Unmarshal(resGenerateRandom, &response)

	if response.Status == 0 {
		return &response, errors.New(response.Message)
	}

	return &response, nil
}

// send transaction to bank api
func (h *TransactionHandler) sendTransaction(userbank_id int, transaction_id int, totalAmount int) error {
	jsonW := &SendTransaction{
		UserBankID:    userbank_id,
		TransactionID: transaction_id,
		Amount:        totalAmount,
	}

	jsonR, _ := json.Marshal(jsonW)
	jsonStr := []byte(string(jsonR))

	req, _ := http.NewRequest("POST", h.Config.BankAPIUrl+"/transaction", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err2 := client.Do(req)

	if err2 != nil {
		return errors.New("Gagal menghubungkan ke server 2")
	}

	return nil
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction model.Transaction
	var productDetail model.ProductDetail
	var userbank model.UserBank

	if err := c.ShouldBind(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	// get product detail
	if check := h.DB.Preload("Product").Preload("ProductVariant").Find(&productDetail, transaction.Order.ProductDetailID).RecordNotFound(); check == true {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": "Produk tidak tersedia.",
		})

		return
	}

	// check avaible qty
	if productDetail.Qty < transaction.Order.Qty {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": "Kuantitas produk tidak mencukupi.",
		})

		return
	}

	// check userbank
	if check := h.DB.Preload("Bank").First(&userbank, transaction.UserBankID).RecordNotFound(); check == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Akun bank tidak ditemukan.",
		})

		return
	}

	// calculate price
	totalAmount := productDetail.Price * transaction.Order.Qty
	totalAmount = totalAmount + transaction.PriceShipment

	// generate payment expired time
	transaction.PaymentExpiredAt = time.Now().Local().Add(time.Hour * 48)

	// get generate random
	generateRandom, err := h.getGenerateRandom(transaction.UserBankID, totalAmount)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	transaction.Amount = generateRandom.Data.Amount
	transaction.PaymentCode = generateRandom.Data.RandomNumber

	// transaction.Amount = totalAmount
	// transaction.PaymentCode = 0

	// generate trasaction code
	token := randstr.String(16)
	transaction.Code = token

	// create transaction
	if err := h.DB.Omit("deleted_at").Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal membuat transaksi, silakan coba kembali.",
		})

		return
	}

	// send email for customer
	productTotalPrice := productDetail.Price * transaction.Order.Qty
	courierWeight := transaction.Order.Qty * productDetail.Product.Weight
	totalPrice := transaction.Amount

	templateData := helper.CreateTransactionTemplateData{
		Fullname:         transaction.CustomerName,
		TransactionCode:  transaction.Code,
		BankName:         userbank.Bank.Name,
		PaymentExpiredAt: h.TimeHelper.TimeFormat(transaction.PaymentExpiredAt),
		TotalAmount:      humanize.FormatFloat("#.###,", float64(transaction.Amount)),
		BankLogo:         userbank.Bank.Logo,
		AccountNumber:    userbank.AccountNumber,
		CheckStatusLink:  fmt.Sprintf("%stransaction/%s", h.Config.WebUrl, transaction.Code),
		DetailOrder: helper.DetailOrder{
			ProductName:               productDetail.Product.Name,
			ProductVariant:            productDetail.ProductVariant.Name,
			ProductPrice:              humanize.FormatFloat("#.###,", float64(productDetail.Price)),
			ProductQty:                transaction.Order.Qty,
			ProductWeight:             h.UnitHelper.WeightFormat(productDetail.Product.Weight),
			ProductTotalPrice:         humanize.FormatFloat("#.###,", float64(productTotalPrice)),
			Courier:                   transaction.Courier,
			CourierType:               transaction.CourierType,
			CourierWeight:             h.UnitHelper.WeightFormat(courierWeight),
			PriceShipment:             humanize.FormatFloat("#.###,", float64(transaction.PriceShipment)),
			TotalPrice:                humanize.FormatFloat("#.###,", float64(totalPrice)),
			TransactionGenerateRandom: transaction.PaymentCode,
			TotalAmount:               humanize.FormatFloat("#.###,", float64(transaction.Amount)),
		},
	}

	log.Println(productDetail.Product)

	go h.Email.SendEmailCreateTransaction(&templateData, transaction.CustomerEmail, productDetail.Product.NonPyshical)

	// payment expired time scheduler
	cron := gocron.NewScheduler()
	cron.Every(48).Hours().Do(func(cronjob *gocron.Scheduler) {
		if transaction.Status == 0 {
			go h.DB.Model(&transaction).Update("status", 7)
			go h.DB.Model(&productDetail).UpdateColumn("qty", gorm.Expr("qty + ?", transaction.Order.Qty))

			templateData := helper.TransactionCancelledTemplateData{
				Fullname:        transaction.CustomerName,
				TransactionCode: transaction.Code,
				TotalAmount:     humanize.FormatFloat("#.###,", float64(transaction.Amount)),
				BankName:        userbank.Bank.Name,
				DetailOrder: helper.DetailOrder{
					ProductName:               productDetail.Product.Name,
					ProductVariant:            productDetail.ProductVariant.Name,
					ProductPrice:              humanize.FormatFloat("#.###,", float64(productDetail.Price)),
					ProductQty:                transaction.Order.Qty,
					ProductWeight:             h.UnitHelper.WeightFormat(productDetail.Product.Weight),
					ProductTotalPrice:         humanize.FormatFloat("#.###,", float64(productTotalPrice)),
					Courier:                   transaction.Courier,
					CourierType:               transaction.CourierType,
					CourierWeight:             h.UnitHelper.WeightFormat(courierWeight),
					PriceShipment:             humanize.FormatFloat("#.###,", float64(transaction.PriceShipment)),
					TotalPrice:                humanize.FormatFloat("#.###,", float64(totalPrice)),
					TransactionGenerateRandom: transaction.PaymentCode,
					TotalAmount:               humanize.FormatFloat("#.###,", float64(transaction.Amount)),
				},
			}

			h.Email.SendEmailTransactionCancelled(&templateData, transaction.CustomerEmail, transaction.Order.Product.NonPyshical)
		}

		cronjob.Clear()
	}, cron)

	go cron.Start()

	// send transaction to bank api
	if err := h.sendTransaction(transaction.UserBankID, int(transaction.ID), transaction.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	// update product qty
	go h.DB.Omit("deleted_at").Model(&productDetail).UpdateColumn("qty", gorm.Expr("qty - ?", transaction.Order.Qty))

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Transaction created.",
		"data": gin.H{
			"transaction_code": transaction.Code,
			"amount":           transaction.Amount,
			"banks": gin.H{
				"bank_id":       userbank.Bank.ID,
				"bank_name":     userbank.Bank.Name,
				"bank_code":     userbank.Bank.Code,
				"accountNumber": userbank.AccountNumber,
				"accountName":   userbank.AccountName,
			},
		},
	})
}

type WebhookTransactionInput struct {
	TransactionID int `form:"transaction_id" json:"transaction_id" binding:"required"`
	Status        int `form:"status" json:"status" binding:"required"`
}

func (h *TransactionHandler) WebhookTransaction(c *gin.Context) {
	var transaction model.Transaction
	var input WebhookTransactionInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": fmt.Sprint("Validation error: ", err.Error()),
		})

		return
	}

	if err := h.DB.Preload("UserBank.Bank").Preload("Shop").Preload("Order.Product").Preload("Order.ProductDetail.ProductVariant").First(&transaction, input.TransactionID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if transaction.Status >= input.Status {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Status transaksi sudah diubah sebelumnya.",
		})

		return
	}

	if err := h.DB.Model(&transaction).Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal mengubah status transaksi.",
		})

		return
	}

	// socket
	event := fmt.Sprintf("transaction_%d", transaction.ID)

	msg := gin.H{
		"transaction_id": transaction.ID,
		"status":         transaction.Status,
	}

	h.SocketServerHandler.Server.BroadcastTo("transaction", event, msg)

	// if status paid
	if input.Status == h.TxConfig.Status.Paid {
		productTotalPrice := transaction.Order.ProductDetail.Price * transaction.Order.Qty
		courierWeight := transaction.Order.Qty * transaction.Order.Product.Weight
		totalPrice := productTotalPrice + transaction.PriceShipment

		templateData := helper.PaymentSuccessTemplateData{
			Fullname:        transaction.CustomerName,
			TransactionCode: transaction.Code,
			TotalAmount:     humanize.FormatFloat("#.###,", float64(transaction.Amount)),
			BankName:        transaction.UserBank.Bank.Name,
			PaymentTime:     h.TimeHelper.TimeFormat(time.Now().Local()),
			CheckStatusLink: fmt.Sprintf("%stransaction/%s", h.Config.WebUrl, transaction.Code),
			DetailOrder: helper.DetailOrder{
				ProductName:               transaction.Order.Product.Name,
				ProductVariant:            transaction.Order.ProductDetail.ProductVariant.Name,
				ProductPrice:              humanize.FormatFloat("#.###,", float64(transaction.Order.ProductDetail.Price)),
				ProductQty:                transaction.Order.Qty,
				ProductWeight:             h.UnitHelper.WeightFormat(transaction.Order.Product.Weight),
				ProductTotalPrice:         humanize.FormatFloat("#.###,", float64(productTotalPrice)),
				Courier:                   transaction.Courier,
				CourierType:               transaction.CourierType,
				CourierWeight:             h.UnitHelper.WeightFormat(courierWeight),
				PriceShipment:             humanize.FormatFloat("#.###,", float64(transaction.PriceShipment)),
				TotalPrice:                humanize.FormatFloat("#.###,", float64(totalPrice)),
				TransactionGenerateRandom: transaction.PaymentCode,
				TotalAmount:               humanize.FormatFloat("#.###,", float64(transaction.Amount)),
			},
		}

		go h.Email.SendEmailPaymentSuccess(&templateData, transaction.CustomerEmail, transaction.Order.Product.NonPyshical)
	}

	if input.Status == h.TxConfig.Status.Done {
		// send email
		productTotalPrice := transaction.Order.ProductDetail.Price * transaction.Order.Qty
		courierWeight := transaction.Order.Qty * transaction.Order.Product.Weight
		totalPrice := productTotalPrice + transaction.PriceShipment

		templateData := helper.TransactionDoneTemplateData{
			Fullname: transaction.CustomerName,
			DetailOrder: helper.DetailOrder{
				ProductName:               transaction.Order.Product.Name,
				ProductVariant:            transaction.Order.ProductDetail.ProductVariant.Name,
				ProductPrice:              humanize.FormatFloat("#.###,", float64(transaction.Order.ProductDetail.Price)),
				ProductQty:                transaction.Order.Qty,
				ProductWeight:             h.UnitHelper.WeightFormat(transaction.Order.Product.Weight),
				ProductTotalPrice:         humanize.FormatFloat("#.###,", float64(productTotalPrice)),
				Courier:                   transaction.Courier,
				CourierType:               transaction.CourierType,
				CourierWeight:             h.UnitHelper.WeightFormat(courierWeight),
				PriceShipment:             humanize.FormatFloat("#.###,", float64(transaction.PriceShipment)),
				TotalPrice:                humanize.FormatFloat("#.###,", float64(totalPrice)),
				TransactionGenerateRandom: transaction.PaymentCode,
				TotalAmount:               humanize.FormatFloat("#.###,", float64(transaction.Amount)),
			},
		}

		log.Println("called")
		log.Println(transaction.Order.Product)

		go h.Email.SendEmailTransactionDone(&templateData, transaction.CustomerEmail, transaction.Order.Product.NonPyshical)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Transaction updated.",
	})
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	var transaction model.Transaction
	var updateTransaction model.Transaction

	c.ShouldBind(&updateTransaction)

	UserID, _ := c.Get("UserID")
	transactionID := c.Param("id")

	if err := h.DB.Preload("Shop").Preload("Order.Product").Preload("Order.ProductDetail.ProductVariant").First(&transaction, transactionID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if transaction.Shop.OwnerID != UserID {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Transaksi ini bukan milik toko kamu.",
		})

		return
	}

	if transaction.Status >= updateTransaction.Status {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Status transaksi sudah diubah sebelumnya.",
		})

		return
	}

	if err := h.DB.Model(&transaction).Updates(&updateTransaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Gagal mengubah status transaksi.",
		})

		return
	}

	// socket
	event := fmt.Sprintf("transaction_%d", transaction.ID)

	msg := gin.H{
		"transaction_id": transaction.ID,
		"status":         transaction.Status,
	}

	h.SocketServerHandler.Server.BroadcastTo("transaction", event, msg)

	switch updateTransaction.Status {
	case h.TxConfig.Status.Sending:
		// schedule
		cron := gocron.NewScheduler()
		cron.Every(1).Days().Do(func(cron *gocron.Scheduler) {
			ch := make(chan WaybillResponse)
			go h.ShippingHandler.checkWaybill(updateTransaction.TrackingNumber, transaction.Courier, ch)

			result := <-ch

			if result.Rajaongkir.Result.Delivered == true {
				// send email
				productTotalPrice := transaction.Order.Price * transaction.Order.Qty
				courierWeight := transaction.Order.Qty * transaction.Order.Product.Weight
				totalPrice := productTotalPrice + transaction.PriceShipment

				templateData := helper.TransactionReceivedTemplateData{
					Fullname:        transaction.CustomerName,
					TransactionCode: transaction.Code,
					TrackingNumber:  transaction.TrackingNumber,
					ReceivedDate:    h.TimeHelper.TimeFormat(time.Now().Local()),
					DetailOrder: helper.DetailOrder{
						ProductName:               transaction.Order.Product.Name,
						ProductVariant:            transaction.Order.ProductDetail.ProductVariant.Name,
						ProductPrice:              humanize.FormatFloat("#.###,", float64(transaction.Order.Price)),
						ProductQty:                transaction.Order.Qty,
						ProductWeight:             h.UnitHelper.WeightFormat(transaction.Order.Product.Weight),
						ProductTotalPrice:         humanize.FormatFloat("#.###,", float64(productTotalPrice)),
						Courier:                   transaction.Courier,
						CourierType:               transaction.CourierType,
						CourierWeight:             h.UnitHelper.WeightFormat(courierWeight),
						TotalPrice:                humanize.FormatFloat("#.###,", float64(totalPrice)),
						TransactionGenerateRandom: transaction.PaymentCode,
						TotalAmount:               humanize.FormatFloat("#.###,", float64(transaction.Amount)),
					},
					ShipmentAddress: helper.ShipmentAddress{
						Fullname:    transaction.CustomerName,
						Address:     transaction.CustomerAddress,
						Subdistrict: transaction.CustomerSubdistrict,
						District:    transaction.CustomerCityOrDistrict,
						PostalCode:  transaction.PostalCode,
						Province:    transaction.CustomerProvince,
						Phone:       transaction.CustomerPhoneNumber,
					},
				}

				go h.Email.SendEmailTransactionReceived(&templateData, transaction.CustomerEmail, transaction.Order.Product.NonPyshical)

				// clear schedule
				cron.Clear()
			}
		}, cron)

		go cron.Start()

		// send email
		productTotalPrice := transaction.Order.ProductDetail.Price * transaction.Order.Qty
		courierWeight := transaction.Order.Qty * transaction.Order.Product.Weight
		totalPrice := productTotalPrice + transaction.PriceShipment

		templateData := helper.TransactionSendingTemplateData{
			Fullname:        transaction.CustomerName,
			TransactionCode: transaction.Code,
			ShopName:        transaction.Shop.Name,
			TrackingNumber:  updateTransaction.TrackingNumber,
			Courier:         transaction.Courier,
			CourierType:     transaction.CourierType,
			ShipmentDate:    h.TimeHelper.TimeFormat(time.Now().Local()),
			DetailOrder: helper.DetailOrder{
				ProductName:               transaction.Order.Product.Name,
				ProductVariant:            transaction.Order.ProductDetail.ProductVariant.Name,
				ProductPrice:              humanize.FormatFloat("#.###,", float64(transaction.Order.ProductDetail.Price)),
				ProductQty:                transaction.Order.Qty,
				ProductWeight:             h.UnitHelper.WeightFormat(transaction.Order.Product.Weight),
				ProductTotalPrice:         humanize.FormatFloat("#.###,", float64(productTotalPrice)),
				Courier:                   transaction.Courier,
				CourierType:               transaction.CourierType,
				CourierWeight:             h.UnitHelper.WeightFormat(courierWeight),
				PriceShipment:             humanize.FormatFloat("#.###,", float64(transaction.PriceShipment)),
				TotalPrice:                humanize.FormatFloat("#.###,", float64(totalPrice)),
				TransactionGenerateRandom: transaction.PaymentCode,
				TotalAmount:               humanize.FormatFloat("#.###,", float64(transaction.Amount)),
			},
			ShipmentAddress: helper.ShipmentAddress{
				Fullname:    transaction.CustomerName,
				Address:     transaction.CustomerAddress,
				Subdistrict: transaction.CustomerSubdistrict,
				District:    transaction.CustomerCityOrDistrict,
				PostalCode:  transaction.PostalCode,
				Province:    transaction.CustomerProvince,
				Phone:       transaction.CustomerPhoneNumber,
			},
		}

		go h.Email.SendEmailTransactionSending(&templateData, transaction.CustomerEmail, transaction.Order.Product.NonPyshical)
	case h.TxConfig.Status.Received:
		cron := gocron.NewScheduler()
		cron.Every(2).Days().Do(func(cron *gocron.Scheduler) {
			// keep
			go h.DB.Model(&transaction).Update("status", h.TxConfig.Status.Done)
		}, cron)

		go cron.Start()
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "Berhasil mengubah data transaksi.",
	})
}

func NewTransactionHandler(
	database *gorm.DB,
	config *config.Config,
	txConfig *config.TransactionConfig,
	email *helper.Email,
	shipping *ShippingHandler,
	socketServerHandler *SocketServerHandler,
	timeHelper *helper.TimeHelper,
	unitHelper *helper.UnitHelper,
) *TransactionHandler {
	return &TransactionHandler{
		DB:                  database,
		Config:              config,
		TxConfig:            txConfig,
		Email:               email,
		ShippingHandler:     shipping,
		SocketServerHandler: socketServerHandler,
		TimeHelper:          timeHelper,
		UnitHelper:          unitHelper,
	}
}
