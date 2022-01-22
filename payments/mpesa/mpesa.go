package mpesa

import (
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

func MakePayment() error {
	url := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"
	method := "POST"

	businessShortCode := 174379
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

	accessToken := getAccessToken()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

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

func getAccessToken() string {
	url := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
	}

	consumerKey := "OS7jGKun8KgrybHTr5koshJPzSUAligA"
	consumerSecret := "WrZhzVxkakde91Ga"

	token := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	fmt.Println(string(body))

	var response Response
	json.Unmarshal(body, &response)
	fmt.Println("Access Token:", response.AccessToken)

	return response.AccessToken
}
