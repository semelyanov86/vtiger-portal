package vtiger

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type WebRequests struct {
	config VtigerConnectionConfig
}

type CrmFetcher interface {
	FetchBytes(ctx context.Context, postfix string) ([]byte, error)
	SendData(ctx context.Context, data RequestData) ([]byte, error)
}

var ErrWrongStatusCode = errors.New("wrong status code")

func NewWebRequest(config VtigerConnectionConfig) WebRequests {
	return WebRequests{config: config}
}

func (w WebRequests) FetchBytes(ctx context.Context, postfix string) ([]byte, error) {
	url := w.config.Url + "?" + postfix
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("User-Agent", "Go-Portal")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != http.StatusOK {
		return []byte{}, e.Wrap("status code: "+strconv.Itoa(res.StatusCode), ErrWrongStatusCode)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func (w WebRequests) SendData(ctx context.Context, data RequestData) ([]byte, error) {
	form := url.Values{
		"operation":   {data.FormParams.Operation},
		"username":    {data.FormParams.Username},
		"accessKey":   {data.FormParams.AccessKey},
		"sessionName": {data.FormParams.SessionName},
	}
	reqBody := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.config.Url, reqBody)

	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("User-Agent", "Go-Portal")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != http.StatusOK {
		return []byte{}, e.Wrap("status code: "+strconv.Itoa(res.StatusCode), ErrWrongStatusCode)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func (w WebRequests) SendObject(ctx context.Context, operation string, session string, elementType string, data map[string]any) ([]byte, error) {
	jsonObject, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	form := url.Values{
		"operation":   {operation},
		"sessionName": {session},
		"element":     {string(jsonObject)},
		"elementType": {elementType},
	}
	reqBody := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.config.Url, reqBody)

	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("User-Agent", "Go-Portal")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != http.StatusOK {
		return []byte{}, e.Wrap("status code: "+strconv.Itoa(res.StatusCode), ErrWrongStatusCode)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
