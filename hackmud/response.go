package hackmud

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// UserName is just a string wrapper
type UserName string

// ChannelName is just a string wrapper
type ChannelName string

// User for the AccountData
type User map[UserName]Channel

// Channel for the AccountData
type Channel map[ChannelName][]UserName

// AccountData for the account
type AccountData struct {
	// OK    bool `json:"ok"`
	successResponse
	Users User `json:"users"`
}

// GetAccountData returns the account data for the current account
func (c *Client) GetAccountData() (*AccountData, error) {
	// TODO: Would be great if we could somehow always start/embed the body with the token
	reqBody := []byte(fmt.Sprintf(`{"chat_token":"%s"}`, c.apiToken))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/mobile/account_data.json", c.baseURL), bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("GetAccountData failed to create request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res := AccountData{}
	if err := c.sendRequest(req, &res); err != nil {
		log.Println("GetAccountData failed to send request:", err)
		return nil, err
	}

	// Pretty print the JSON response
	// if json, err := InterfaceToJSON(res); err == nil {
	// 	log.Println(json)
	// }

	// log.Println("Returning GetAccountData response:", res)
	return &res, nil
}

// TODO: Can we somehow do this globally for everything?
// Override default stringer for AcccountData
func (a *AccountData) String() string {
	return a.ToJSON()
}

// ToJSON attempts to return a prettyfied JSON string
func (a *AccountData) ToJSON() string {
	json, err := InterfaceToJSON(a)
	if err == nil {
		return json
	}
	log.Fatal(err)
	return ""
}
