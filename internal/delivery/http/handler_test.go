package http_test

import (
	"github.com/semelyanov86/vtiger-portal/internal/config"
	http2 "github.com/semelyanov86/vtiger-portal/internal/delivery/http"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	h := http2.NewHandler(&service.Services{}, &config.Config{})

	require.IsType(t, &http2.Handler{}, h)
}

func TestNewHandler_Init(t *testing.T) {
	h := http2.NewHandler(&service.Services{}, &config.Config{
		Limiter: config.Limiter{
			Rps:   2,
			Burst: 4,
			TTL:   10 * time.Minute,
		},
	})

	router := h.Init()

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
