package vtiger

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const TokenKey = "portal_vtiger_token"

const CacheTTL = 500

const TimeoutInSec = 5

var ErrMaxRetriesConnection = errors.New("could not process with max retries")

var ErrResponseError = errors.New("response error from vtiger code 3")

var ErrNoCacheKey = errors.New("code 8 There is no key for index in cache")

var ErrWrongCredentials = errors.New("invalid credentials due to auth")

var ErrCanNotParseCountObject = errors.New("can not parse count object")

type VtigerConnector struct {
	cache      cache.Cache
	connection VtigerConnectionConfig
	fetcher    CrmFetcher
}

var sessionIdMutex = sync.Mutex{}

type VtigerConnectionConfig struct {
	Url                   string `yaml:"url"`
	Login                 string `yaml:"login"`
	Password              string `yaml:"password"`
	PersistenceConnection bool   `yaml:"persistenceConnection"`
	MaxRetries            int    `yaml:"maxRetries"`
}

type SessionData struct {
	SessionID     string `json:"sessionName"`
	ExpireTime    int64  `json:"expireTime"`
	ServerTime    int64  `json:"serverTime"`
	Token         string `json:"token"`
	UserId        string `json:"userId"`
	VtigerVersion string `json:"vtigerVersion"`
}

type ResultData interface {
	SessionData | Module | map[string]any | []map[string]any | []File
}

type ErrorData struct {
	Message string
	Code    string
}

type VtigerResponse[T ResultData] struct {
	Success bool      `json:"success"`
	Result  T         `json:"result"`
	Error   ErrorData `json:"error"`
}

type RequestData struct {
	FormParams FormParamsData `json:"form_params"`
}

type FormParamsData struct {
	Operation   string `json:"operation"`
	Username    string `json:"username"`
	AccessKey   string `json:"access_key"`
	SessionName string `json:"sessionName"`
}

type File struct {
	Fileid       string `json:"fileid"`
	Filename     string `json:"filename"`
	Filetype     string `json:"filetype"`
	Filesize     int    `json:"filesize"`
	Filecontents string `json:"filecontents"`
}

func NewVtigerConnector(cache cache.Cache, config VtigerConnectionConfig, fetcher CrmFetcher) VtigerConnector {
	return VtigerConnector{
		cache:      cache,
		connection: config,
		fetcher:    fetcher,
	}
}

func (c VtigerConnector) Lookup(ctx context.Context, dataType, value, module string, columns []string) (*VtigerResponse[[]map[string]any], error) {
	sessionID, err := c.sessionId()
	if err != nil {
		return nil, err
	}

	// Update columns into the proper format
	var columnsText string
	for _, column := range columns {
		columnsText += "\"" + column + "\","
	}
	columnsText = strings.TrimSuffix(columnsText, ",")

	resp, err := c.fetcher.FetchBytes(ctx, url.Values{
		"operation":   {"lookup"},
		"sessionName": {sessionID},
		"type":        {dataType},
		"value":       {value},
		"searchIn":    {"{\"" + module + "\":[" + columnsText + "]}"},
	}.Encode())

	if err != nil {
		return nil, e.Wrap("code 7", err)
	}

	err = c.close(sessionID)
	if err != nil {
		return nil, err
	}
	vtigerResponse := &VtigerResponse[[]map[string]any]{}

	return processVtigerResponse[[]map[string]any](resp, vtigerResponse)
}

func (c VtigerConnector) AddRelated(ctx context.Context, source string, related string, label string) (*VtigerResponse[[]map[string]any], error) {
	sessionID, err := c.sessionId()
	if err != nil {
		return nil, err
	}

	resp, err := c.fetcher.FetchBytes(ctx, url.Values{
		"operation":       {"add_related"},
		"sessionName":     {sessionID},
		"sourceRecordId":  {source},
		"relatedRecordId": {related},
		"relationIdLabel": {label},
	}.Encode())

	if err != nil {
		return nil, e.Wrap("code 7", err)
	}

	err = c.close(sessionID)
	if err != nil {
		return nil, err
	}
	vtigerResponse := &VtigerResponse[[]map[string]any]{}

	return processVtigerResponse[[]map[string]any](resp, vtigerResponse)
}

func (c VtigerConnector) Query(ctx context.Context, query string) (*VtigerResponse[[]map[string]any], error) {
	sessionID, err := c.sessionId()
	if err != nil {
		return nil, err
	}

	resp, err := c.fetcher.FetchBytes(ctx, url.Values{
		"operation":   {"query"},
		"sessionName": {sessionID},
	}.Encode()+"&query="+url.QueryEscape(query))

	if err != nil {
		return nil, e.Wrap("code 7", err)
	}

	err = c.close(sessionID)
	if err != nil {
		return nil, err
	}
	vtigerResponse := &VtigerResponse[[]map[string]any]{}

	return processVtigerResponse[[]map[string]any](resp, vtigerResponse)
}

func (c VtigerConnector) RetrieveRelated(ctx context.Context, id string, module string) (*VtigerResponse[[]map[string]any], error) {
	sessionID, err := c.sessionId()
	if err != nil {
		return nil, err
	}

	resp, err := c.fetcher.FetchBytes(ctx, url.Values{
		"operation":    {"retrieve_related"},
		"sessionName":  {sessionID},
		"id":           {id},
		"relatedLabel": {module},
		"relatedType":  {module},
	}.Encode())

	if err != nil {
		return nil, e.Wrap("code 7", err)
	}

	err = c.close(sessionID)
	if err != nil {
		return nil, err
	}
	vtigerResponse := &VtigerResponse[[]map[string]any]{}

	return processVtigerResponse[[]map[string]any](resp, vtigerResponse)
}

func (c VtigerConnector) Retrieve(ctx context.Context, id string) (*VtigerResponse[map[string]any], error) {
	sessionID, err := c.sessionId()
	if err != nil {
		return nil, err
	}

	webRequest := NewWebRequest(c.connection)

	// send a request to retrieve a record
	resp, err := webRequest.FetchBytes(ctx, url.Values{
		"operation":   {"retrieve"},
		"sessionName": {sessionID},
		"id":          {id},
	}.Encode())
	if err != nil {
		return nil, e.Wrap("code 7", err)
	}

	err = c.close(sessionID)
	if err != nil {
		return nil, err
	}
	vtigerResponse := &VtigerResponse[map[string]any]{}

	return processVtigerResponse[map[string]any](resp, vtigerResponse)
}

func (c VtigerConnector) Describe(ctx context.Context, element string) (*VtigerResponse[Module], error) {
	sessionID, err := c.sessionId()
	if err != nil {
		return nil, err
	}

	webRequest := NewWebRequest(c.connection)

	// send a request to retrieve a record
	resp, err := webRequest.FetchBytes(ctx, url.Values{
		"operation":   {"describe"},
		"sessionName": {sessionID},
		"elementType": {element},
	}.Encode())
	if err != nil {
		return nil, e.Wrap("code 7", err)
	}

	err = c.close(sessionID)
	if err != nil {
		return nil, err
	}
	vtigerResponse := &VtigerResponse[Module]{}

	return processVtigerResponse[Module](resp, vtigerResponse)
}

func (c VtigerConnector) Delete(ctx context.Context, element string) error {
	sessionID, err := c.sessionId()
	if err != nil {
		return err
	}

	webRequest := NewWebRequest(c.connection)

	// send a request to retrieve a record
	resp, err := webRequest.FetchBytes(ctx, url.Values{
		"operation":   {"delete"},
		"sessionName": {sessionID},
		"id":          {element},
	}.Encode())
	if err != nil {
		return e.Wrap("code 7", err)
	}
	m := VtigerResponse[Module]{}
	err = json.Unmarshal(resp, &m)

	if err != nil {
		return e.Wrap("code 12", err)
	}

	if !m.Success {
		return e.Wrap("code 15: "+m.Error.Message, ErrResponseError)
	}

	err = c.close(sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (c VtigerConnector) RetrieveFiles(ctx context.Context, id string) (*VtigerResponse[[]File], error) {
	sessionID, err := c.sessionId()
	if err != nil {
		return nil, err
	}

	webRequest := NewWebRequest(c.connection)

	// send a request to retrieve a record
	resp, err := webRequest.FetchBytes(ctx, url.Values{
		"operation":   {"files_retrieve"},
		"sessionName": {sessionID},
		"id":          {id},
	}.Encode())
	if err != nil {
		return nil, e.Wrap("code 7", err)
	}

	err = c.close(sessionID)
	if err != nil {
		return nil, err
	}
	vtigerResponse := &VtigerResponse[[]File]{}

	return processVtigerResponse[[]File](resp, vtigerResponse)
}

func (c VtigerConnector) Update(ctx context.Context, data map[string]any) (*VtigerResponse[map[string]any], error) {
	return c.doUpdate(ctx, data, "update")
}

func (c VtigerConnector) Revise(ctx context.Context, data map[string]any) (*VtigerResponse[map[string]any], error) {
	return c.doUpdate(ctx, data, "revise")
}

func (c VtigerConnector) doUpdate(ctx context.Context, data map[string]any, operation string) (*VtigerResponse[map[string]any], error) {
	sessionID, err := c.sessionId()
	if err != nil {
		return nil, err
	}

	webRequest := NewWebRequest(c.connection)
	// send a request to retrieve a record
	resp, err := webRequest.SendObject(ctx, operation, sessionID, "", data)
	if err != nil {
		return nil, e.Wrap("code 7", err)
	}

	err = c.close(sessionID)
	if err != nil {
		return nil, err
	}
	vtigerResponse := &VtigerResponse[map[string]any]{}

	return processVtigerResponse[map[string]any](resp, vtigerResponse)
}

func (c VtigerConnector) Create(ctx context.Context, element string, data map[string]any) (*VtigerResponse[map[string]any], error) {
	sessionID, err := c.sessionId()
	if err != nil {
		return nil, err
	}

	webRequest := NewWebRequest(c.connection)

	// send a request to retrieve a record
	resp, err := webRequest.SendObject(ctx, "create", sessionID, element, data)
	if err != nil {
		return nil, e.Wrap("code 7", err)
	}

	err = c.close(sessionID)
	if err != nil {
		return nil, err
	}
	vtigerResponse := &VtigerResponse[map[string]any]{}

	return processVtigerResponse[map[string]any](resp, vtigerResponse)
}

func (c VtigerConnector) Count(ctx context.Context, module string, filters map[string]string) (int, error) {
	query := "SELECT COUNT(*) FROM " + module + " WHERE "
	for field, value := range filters {
		field = strings.TrimPrefix(field, "_")
		query += field + " = '" + value + "' OR "
	}
	query = strings.TrimSuffix(query, " OR ")
	query += ";"
	return c.ExecuteCount(ctx, query)
}

func (c VtigerConnector) ExecuteCount(ctx context.Context, query string) (int, error) {
	result, err := c.Query(ctx, query)
	if err != nil {
		return 0, e.Wrap("can not execute query "+query+", got error", err)
	}
	countObject := result.Result[0]
	if countObject == nil {
		return 0, ErrCanNotParseCountObject
	}
	count, err := strconv.Atoi(countObject["count"].(string))
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c VtigerConnector) getToken() (SessionData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*TimeoutInSec)
	defer cancel()
	tryCounter := 1

	var result *VtigerResponse[SessionData]

	for {
		response, err := c.fetcher.FetchBytes(ctx, "operation=getchallenge&username="+c.connection.Login)
		if err != nil {
			return SessionData{}, e.Wrap("error code 7 from vtiger connector", err)
		}
		tryCounter++
		result, err = processVtigerResponse(response, result)

		if err != nil {
			return SessionData{}, err
		}

		if tryCounter > c.connection.MaxRetries || result.Success {
			break
		}
	}

	if tryCounter >= c.connection.MaxRetries {
		return SessionData{}, e.Wrap("code 6:", ErrMaxRetriesConnection)
	}

	return result.Result, nil
}

func (c VtigerConnector) storeSession() (SessionData, error) {
	token, err := c.getToken()
	if err != nil {
		return SessionData{}, err
	}
	cachedValue, err := json.Marshal(token)
	if err != nil {
		return SessionData{}, err
	}
	c.cache.Set(TokenKey, cachedValue, CacheTTL)
	return token, nil
}

func processVtigerResponse[T ResultData](response []byte, data *VtigerResponse[T]) (*VtigerResponse[T], error) {
	if err := json.Unmarshal(response, &data); err != nil {
		return data, err
	}
	if !data.Success {
		return data, e.Wrap(data.Error.Code+" - "+data.Error.Message, ErrResponseError)
	}
	return data, nil
}

// sessionId - Get the session id for a login either from a stored session id or fresh from the API
func (c VtigerConnector) sessionId() (string, error) {
	// Get the sessionData from the cache
	var sessionData SessionData
	var decodedSessionData *SessionData
	sessionIdMutex.Lock()
	defer sessionIdMutex.Unlock()
	cachedSessionData, err := c.cache.Get(TokenKey)

	if errors.Is(cache.ErrItemNotFound, err) {
		sessionData, err = c.storeSession()
		if err != nil {
			return "", e.Wrap("can not store a session", err)
		}
		loginResult, err := c.login(sessionData)
		if err != nil {
			return "", e.Wrap("can not login", err)
		}
		return loginResult.SessionID, nil
	}

	if cachedSessionData != nil {
		decodedSessionData = &SessionData{}
		err = json.Unmarshal(cachedSessionData, decodedSessionData)
		if err != nil {
			return "", e.Wrap("can not convert caches data to session", err)
		}

		if decodedSessionData.ExpireTime > time.Now().Unix() || decodedSessionData.Token == "" {
			sessionData, err = c.storeSession()
			if err != nil {
				return "", e.Wrap("can not receive session data", err)
			}
		}
	} else {
		sessionData, err = c.storeSession()
		if err != nil {
			return "", e.Wrap("can not receive session data", err)
		}
	}

	if sessionData.SessionID == "" {
		loginResult, err := c.login(sessionData)
		if err != nil {
			return "", e.Wrap("can not login", err)
		}
		return loginResult.SessionID, nil
	} else {
		return sessionData.SessionID, nil
	}
}

func (c VtigerConnector) login(session SessionData) (SessionData, error) {
	var sessionData SessionData
	var responseData *VtigerResponse[SessionData]
	var token = session.Token
	generatedKey := fmt.Sprintf("%x", md5.Sum([]byte(token+c.connection.Password)))
	var tryCounter = 1

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*TimeoutInSec)
	defer cancel()
	requestData := RequestData{FormParams: FormParamsData{
		Operation: "login",
		Username:  c.connection.Login,
		AccessKey: generatedKey,
	}}

	for {
		// login using username and accesskey
		resp, err := c.fetcher.SendData(ctx, requestData)

		if err != nil {
			return sessionData, e.Wrap("code 7", err)
		}
		// decode the response
		loginResult, err := processVtigerResponse[SessionData](resp, responseData)
		loginResult.Result.Token = session.Token
		loginResult.Result.ExpireTime = session.ExpireTime
		loginResult.Result.ServerTime = session.ServerTime
		responseData = loginResult
		if err != nil {
			c.cache.Set(TokenKey, []byte{}, 0)
			return sessionData, e.Wrap("wrong response received from vtiger during login", err)
		}
		tryCounter++
		if loginResult.Success || tryCounter > c.connection.MaxRetries {
			break
		}
	}
	if tryCounter >= c.connection.MaxRetries {
		return sessionData, e.Wrap(fmt.Sprintf("Could not complete login request within %d tries", c.connection.MaxRetries), ErrMaxRetriesConnection)
	}

	if responseData.Success {
		responseSessionData := responseData.Result
		encodedSession, err := json.Marshal(responseSessionData)
		if err != nil {
			return sessionData, err
		}

		err = c.cache.Set(TokenKey, encodedSession, CacheTTL)
		if err != nil {
			return sessionData, e.Wrap("can not store cache session", err)
		}
		return responseSessionData, nil
	}

	if responseData.Error.Code == "INVALID_USER_CREDENTIALS" || responseData.Error.Code == "INVALID_SESSIONID" {
		c.cache.Set(TokenKey, []byte{}, 0)
		return sessionData, ErrWrongCredentials
	}

	return sessionData, e.Wrap(responseData.Error.Message, ErrResponseError)
}

func (c VtigerConnector) close(sessionID string) error {
	if c.connection.PersistenceConnection {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*TimeoutInSec)
	defer cancel()

	requestData := RequestData{FormParams: FormParamsData{
		Operation:   "logout",
		SessionName: sessionID,
	}}

	_, err := c.fetcher.SendData(ctx, requestData)
	return err
}
