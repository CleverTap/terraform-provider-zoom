package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Email          string        `json:"email,omitempty"`
	FirstName      string        `json:"first_name,omitempty"`
	LastName       string        `json:"last_name,omitempty"`
	Status         string        `json:"status,omitempty"`
	Type           int           `json:"type,omitempty"`
	Pmi            int           `json:"pmi,omitempty"`
	UsePmi         *bool         `json:"use_pmi,omitempty"`
	Timezone       string        `json:"timezone,omitempty"`
	Language       string        `json:"language,omitempty"`
	VanityName     string        `json:"vanity_name,omitempty"`
	HostKey        string        `json:"host_key,omitempty"`
	CmsUserId      string        `json:"cms_user_id,omitempty"`
	Company        string        `json:"company,omitempty"`
	GroupId        string        `json:"group_id,omitempty"`
	Manager        string        `json:"manager,omitempty"`
	Pronouns       string        `json:"pronouns,omitempty"`
	PhoneNumbers   []PhoneNumber `json:"phone_numbers,omitempty"`
	PronounsOption int           `json:"pronouns_option,omitempty"`
	RoleName       string        `json:"role_name,omitempty"`
	Department     string        `json:"dept,omitempty"`
	JobTitle       string        `json:"job_title,omitempty"`
	Location       string        `json:"location,omitempty"`
}

type PhoneNumber struct {
	Country string `json:"country,omitempty"`
	Code    string `json:"code,omitempty"`
	Number  string `json:"number,omitempty"`
	Label   string `json:"label,omitempty"`
}

type Client struct {
	authToken      string
	TimeoutMinutes int
	httpClient     *http.Client
}

func NewClient(token string, timeoutMinutes int) *Client {
	return &Client{
		authToken:      token,
		TimeoutMinutes: timeoutMinutes,
		httpClient:     &http.Client{},
	}
}

func (c *Client) NewUser(user *User) error {
	userInfo, err := json.Marshal(user)
	if err != nil {
		return err
	}
	body := strings.NewReader(fmt.Sprintf("{\"action\":\"create\",\"user_info\":" + string(userInfo) + "}"))
	_, err = c.httpRequest("POST", body, "")
	if err != nil {
		log.Println("[CREATE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) GetUser(email string) (*User, error) {
	body, err := c.httpRequest("GET", &strings.Reader{}, fmt.Sprintf("/%v", email))
	if err != nil {
		log.Println("[READ ERROR]: ", err)
		return nil, err
	}
	user := &User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Println("[READ ERROR]: ", err)
		return nil, err
	}
	return user, nil
}

func (c *Client) UpdateUser(email string, user *User) error {
	userInfo, err := json.Marshal(user)
	if err != nil {
		return err
	}
	body := strings.NewReader(string(userInfo))
	_, err = c.httpRequest("PATCH", body, fmt.Sprintf("/%v", email))
	if err != nil {
		log.Println("[UPDATE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) DeleteUser(email, status string) error {
	var err error
	if status == "pending" {
		_, err = c.httpRequest("DELETE", &strings.Reader{}, fmt.Sprintf("/%s", email))
	} else {
		_, err = c.httpRequest("DELETE", &strings.Reader{}, fmt.Sprintf("/%s?action=delete", email))
	}
	if err != nil {
		log.Println("[DELETE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) ChangeUserStatus(email, action string) error {
	action = fmt.Sprintf("{\"action\":\"%s\"}", action)
	body := strings.NewReader(action)
	_, err := c.httpRequest("PUT", body, fmt.Sprintf("/%s/status", email))
	if err != nil {
		log.Println("[DEACTIVATE/ACTIVATE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) ChangeEmail(oldEmail, newEmail string) error {
	body := strings.NewReader(fmt.Sprintf("{\"email\":\"" + newEmail + "\"}"))
	_, err := c.httpRequest("PUT", body, fmt.Sprintf("/%v/email", oldEmail))
	if err != nil {
		log.Println("[UPDATE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) httpRequest(method string, body *strings.Reader, path string) ([]byte, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("https://api.zoom.us/v2/users"+path), body)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return nil, err
	}
	authtoken := "Bearer " + c.authToken
	req.Header.Add("Authorization", authtoken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return respBody, nil
	}
	return nil, fmt.Errorf(string(respBody) + fmt.Sprintf(", StatusCode: %v", resp.StatusCode))
}

func (c *Client) IsRetry(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "429") == true {
			return true
		}
	}
	return false
}
