package helper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
	"udabayar-go-api-di/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Email struct {
	Config *config.EmailConfig
	SES    *ses.SES
}

type EmailInput struct {
	Template     string
	Recipient    string
	TemplateData string
}

func (h *Email) CreateTemplate() {
	h.DeleteTemplate()

	registerBody, _ := ioutil.ReadFile("helper/emailtemplate/register.html")
	activatedAccountBody, _ := ioutil.ReadFile("helper/emailtemplate/activated_account.html")
	createTransactionBody, _ := ioutil.ReadFile("helper/emailtemplate/create_transaction.html")
	paymentSuccessBody, _ := ioutil.ReadFile("helper/emailtemplate/payment_success.html")
	transactionSendingBody, _ := ioutil.ReadFile("helper/emailtemplate/transaction_sending.html")
	transactionReceivedBody, _ := ioutil.ReadFile("helper/emailtemplate/transaction_received.html")
	transactionDoneBody, _ := ioutil.ReadFile("helper/emailtemplate/transaction_done.html")
	transactionCancelledBody, _ := ioutil.ReadFile("helper/emailtemplate/transaction_cancelled.html")

	templates := []ses.CreateTemplateInput{
		{
			Template: &ses.Template{
				HtmlPart:     aws.String(string(registerBody)),
				SubjectPart:  aws.String("Selamat Datang di Udabayar"),
				TemplateName: aws.String("RegistrationTemplate"),
			},
		},
		{
			Template: &ses.Template{
				HtmlPart:     aws.String(string(activatedAccountBody)),
				SubjectPart:  aws.String("Akun anda telah aktif"),
				TemplateName: aws.String("ActivatedAccountTemplate"),
			},
		},
		{
			Template: &ses.Template{
				HtmlPart:     aws.String(string(createTransactionBody)),
				SubjectPart:  aws.String("Segera selesaikan pembayaran anda"),
				TemplateName: aws.String("CreateTransactionTemplate"),
			},
		},
		{
			Template: &ses.Template{
				HtmlPart:     aws.String(string(paymentSuccessBody)),
				SubjectPart:  aws.String("Pembayaran anda telah kami terima"),
				TemplateName: aws.String("PaymentSuccessTemplate"),
			},
		},
		{
			Template: &ses.Template{
				HtmlPart:     aws.String(string(paymentSuccessBody)),
				SubjectPart:  aws.String("Pembayaran anda telah kami terima"),
				TemplateName: aws.String("PaymentSuccessTemplate"),
			},
		},
		{
			Template: &ses.Template{
				HtmlPart:     aws.String(string(transactionSendingBody)),
				SubjectPart:  aws.String("Pesanan anda sedang dalam perjalanan."),
				TemplateName: aws.String("TransactionSendingTemplate"),
			},
		},
		{
			Template: &ses.Template{
				HtmlPart:     aws.String(string(transactionReceivedBody)),
				SubjectPart:  aws.String("Pesanan telah diterima."),
				TemplateName: aws.String("TransactionReceivedTemplate"),
			},
		},
		{
			Template: &ses.Template{
				HtmlPart:     aws.String(string(transactionDoneBody)),
				SubjectPart:  aws.String("Transaksi selesai."),
				TemplateName: aws.String("TransactionDoneTemplate"),
			},
		},
		{
			Template: &ses.Template{
				HtmlPart:     aws.String(string(transactionCancelledBody)),
				SubjectPart:  aws.String("Transaksi anda telah kami dibatalkan."),
				TemplateName: aws.String("TransactionCancelledTemplate"),
			},
		},
	}

	for _, template := range templates {
		h.SES.CreateTemplate(&template)
	}
}

func (h *Email) DeleteTemplate() {
	templates := []ses.DeleteTemplateInput{
		{
			TemplateName: aws.String("RegistrationTemplate"),
		},
		{
			TemplateName: aws.String("ActivatedAccountTemplate"),
		},
		{
			TemplateName: aws.String("CreateTransactionTemplate"),
		},
		{
			TemplateName: aws.String("PaymentSuccessTemplate"),
		},
		{
			TemplateName: aws.String("TransactionSendingTemplate"),
		},
		{
			TemplateName: aws.String("TransactionReceivedTemplate"),
		},
		{
			TemplateName: aws.String("TransactionDoneTemplate"),
		},
		{
			TemplateName: aws.String("TransactionCancelledTemplate"),
		},
	}

	for _, template := range templates {
		if _, err := h.SES.DeleteTemplate(&template); err != nil {
			panic(err)
		}
	}
}

type DetailOrder struct {
	ProductName               string `json:"product_name"`
	ProductVariant            string `json:"product_variant"`
	ProductPrice              string `json:"product_price"`
	ProductQty                int    `json:"product_qty"`
	ProductWeight             string `json:"product_weight"`
	ProductTotalPrice         string `json:"product_total_price"`
	Courier                   string `json:"courier"`
	CourierType               string `json:"courier_type"`
	CourierWeight             string `json:"courier_weight"`
	PriceShipment             string `json:"price_shipment"`
	TotalPrice                string `json:"total_price"`
	TransactionGenerateRandom int    `json:"transaction_generate_random"`
	TotalAmount               string `json:"total_amount"`
}

type ShipmentAddress struct {
	Fullname    string `json:"fullname"`
	Address     string `json:"customer_address"`
	Subdistrict string `json:"subdistrict"`
	District    string `json:"district"`
	PostalCode  int    `json:"postal_code"`
	Province    string `json:"province"`
	Phone       string `json:"customer_phone"`
}

type RegistrationTemplateData struct {
	Fullname       string    `json:"fullname"`
	ActivationLink string    `json:"activation_link"`
	ExpiredAt      time.Time `json:"expired_at"`
}

func (h *Email) SendEmailRegistration(templateData *RegistrationTemplateData, recipient string) {
	templateDataJSON, _ := json.Marshal(templateData)

	input := &EmailInput{
		Template:     "RegistrationTemplate",
		Recipient:    recipient,
		TemplateData: string(templateDataJSON),
	}

	h.SendEmail(input)
}

type ActivatedAccountTemplateData struct {
	Fullname  string `json:"fullname"`
	LoginLink string `json:"login_link"`
}

func (h *Email) SendEmailActivatedAccount(templateData *ActivatedAccountTemplateData, recipient string) {
	templateDataJSON, _ := json.Marshal(templateData)

	input := &EmailInput{
		Template:     "ActivatedAccountTemplate",
		Recipient:    recipient,
		TemplateData: string(templateDataJSON),
	}

	h.SendEmail(input)
}

type CreateTransactionTemplateData struct {
	Fullname         string `json:"fullname"`
	TransactionCode  string `json:"transaction_code"`
	BankName         string `json:"bank_name"`
	PaymentExpiredAt string `json:"payment_expired_at"`
	TotalAmount      string `json:"total_amount"`
	BankLogo         string `json:"bank_logo"`
	AccountNumber    string `json:"account_number"`
	CheckStatusLink  string `json:"check_status_link"`
	DetailOrder
}

type CreateTransactionPsyhicalTemplateData struct {
	Pyshical                      int `json:"pyshical"`
	CreateTransactionTemplateData *CreateTransactionTemplateData
}

func (h *Email) SendEmailCreateTransaction(templateData *CreateTransactionTemplateData, recipient string, nonPyshical int) {
	var templateDataJSON []byte

	if nonPyshical == 0 {
		templateDataPyshical := CreateTransactionPsyhicalTemplateData{
			Pyshical:                      1,
			CreateTransactionTemplateData: templateData,
		}

		templateDataJSON, _ = json.Marshal(templateDataPyshical)
	} else {
		templateDataJSON, _ = json.Marshal(templateData)
	}

	log.Println(string(templateDataJSON))

	input := &EmailInput{
		Template:     "CreateTransactionTemplate",
		Recipient:    recipient,
		TemplateData: string(templateDataJSON),
	}

	h.SendEmail(input)
}

type PaymentSuccessTemplateData struct {
	Fullname        string `json:"fullname"`
	TransactionCode string `json:"transaction_code"`
	TotalAmount     string `json:"total_amount"`
	BankName        string `json:"bank_name"`
	PaymentTime     string `json:"payment_time"`
	CheckStatusLink string `json:"check_status_link"`
	DetailOrder
}

type PaymentSuccessPyshicalTemplateData struct {
	Pyshical                   int `json:"pyshical"`
	PaymentSuccessTemplateData *PaymentSuccessTemplateData
}

func (h *Email) SendEmailPaymentSuccess(templateData *PaymentSuccessTemplateData, recipient string, nonPyshical int) {
	var templateDataJSON []byte

	if nonPyshical == 0 {
		templateDataPyshical := PaymentSuccessPyshicalTemplateData{
			Pyshical:                   1,
			PaymentSuccessTemplateData: templateData,
		}

		templateDataJSON, _ = json.Marshal(templateDataPyshical)
	} else {
		templateDataJSON, _ = json.Marshal(templateData)
	}

	input := &EmailInput{
		Template:     "PaymentSuccessTemplate",
		Recipient:    recipient,
		TemplateData: string(templateDataJSON),
	}

	h.SendEmail(input)
}

type TransactionSendingTemplateData struct {
	Fullname        string `json:"fullname"`
	TransactionCode string `json:"transaction_code"`
	ShopName        string `json:"shop_name"`
	TrackingNumber  string `json:"tracking_number"`
	Courier         string `json:"courier"`
	CourierType     string `json:"courier_type"`
	ShipmentDate    string `json:"shipment_date"`
	DetailOrder
	ShipmentAddress
}

type TransactionSendingPyshicalTemplateData struct {
	Pyshical                       int `json:"pyshical"`
	TransactionSendingTemplateData *TransactionSendingTemplateData
}

func (h *Email) SendEmailTransactionSending(templateData *TransactionSendingTemplateData, recipient string, nonPyshical int) {
	var templateDataJSON []byte

	if nonPyshical == 0 {
		templateDataPyshical := TransactionSendingPyshicalTemplateData{
			Pyshical:                       1,
			TransactionSendingTemplateData: templateData,
		}

		templateDataJSON, _ = json.Marshal(templateDataPyshical)
	} else {
		templateDataJSON, _ = json.Marshal(templateData)
	}

	input := &EmailInput{
		Template:     "TransactionSendingTemplate",
		Recipient:    recipient,
		TemplateData: string(templateDataJSON),
	}

	log.Println(recipient)

	h.SendEmail(input)
}

type TransactionReceivedTemplateData struct {
	Fullname        string `json:"fullname"`
	TransactionCode string `json:"transaction_code"`
	ShopName        string `json:"shop_name"`
	TrackingNumber  string `json:"tracking_number"`
	ReceivedDate    string `json:"date_received"`
	DetailOrder
	ShipmentAddress
}

type TransactionReceivedPyshicalTemplateData struct {
	Pyshical                        int `json:"pyshical"`
	TransactionReceivedTemplateData *TransactionReceivedTemplateData
}

func (h *Email) SendEmailTransactionReceived(templateData *TransactionReceivedTemplateData, recipient string, nonPyshical int) {
	var templateDataJSON []byte

	if nonPyshical == 0 {
		templateDataPyshical := TransactionReceivedPyshicalTemplateData{
			Pyshical:                        1,
			TransactionReceivedTemplateData: templateData,
		}

		templateDataJSON, _ = json.Marshal(templateDataPyshical)
	} else {
		templateDataJSON, _ = json.Marshal(templateData)
	}

	input := &EmailInput{
		Template:     "TransactionReceivedTemplate",
		Recipient:    recipient,
		TemplateData: string(templateDataJSON),
	}

	h.SendEmail(input)
}

type TransactionDoneTemplateData struct {
	Fullname string `json:"fullname"`
	DetailOrder
}

type TransactionDonePyshicalTemplateData struct {
	Pyshical                    int `json:"pyshical"`
	TransactionDoneTemplateData *TransactionDoneTemplateData
}

func (h *Email) SendEmailTransactionDone(templateData *TransactionDoneTemplateData, recipient string, nonPyshical int) {
	var templateDataJSON []byte

	if nonPyshical == 0 {
		templateDataPyshical := TransactionDonePyshicalTemplateData{
			Pyshical:                    1,
			TransactionDoneTemplateData: templateData,
		}

		templateDataJSON, _ = json.Marshal(templateDataPyshical)
	} else {
		templateDataJSON, _ = json.Marshal(templateData)
	}

	input := &EmailInput{
		Template:     "TransactionDoneTemplate",
		Recipient:    recipient,
		TemplateData: string(templateDataJSON),
	}

	log.Println(recipient)

	h.SendEmail(input)
}

type TransactionCancelledTemplateData struct {
	Fullname        string `json:"fullname"`
	TransactionCode string `json:"transaction_code"`
	TotalAmount     string `json:"total_amount"`
	BankName        string `json:"bank_name"`
	DetailOrder
}

type TransactionCancelledPyshicalTemplateData struct {
	Pyshical                         int `json:"pyshical"`
	TransactionCancelledTemplateData *TransactionCancelledTemplateData
}

func (h *Email) SendEmailTransactionCancelled(templateData *TransactionCancelledTemplateData, recipient string, nonPyshical int) {
	var templateDataJSON []byte

	if nonPyshical == 0 {
		templateDataPyshical := TransactionCancelledPyshicalTemplateData{
			Pyshical:                         1,
			TransactionCancelledTemplateData: templateData,
		}

		templateDataJSON, _ = json.Marshal(templateDataPyshical)
	} else {
		templateDataJSON, _ = json.Marshal(templateData)
	}

	input := &EmailInput{
		Template:     "TransactionCancelledTemplate",
		Recipient:    recipient,
		TemplateData: string(templateDataJSON),
	}

	h.SendEmail(input)
}

func (h *Email) SendEmail(dataInput *EmailInput) {
	input := &ses.SendTemplatedEmailInput{
		Source:   aws.String(h.Config.Sender),
		Template: aws.String(dataInput.Template),
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(dataInput.Recipient),
			},
		},
		TemplateData: aws.String(dataInput.TemplateData),
	}

	_, err := h.SES.SendTemplatedEmail(input)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func NewEmailHelper(emailconfig *config.EmailConfig) *Email {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	if err != nil {
		panic(err)
	}

	ses := ses.New(session)

	return &Email{
		Config: emailconfig,
		SES:    ses,
	}
}
