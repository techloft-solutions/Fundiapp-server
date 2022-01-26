package mpesa

type registerURL struct {
	ShortCode       int
	ResponseType    string
	ConfirmationURL string
	ValidationURL   string
}

type C2B struct {
	ShortCode     int
	CommandID     string
	Amount        int
	Msisdn        string
	BillRefNumber string
}
