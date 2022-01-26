package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Payload struct {
	BusinessShortCode int    `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            int    `json:"Amount"`
	PartyA            int    `json:"PartyA"`
	PartyB            int    `json:"PartyB"`
	PhoneNumber       int    `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

type Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

const businessShortCode int = 174379

func STKPush() error {
	url := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
	method := "POST"

	passKey := "bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919"
	timestamp := time.Now().Format("20060102150405")

	pw := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(businessShortCode) + passKey + timestamp))

	payload := &Payload{
		BusinessShortCode: businessShortCode,
		Password:          pw,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            1,
		PartyA:            254708374149,
		PartyB:            businessShortCode,
		PhoneNumber:       254719825151,
		CallBackURL:       "https://hudumaapp.herokuapp.com/transactions/confirm",
		AccountReference:  "CompanyXLTD",
		TransactionDesc:   "Payment of X",
	}

	pl, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return err
	}
	pldStr := strings.NewReader(string(pl[:]))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, pldStr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	token, err := accessToken()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}

func GetPayment() error {
	return nil
}

func accessToken() (string, error) {
	url := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	consumerKey := "OS7jGKun8KgrybHTr5koshJPzSUAligA"
	consumerSecret := "WrZhzVxkakde91Ga"

	token := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	fmt.Println(string(body))

	var response Response
	json.Unmarshal(body, &response)
	fmt.Println("Access Token:", response.AccessToken)

	return response.AccessToken, nil
}

func C2BRegisterURL() error {
	body, err := json.Marshal(registerURL{
		ConfirmationURL: "https://hudumaapp.herokuapp.com/transactions/confirm",
		ValidationURL:   "https://hudumaapp.herokuapp.com/transactions/validate",
		ShortCode:       600979,
		ResponseType:    "completed",
	})
	if err != nil {
		return err
	}

	auth, err := accessToken()
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["Cache-Control"] = "no-cache"

	url := baseURL() + "mpesa/c2b/v1/registerurl"
	fmt.Println(newReq(url, body, headers))
	return nil
}

func C2BSimulate() error {
	body, err := json.Marshal(C2B{
		ShortCode:     600576,
		CommandID:     "CustomerPayBillOnline",
		Amount:        2,
		Msisdn:        "254708374149",
		BillRefNumber: "hkjhjkhjkh"})
	if err != nil {
		return err
	}

	auth, err := accessToken()
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth

	url := baseURL() + "mpesa/c2b/v1/simulate"
	fmt.Println(newReq(url, body, headers))
	return nil
}

func newReq(url string, body []byte, headers map[string]string) (string, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", nil
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(request)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}

	stringBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(stringBody), nil
}

func baseURL() string {
	return "https://sandbox.safaricom.co.ke/"
}
