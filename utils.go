package GoSDK

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clearblade/mqttclient"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	addr string
)

//Client is a convience interface for API consumers, if they want to use the same functions for both
//Dev Users and unprivleged users, such as tiny helper functions. please use this with care
type Client interface {
	InsertData(string, interface{}) error
	UpdateData(string, [][]map[string]interface{}, map[string]interface{}) error
	GetData(string, [][]map[string]interface{}) (map[string]interface{}, error)
	DeleteData(string, [][]map[string]interface{}) (map[string]interface{}, error)
	Authenticate(string, string) error
	Register(username, password string) error
	Logout() error
}

//cbClient will supply various information that differs between privleged and unprivleged users
type cbClient interface {
	credentials() [][]string    //the inner slice is a tuple of "Header":"Value"
	authInfo() (string, string) //username,password
	preamble(bool) string       // "api/v/1" ||"admin", the bool is whether or not it's a user call
	setToken(string)
	getToken(string)
	getKeySecret() (string, string)
}

type client struct {
	mrand *rand.Rand
	MQTTClient
	SystemKey    string
	SystemSecret string
}

type UserClient struct {
	UserToken string
	mrand     *rand.Rand
	MQTTClient
	SystemKey    string
	SystemSecret string
}

type DevClient struct {
	DevToken string
	mrand    *rand.Rand
	MQTTClient
	SystemKey    string
	SystemSecret string
}

func authenticate(c cbClient, username, password string) error {
	resp, err := post(c.preamble(true)+"/auth", map[string]interface{}{
		"username": username,
		"password": password,
	}, c.creds())
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error in authenticating %v\n", resp.Body)
	}

	var token string = ""
	switch c.(type) {
	case *UserClient:
		token = resp.Body.(map[string]interface{})["user_token"].(string)
	case *DevClient:
		token = resp.Body.(map[string]interface{})["dev_token"].(string)
	}
	if token == "" {
		return fmt.Errorf("Token not present i response from platform %+v", resp.Body)
	}
	c.setToken(token)
	return nil
}

func register(c cbClient, username, password string) error {
	resp, err := post(c.preamble(true)+"/reg", map[string]interface{}{
		"username": username,
		"password": password,
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error in authenticating %v\n", resp.Body)
	}

	var token string = ""
	switch c.(type) {
	case *UserClient:
		token = resp.Body.(map[string]interface{})["user_token"].(string)
	case *DevClient:
		token = resp.Body.(map[string]interface{})["dev_token"].(string)
	}
	if token == "" {
		return fmt.Errorf("Token not present i response from platform %+v", resp.Body)
	}
	//there isn't really a decent response to this one
	return nil
}

func logout(c cbClient) error {
	resp, err := post(c.preamble(true)+"/logout", nil, c.creds())
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error in authenticating %v\n", resp.Body)
	}
	return nil
}

func (u *UserClient) Authenticate(username, password string) error {
	return authenticate(u, username, password)
}

func (d *DevClient) Authenticate(username, password string) error {
	return authenticate(d, username, password)
}

func (u *UserClient) Register(username, password string) error {
	return register(u, username, password)
}

func (d *DevClient) Register(username, password string) error {
	return register(d, username, password)
}

func (u *UserClient) Logout() error {
	return logout(u)
}

func (d *DevClient) Logout() error {
	return logout(d)
}

type Client struct {
	URL        string
	Headers    map[string]string
	MQTTClient *mqttclient.Client
	//We are redundantly storing these so that they presist after the usertoken
	//is added to the header, and the SystemKey and SystemSecret are removed
	//as the MQTT Client requires them
	SystemKey    string
	SystemSecret string
	mrand        *rand.Rand
}

type CbReq struct {
	Body        interface{}
	Method      string
	Endpoint    string
	QueryString string
}

type CbResp struct {
	Body       interface{}
	StatusCode int
}

func NewClient() *Client {
	return &Client{
		URL:        "https://platform.clearblade.com",
		Headers:    map[string]string{},
		MQTTClient: nil,
		mrand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *Client) AddHeader(key, value string) {
	c.Headers[key] = value
}

func (c *Client) RemoveHeader(key string) {
	delete(c.Headers, key)
}

func (c *Client) GetHeader(key string) string {
	s, _ := c.Headers[key]
	return s
}

func (c *Client) SetSystem(key, secret string) {
	c.SystemKey = key
	c.SystemSecret = secret
	c.AddHeader("ClearBlade-SystemKey", key)
	c.AddHeader("ClearBlade-SystemSecret", secret)
}

func (c *Client) SetDevToken(tok string) {
	c.RemoveHeader("ClearBlade-SystemKey")
	c.RemoveHeader("ClearBlade-SystemSecret")
	c.RemoveHeader("ClearBlade-UserToken") // just in case
	c.AddHeader("ClearBlade-DevToken", tok)
}

func (c *Client) SetUserToken(tok string) {
	c.RemoveHeader("ClearBlade-SystemKey")
	c.RemoveHeader("ClearBlade-SystemSecret")
	c.RemoveHeader("ClearBlade-DevToken") // just in case
	c.AddHeader("ClearBlade-UserToken", tok)
}

func (c *Client) GetSystemInfo() (string, string) {
	k := c.GetHeader("ClearBlade-SystemKey")
	s := c.GetHeader("ClearBlade-SystemSecret")
	return k, s
}

func (c *Client) GetUserToken() string {
	tok := c.GetHeader("ClearBlade-UserToken")
	return tok
}

func (c *Client) GetDevToken() string {
	tok := c.GetHeader("ClearBlade-DevToken")
	return tok
}

func do(r *CbReq, creds [][]string) (*CbResp, error) {
	var bodyToSend *bytes.Buffer
	if r.Body != nil {
		b, jsonErr := json.Marshal(r.Body)
		if jsonErr != nil {
			return nil, fmt.Errorf("JSON Encoding Error: %v", jsonErr)
		}
		bodyToSend = bytes.NewBuffer(b)
	} else {
		bodyToSend = nil
	}
	url := c.URL + r.Endpoint
	if r.QueryString != "" {
		url += "?" + r.QueryString
	}
	var req *http.Request
	var reqErr error
	if bodyToSend != nil {
		req, reqErr = http.NewRequest(r.Method, url, bodyToSend)
	} else {
		req, reqErr = http.NewRequest(r.Method, url, nil)
	}
	if reqErr != nil {
		return nil, fmt.Errorf("Request Creation Error: %v", reqErr)
	}
	for _, c := range creds {
		if len(c) != 2 {
			return nil, fmt.Errorf("Request Creation Error: Invalid credential header supplied")
		}
		req.Header.Add(c[0], c[1])
	}

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error Making Request: %v", err)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("Error Reading Response Body: %v", readErr)
	}
	var d interface{}
	if len(body) == 0 {
		return &CbResp{
			Body:       nil,
			StatusCode: resp.StatusCode,
		}, nil
	}
	buf := bytes.NewBuffer(body)
	dec := json.NewDecoder(buf)
	decErr := dec.Decode(&d)
	var bod interface{}
	if decErr != nil {
		//		return nil, fmt.Errorf("JSON Decoding Error: %v\n With Body: %v\n", decErr, string(body))
		bod = string(body)
	}
	switch d.(type) {
	case []interface{}:
		bod = d
	case map[string]interface{}:
		bod = d
	default:
		bod = string(body)
	}
	return &CbResp{
		Body:       bod,
		StatusCode: resp.StatusCode,
	}, nil
}

func get(endpoint string, query map[string]string, creds [][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        nil,
		Method:      "GET",
		Endpoint:    endpoint,
		QueryString: query_to_string(query),
	}
	return do(req, creds)
}

func post(endpoint string, body interface{}, creds [][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        body,
		Method:      "POST",
		Endpoint:    endpoint,
		QueryString: "",
	}
	return do(req, creds)
}

func put(endpoint string, body interface{}, heads [][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        body,
		Method:      "PUT",
		Endpoint:    endpoint,
		QueryString: "",
	}
	return do(req, heads)
}

func delete(endpoint string, query map[string]string, heds [][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        nil,
		Method:      "DELETE",
		Endpoint:    endpoint,
		QueryString: query_to_string(query),
	}
	return do(req, heds)
}

func query_to_string(query map[string]string) string {
	qryStr := ""
	for k, v := range query {
		qryStr += k + "=" + v + "&"
	}
	return strings.TrimSuffix(qryStr, "&")
}
