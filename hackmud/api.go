package hackmud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Client for API operations
type Client struct {
	apiToken   string
	baseURL    string
	HTTPClient *http.Client
	Debug      bool
}

// NewClient creates new Hackmud Chat API client with given API key
func NewClient(apiToken string) *Client {
	return &Client{
		apiToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL: "https://www.hackmud.com",
		Debug:   false,
	}
}

type errorResponse struct {
	Code    int    `json:"status"`
	Message string `json:"error"`
}

type successResponse struct {
	OK bool `json:"ok"`
	// Data interface{} `json:"data"`
}

// Content-type and body should be already added to req
func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	// req.Header.Set("chat_token", c.apiKey)
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		if c.Debug {
			log.Println("HTTP call failed")
		}
		return err
	}

	defer res.Body.Close()

	// Try to unmarshall into errorResponse
	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			log.Println("Early response decoding failed:", errRes)
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("Unknown error with status code: %d", res.StatusCode)
	}

	// TODO: Would be great if we could somehow validate that the
	//       incoming interface embeds "successResponse" or "errorResponse"..
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		if c.Debug {
			log.Println("Late response decoding failed")
		}
		return err
	}

	// TODO: Somehow validate that the decoded response has "ok: true" set (also look for "ok: false")

	// Pretty print the JSON response
	// if c.Debug {
	// 	if json, err := InterfaceToJSON(v); err == nil {
	// 		log.Println("Response:", json)
	// 	}
	// }

	return nil
}

// InterfaceToJSON converts any interface or struct to a prettyfied JSON string
func InterfaceToJSON(v interface{}) (string, error) {
	res, err := json.MarshalIndent(v, "", "  ")
	return string(res), err
}
