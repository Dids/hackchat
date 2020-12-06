package hackmud

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
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
	Code    int    `json:"status"` // TODO: Is this actually part of the response?
	Message string `json:"error"`  // TODO: Is this actually part of the response?
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
			log.Error("HTTP call failed")
		}
		return err
	}

	defer res.Body.Close()

	// Check for error code, then unmarshal and validate the error response
	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			if c.Debug {
				log.Error("Error response decoding failed:", errRes)
			}
			return errors.New(errRes.Message)
		}

		// FIXME: Return errRes.Code and errRes.Message, but check that they're part of the API first?
		return fmt.Errorf("Unknown error with status code: %d", res.StatusCode)
	}

	// Unmarshal and validate the success response
	var succRes successResponse
	if err = json.NewDecoder(res.Body).Decode(&succRes); err != nil {
		if c.Debug {
			log.Error("Success response decoding failed")
		}
		return err
	}
	if !succRes.OK {
		return errors.New("Success response not OK")
	}

	// Unmarshal and validate the final response
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		if c.Debug {
			log.Error("Final response decoding failed")
		}
		return err
	}

	// Pretty print the JSON response
	// if c.Debug {
	// 	if json, err := InterfaceToJSON(v); err == nil {
	// 		log.Info("Response:", json)
	// 	}
	// }

	return nil
}

// InterfaceToJSON converts any interface or struct to a prettyfied JSON string
func InterfaceToJSON(v interface{}) (string, error) {
	res, err := json.MarshalIndent(v, "", "  ")
	return string(res), err
}
