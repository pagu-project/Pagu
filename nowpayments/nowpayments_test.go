package nowpayments

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetworkInfo(t *testing.T) {
	cfg := Config{
		IPNSecret: "wQRl2P7/nvgjlqoRhcARpJbFQR6/hZ92",
	}

	nowPayments, err := NewNowPayments(&cfg)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	data := `{"actually_paid":0,"actually_paid_at_fiat":0,"fee":{"currency":"usdtbsc","depositFee":0,"serviceFee":0,"withdrawalFee":0},"invoice_id":4978049764,"order_description":null,"order_id":"181563","outcome_amount":46.5258178,"outcome_currency":"usdtbsc","parent_payment_id":null,"pay_address":"35AW2C7VeU6z6dxG1rrWEDR3qHBarkHGYw","pay_amount":0.00096324,"pay_currency":"btc","payin_extra_id":null,"payment_extra_ids":null,"payment_id":5613681154,"payment_status":"finished","price_amount":50,"price_currency":"usd","purchase_id":"4621639450","updated_at":1708011564521}`
	jsonRaeder := strings.NewReader(data)
	request, err := http.NewRequest("POST", "something-url", jsonRaeder)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("x-nowpayments-sig", "4f0d652781345b093c82ad1bf2812ab4cfcde9cf670d5d11147e908d4ea99d9f0ac8425346c4f917be97c441546f8b102be86dd4eb57127d8642ae162851627d")
	nowPayments.webhookFunc(w, request)

	assert.Equal(t, w.Code, 200)
}
