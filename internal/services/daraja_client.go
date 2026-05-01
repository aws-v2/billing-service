package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Qarani-m/billing-service/internal/domain"
)

func baseURL(c domain.DarajaConfig) string {
	if c.Environment == "prod" {
		return "https://api.safaricom.co.ke"
	}
	return "https://sandbox.safaricom.co.ke"
}
// darajaClient manages OAuth tokens and wraps Daraja API calls.
type DarajaClient struct {
	cfg        domain.DarajaConfig
	httpClient *http.Client

	mu          sync.Mutex
	accessToken string
	tokenExpiry time.Time
}

func NewDarajaClient(cfg domain.DarajaConfig) *DarajaClient {
	return &DarajaClient{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ── OAuth ──────────────────────────────────────────────────────────────────

type oauthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

// token returns a valid access token, refreshing if expired.
func (d *DarajaClient) token() (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.accessToken != "" && time.Now().Before(d.tokenExpiry) {
		return d.accessToken, nil
	}

	creds := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", d.cfg.ConsumerKey, d.cfg.ConsumerSecret)),
	)

	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/oauth/v1/generate?grant_type=client_credentials", d.cfg.DarajaBaseURL),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("daraja: build token request: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+creds)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("daraja: token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("daraja: token HTTP %d: %s", resp.StatusCode, body)
	}

	var result oauthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("daraja: decode token response: %w", err)
	}

	// Expire 60 s before the real TTL to avoid edge-case races
	d.accessToken = result.AccessToken
	d.tokenExpiry = time.Now().Add(55 * time.Minute)
	return d.accessToken, nil
}

// ── STK Push ───────────────────────────────────────────────────────────────

type stkPushPayload struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            int    `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

// StkPushResult is the raw Daraja STK push response.
type StkPushResult struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

// SendStkPush dispatches an STK push via Daraja and returns the result.
func (d *DarajaClient) SendStkPush(phone string, amount int) (*StkPushResult, error) {
	token, err := d.token()
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("20060102150405")
	rawPassword := fmt.Sprintf("%s%s%s", d.cfg.ShortCode, d.cfg.PassKey, timestamp)
	password := base64.StdEncoding.EncodeToString([]byte(rawPassword))

	// Normalize phone: strip leading + so Daraja gets "254XXXXXXXXX"
	normalizedPhone := strings.TrimPrefix(phone, "+")

	payload := stkPushPayload{
		BusinessShortCode: d.cfg.ShortCode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            amount,
		PartyA:            normalizedPhone,
		PartyB:            d.cfg.ShortCode,
		PhoneNumber:       normalizedPhone,
		CallBackURL:       d.cfg.CallbackURL,
		AccountReference:  "SERWIN-VERIFY",
		TransactionDesc:   "Billing verification",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("daraja: marshal stk payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/mpesa/stkpush/v1/processrequest", baseURL(d.cfg)),
		strings.NewReader(string(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("daraja: build stk request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("daraja: stk push request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("daraja: stk push HTTP %d: %s", resp.StatusCode, raw)
	}

	var result StkPushResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("daraja: decode stk response: %w", err)
	}
	if result.ResponseCode != "0" {
		return nil, fmt.Errorf("daraja: stk push rejected: %s", result.ResponseDescription)
	}

	return &result, nil
}

// ── STK Query ──────────────────────────────────────────────────────────────

type stkQueryPayload struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	CheckoutRequestID string `json:"CheckoutRequestID"`
}

 

// QueryStkStatus polls Daraja for the status of a previously initiated push.
func (d *DarajaClient) QueryStkStatus(checkoutRequestID string) (*domain.StkQueryResult, error) {
	token, err := d.token()
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("20060102150405")
	rawPassword := fmt.Sprintf("%s%s%s", d.cfg.ShortCode, d.cfg.PassKey, timestamp)
	password := base64.StdEncoding.EncodeToString([]byte(rawPassword))

	payload := stkQueryPayload{
		BusinessShortCode: d.cfg.ShortCode,
		Password:          password,
		Timestamp:         timestamp,
		CheckoutRequestID: checkoutRequestID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("daraja: marshal query payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/mpesa/stkpushquery/v1/query", baseURL(d.cfg)),
		strings.NewReader(string(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("daraja: build query request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("daraja: query request failed: %w", err)
	}
	defer resp.Body.Close()

	var result domain.StkQueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("daraja: decode query response: %w", err)
	}

	return &result, nil
}