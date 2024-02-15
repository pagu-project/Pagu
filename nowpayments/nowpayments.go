package nowpayments

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kehiy/RoboPac/store"
	"github.com/pactus-project/pactus/util/logger"
)

type NowPayments struct {
	apiToken  string
	ipnSecret []byte
	webhook   string
	apiURL    string
	username  string
	password  string
}

func NewNowPayments(cfg *Config) (*NowPayments, error) {
	ipnSecret, err := base64.StdEncoding.DecodeString(cfg.IPNSecret)
	if err != nil {
		return nil, err
	}
	s := &NowPayments{
		apiToken:  cfg.APIToken,
		ipnSecret: ipnSecret,
		apiURL:    cfg.APIUrl,
		webhook:   cfg.Webhook,
		username:  cfg.Username,
		password:  cfg.Password,
	}
	http.HandleFunc("/nowpayments", s.webhookFunc)

	go func() {
		for {
			logger.Info("starting NowPayments webhook", "port", cfg.ListenPort)
			err = http.ListenAndServe(fmt.Sprintf(":%v", cfg.ListenPort), nil)
			if err != nil {
				logger.Error("unable to start NowPayments webhook", "error", err)
			}
		}
	}()

	return s, nil
}

func (s *NowPayments) webhookFunc(w http.ResponseWriter, r *http.Request) {
	logger.Debug("NowPayment webhook called")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("Callback read error", "error", err)
		return
	}

	logger.Debug("Callback result", "data", data)
	msgMACHex := r.Header.Get("x-nowpayments-sig")
	msgMAC, err := hex.DecodeString(msgMACHex)
	if err != nil {
		logger.Error("Invalid sig hex", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mac := hmac.New(sha512.New, s.ipnSecret)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("json.Unmarshal read error", "error", err)
		return
	}
	// marshal it again, it will be sorted
	sortedData, _ := json.Marshal(result)

	if !bytes.Equal(sortedData, data) {
		logger.Debug("data was not sorted")
	}

	_, err = mac.Write(sortedData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("mac.Write read error", "error", err)
		return
	}
	expectedMAC := mac.Sum(nil)
	if !hmac.Equal(expectedMAC, msgMAC) {
		/// TODO: fix me
		// w.WriteHeader(http.StatusBadRequest)
		logger.Error("HMAC is invalid", "expectedMAC", expectedMAC, "msgMAC", msgMAC)
		// return
	}

	w.WriteHeader(http.StatusOK)
}

// curl --location 'https://api.nowpayments.io/v1/invoice' \
// --header 'x-api-key: {{api-key}}' \
// --header 'Content-Type: application/json' \
//
//	--data '{
//	  "price_amount": 1000,
//	  "price_currency": "usd",
//	  "order_id": "RGDBP-21314",
//	  "order_description": "Apple Macbook Pro 2019 x 1",
//	  "ipn_callback_url": "https://nowpayments.io",
//	  "success_url": "https://nowpayments.io",
//	  "cancel_url": "https://nowpayments.io"
//	}

// return
//
//	{
//		"id": "4522625843",
//		"order_id": "RGDBP-21314",
//		"order_description": "Apple Macbook Pro 2019 x 1",
//		"price_amount": "1000",
//		"price_currency": "usd",
//		"pay_currency": null,
//		"ipn_callback_url": "https://nowpayments.io",
//		"invoice_url": "https://nowpayments.io/payment/?iid=4522625843",
//		"success_url": "https://nowpayments.io",
//		"cancel_url": "https://nowpayments.io",
//		"created_at": "2020-12-22T15:05:58.290Z",
//		"updated_at": "2020-12-22T15:05:58.290Z"
//	  }
func (s *NowPayments) CreatePayment(party *store.TwitterParty) error {
	url := fmt.Sprintf("%v/v1/invoice", s.apiURL)
	// jsonStr := fmt.Sprintf(`{"price_amount":%v,"price_currency":"usd","ipn_callback_url":"%v","order_id":"%v"}`,
	// 	party.TotalPrice, s.webhook, party.DiscountCode)

	jsonStr := fmt.Sprintf(`{"price_amount":%v,"price_currency":"usd","order_id":"%v"}`,
		party.TotalPrice, party.DiscountCode)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiToken)

	logger.Info("calling NowPayments:CreatePayment", "twitter", party.TwitterName, "json", jsonStr)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logger.Debug("CreatePayment Response", "res", string(data))
	// http.StatusOK = 200
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to call NowPayments. Status code: %v, status: %v", resp.StatusCode, resp.Status)
	}

	var resultJSON map[string]interface{}
	err = json.Unmarshal(data, &resultJSON)
	if err != nil {
		return err
	}

	// fmt.Println(string(data))
	party.NowPaymentsInvoiceID = resultJSON["id"].(string)

	return nil
}

func (s *NowPayments) UpdatePayment(party *store.TwitterParty) error {
	token, err := s.getJWTToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%v/v1/payment/?invoiceId=%v",
		s.apiURL, party.NowPaymentsInvoiceID)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", s.apiToken)
	req.Header.Set("Authorization", "Bearer "+token)

	logger.Info("calling NowPayments:ListOfPayments", "Twitter", party.TwitterName, "NowPaymentsInvoiceID", party.NowPaymentsInvoiceID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logger.Debug("ListOfPayments Response", "res", string(data))
	// http.StatusOK = 200
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to call NowPayments:Payment. Status code: %v", resp.StatusCode)
	}

	var resultJSON map[string]interface{}
	err = json.Unmarshal(data, &resultJSON)
	if err != nil {
		return err
	}

	results := resultJSON["data"].([]interface{})
	for _, payment := range results {
		paymentStatus := payment.(map[string]interface{})["payment_status"]

		if paymentStatus == "finished" {
			party.NowPaymentsFinished = true
		}
	}

	return nil
}

func (s *NowPayments) getJWTToken() (string, error) {
	url := fmt.Sprintf("%v/v1/auth", s.apiURL)
	jsonStr := fmt.Sprintf(`{"email":"%v","password":"%v"}`, s.username, s.password)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	logger.Info("calling NowPayments:auth")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// http.StatusOK = 200
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to call Auth. Status code: %v, status: %v", resp.StatusCode, resp.Status)
	}

	var resultJSON map[string]interface{}
	err = json.Unmarshal(data, &resultJSON)
	if err != nil {
		return "", err
	}

	return resultJSON["token"].(string), nil
}
