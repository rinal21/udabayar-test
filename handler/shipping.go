package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"udabayar-go-api-di/config"
	"udabayar-go-api-di/model"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

type ShippingHandler struct {
	Config   *config.Config
	Database *gorm.DB
}

type ProvinceResponse struct {
	Rajaongkir struct {
		Query struct {
			ID string `json:"id"`
		} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Results struct {
			ProvinceID string `json:"province_id"`
			Province   string `json:"province"`
		} `json:"results"`
	} `json:"rajaongkir"`
}

func (h *ShippingHandler) Province(c *gin.Context) {
	var response ProvinceResponse

	url := h.Config.RajaOngkirAPIUrl + "/starter/cost"
	payload := strings.NewReader("origin=501&destination=114&weight=1700&courier=jne")

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("key", h.Config.RajaOngkirAPIKey)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "Gagal mengambil harga pengiriman, silakan coba kembali.",
		})

		return
	}

	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &response)

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   response.Rajaongkir.Results,
	})
}

type CostRequest struct {
	ShopID          int `form:"shop_id" json:"shop_id" binding:"required"`
	DestinationCode int `form:"destination_code" json:"destination_code" binding:"required"`
	Weight          int `form:"weight" json:"weight" binding:"required"`
}

type CostResponse struct {
	Rajaongkir struct {
		Query struct {
			Origin      string `json:"origin"`
			Destination string `json:"destination"`
			Weight      int    `json:"weight"`
			Courier     string `json:"courier"`
		} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		OriginDetails struct {
			CityID     string `json:"city_id"`
			ProvinceID string `json:"province_id"`
			Province   string `json:"province"`
			Type       string `json:"type"`
			CityName   string `json:"city_name"`
			PostalCode string `json:"postal_code"`
		} `json:"origin_details"`
		DestinationDetails struct {
			CityID     string `json:"city_id"`
			ProvinceID string `json:"province_id"`
			Province   string `json:"province"`
			Type       string `json:"type"`
			CityName   string `json:"city_name"`
			PostalCode string `json:"postal_code"`
		} `json:"destination_details"`
		Results []struct {
			Code  string `json:"code"`
			Name  string `json:"name"`
			Costs []struct {
				Service     string `json:"service"`
				Description string `json:"description"`
				Cost        []struct {
					Value int    `json:"value"`
					Etd   string `json:"etd"`
					Note  string `json:"note"`
				} `json:"cost"`
			} `json:"costs"`
		} `json:"results"`
	} `json:"rajaongkir"`
}

type GetCostInput struct {
	Origin          int    `json:"origin"`
	OriginType      string `json:"originType"`
	Destination     int    `json:"destination"`
	DestinationType string `json:"destinationType"`
	Weight          int    `json:"weight"`
	Courier         string `json:"courier"`
}

func (h *ShippingHandler) getCost(origin int, destination int, weight int, courier string, c chan CostResponse) {
	var response CostResponse

	apiUrl := h.Config.RajaOngkirAPIUrl + "/api/cost"
	payload := strings.NewReader(
		fmt.Sprintf(
			"origin=%d&originType=city&destination=%d&destinationType=subdistrict&weight=%d&courier=%s",
			origin, destination, weight, strings.ToLower(courier),
		),
	)

	req, _ := http.NewRequest("POST", apiUrl, payload)
	req.Header.Add("key", h.Config.RajaOngkirAPIKey)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &response)

	c <- response
}

func (h *ShippingHandler) Cost(c *gin.Context) {
	var request CostRequest
	var shop model.Shop
	var address model.Address

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if check := h.Database.Preload("ShopCouriers.Courier").First(&shop, request.ShopID).RecordNotFound(); check == true {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "Toko tidak tersedia.",
		})

		return
	}

	if check := h.Database.First(&address, shop).RecordNotFound(); check == true {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "Alamat tidak tersedia.",
		})

		return
	}

	ch := make(chan CostResponse, 1)

	for _, shopCourier := range shop.ShopCouriers {
		go h.getCost(address.OriginCode, request.DestinationCode, request.Weight, shopCourier.Courier.Code, ch)
	}

	var response []CostResponse
	// status := true

	for i := 0; i < len(shop.ShopCouriers); i++ {
		res := <-ch
		response = append(response, res)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   response,
	})
}

type WaybillRequest struct {
	Waybill string
	Courier string
}

type WaybillResponse struct {
	Rajaongkir struct {
		Query struct {
			Waybill string `json:"waybill"`
			Courier string `json:"courier"`
		} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Result struct {
			Delivered bool `json:"delivered"`
			Summary   struct {
				CourierCode   string `json:"courier_code"`
				CourierName   string `json:"courier_name"`
				WaybillNumber string `json:"waybill_number"`
				ServiceCode   string `json:"service_code"`
				WaybillDate   string `json:"waybill_date"`
				ShipperName   string `json:"shipper_name"`
				ReceiverName  string `json:"receiver_name"`
				Origin        string `json:"origin"`
				Destination   string `json:"destination"`
				Status        string `json:"status"`
			} `json:"summary"`
			Details struct {
				WaybillNumber    string `json:"waybill_number"`
				WaybillDate      string `json:"waybill_date"`
				WaybillTime      string `json:"waybill_time"`
				Weight           string `json:"weight"`
				Origin           string `json:"origin"`
				Destination      string `json:"destination"`
				ShippperName     string `json:"shippper_name"`
				ShipperAddress1  string `json:"shipper_address1"`
				ShipperAddress2  string `json:"shipper_address2"`
				ShipperAddress3  string `json:"shipper_address3"`
				ShipperCity      string `json:"shipper_city"`
				ReceiverName     string `json:"receiver_name"`
				ReceiverAddress1 string `json:"receiver_address1"`
				ReceiverAddress2 string `json:"receiver_address2"`
				ReceiverAddress3 string `json:"receiver_address3"`
				ReceiverCity     string `json:"receiver_city"`
			} `json:"details"`
			DeliveryStatus struct {
				Status      string      `json:"status"`
				PodReceiver interface{} `json:"pod_receiver"`
				PodDate     interface{} `json:"pod_date"`
				PodTime     interface{} `json:"pod_time"`
			} `json:"delivery_status"`
			Manifest []struct {
				ManifestCode        string `json:"manifest_code"`
				ManifestDescription string `json:"manifest_description"`
				ManifestDate        string `json:"manifest_date"`
				ManifestTime        string `json:"manifest_time"`
				CityName            string `json:"city_name"`
			} `json:"manifest"`
		} `json:"result"`
	} `json:"rajaongkir"`
}

func (h *ShippingHandler) checkWaybill(waybill string, courier string, c chan WaybillResponse) {
	var response WaybillResponse

	apiUrl := h.Config.RajaOngkirAPIUrl + "/api/waybill"
	payload := strings.NewReader(
		fmt.Sprintf(
			"waybill=%s&courier=%s",
			waybill, strings.ToLower(courier),
		),
	)

	req, _ := http.NewRequest("POST", apiUrl, payload)
	req.Header.Add("key", h.Config.RajaOngkirAPIKey)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &response)

	c <- response
}

type WaybillInput struct {
	TransactionID int `form:"transaction_id" json:"transaction_id" binding:"required"`
}

func (h *ShippingHandler) Waybill(c *gin.Context) {
	var request WaybillInput
	var transaction model.Transaction

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  0,
			"message": err.Error(),
		})

		return
	}

	if err := h.Database.First(&transaction, request.TransactionID).RecordNotFound(); err == true {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Transaksi tidak tersedia.",
		})
	}

	if transaction.TrackingNumber == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "Transaksi ini tidak ada nomor resi.",
		})
	}

	ch := make(chan WaybillResponse)
	go h.checkWaybill(transaction.TrackingNumber, transaction.Courier, ch)

	response := <-ch

	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   response.Rajaongkir,
	})
}

func NewShippingHandler(config *config.Config, database *gorm.DB) *ShippingHandler {
	return &ShippingHandler{
		Config:   config,
		Database: database,
	}
}
