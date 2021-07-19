package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"log"
)

type User struct {
	Id        string `json:"id"`
	Email   string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    string  `json:"status"`
	Type      int `json:"type"`
	Pmi       int `json:"pmi"`
	RoleName    string `json:"role_name"`
	Department    string `json:"dept"`
	JobTitle    string `json:"job_title"`
	Location    string `json:"location"`
}

type NewUser struct {
	Action   string   `json:"action"`
	UserInfo UserInfo `json:"user_info"`
}

type UserInfo struct {
	Email   string `json:"email"`
	Type      int    `json:"type"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      int `json:"type"`
	Department    string `json:"dept"`
	JobTitle    string `json:"job_title"`
	Location    string `json:"location"`
}


var (
    Errors = make(map[int]string)
)

func init() {
	Errors[400] = "Bad Request, StatusCode = 400"
	Errors[404] = "User Does Not Exist , StatusCode = 404"
	Errors[409] = "User Already Exist, StatusCode = 409"
	Errors[401] = "Unautharized Access, StatusCode = 401"
	Errors[429] = "User Has Sent Too Many Request, StatusCode = 429"
}

type Client struct {
	authToken  string
	httpClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		authToken:  token,
		httpClient: &http.Client{},
	}
}

func (c *Client) NewItem(item *User) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		log.Println("[CREATE ERROR]: ", err)
		return err
	}
	_, err = c.httpRequest("POST", buf, item)
	if err != nil {
		log.Println("[CREATE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) httpRequest(method string, body bytes.Buffer, item *User) (closer io.ReadCloser, err error) {

	userjson := NewUser{
		Action: "create",
		UserInfo: UserInfo{
			Email:   item.Email,
			Type:    item.Type,
			FirstName: item.FirstName,
			LastName:  item.LastName,
		},
	}
	reqjson, _ := json.Marshal(userjson)
	payload := strings.NewReader(string(reqjson))
	req, err := http.NewRequest(method, fmt.Sprintf("%s?access_token=%s", "https://api.zoom.us/v2/users", c.authToken), payload)
	authtoken := "Bearer "+c.authToken
	req.Header.Add("Authorization", authtoken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println("[ERROR]: ",err)
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("[ERROR]: ",err)
		return nil, err
	}
	//var data Data
	var data map[string]interface{}
	newbody, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(newbody), &data)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return resp.Body, nil
    } else {
		return nil, fmt.Errorf("Error : %v",data["message"])
    }

}

func (c *Client) GetItem(name string) (*User, error) {
	body, err := c.gethttpRequest(fmt.Sprintf("%v", name), "GET", bytes.Buffer{})
	if err != nil {
		log.Println("[READ ERROR]: ", err)
		return nil, err
	}
	item := &User{}
	err = json.NewDecoder(body).Decode(item)
	if err != nil {
		log.Println("[READ ERROR]: ", err)
		return nil, err
	}
	return item, nil
}

func (c *Client) gethttpRequest(emailid, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, "https://api.zoom.us/v2/users/"+emailid, &body)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return nil, err
	}
	authtoken := "Bearer "+c.authToken
	req.Header.Add("Authorization", authtoken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR]: ",err)
		return nil, err
	}
	var data map[string]interface{}
	newbody, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(newbody), &data)
	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Error : %v",data["message"])
		}
		return nil, fmt.Errorf("Error : %v ", data["message"])
	}
	return resp.Body, nil
}

func (c *Client) UpdateItem(item *User) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(item)
	if err != nil {
		log.Println("[UPDATE ERROR]: ", err)
		return err
	}
	_, err = c.updatehttpRequest(fmt.Sprintf("%s", item.Email), "PATCH", buf,item)
	if err != nil {
		log.Println("[UPDATE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) updatehttpRequest(path,method string, body bytes.Buffer, item *User) (closer io.ReadCloser, err error) {
	updateuserjson := UpdateUser{
		FirstName: item.FirstName,
		LastName:  item.LastName,
		Type:  item.Type,
		Department: item.Department,
		JobTitle:  item.JobTitle,
		Location:  item.Location,
	}
	updatejson, _ := json.Marshal(updateuserjson)
	payload := strings.NewReader(string(updatejson))
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s","https://api.zoom.us/v2/users", path), payload)
	authtoken := "Bearer "+c.authToken
	req.Header.Add("Authorization", authtoken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println("[ERROR]: ", err)
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("[ERROR]: ", err)
		return nil, err
	}
	var data map[string]interface{}
	newbody, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(newbody), &data)
	if resp.StatusCode >= 200 && resp.StatusCode <= 400 {
		return resp.Body, nil
    } else {
		return nil, fmt.Errorf("Error : %v",data["message"])
    }
}

func (c *Client) DeleteItem(userId string) error {
	_, err := c.deletehttpRequest(fmt.Sprintf("%s", userId), "DELETE", bytes.Buffer{})
	if err != nil {
		log.Println("[DELETE ERROR]: ", err)
		return err
	}
	return nil
}

func (c *Client) deletehttpRequest(path, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method,fmt.Sprintf("%s/%s", "https://api.zoom.us/v2/users", path), &body)
	if err != nil {
		log.Println("[DELETE ERROR]: ", err)
		return nil, err
	}
	authtoken := "Bearer "+c.authToken
	req.Header.Add("Authorization", authtoken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("[DELETE ERROR]: ", err)
		return nil, err
	}
	var data map[string]interface{}
	newbody, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(newbody), &data)
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return resp.Body, nil
    } else {
		log.Println("Broken Request")
		return nil, fmt.Errorf("Error : %v",data["message"])
    }


}

func (c *Client) DeactivateUser(userId string, status string) error {
	log.Println("Changing Status of User : ", userId)
	url := fmt.Sprintf("https://api.zoom.us/v2/users/%s/status", userId)
	data := fmt.Sprintf("{\"action\":\"%s\"}", status)
	payload := strings.NewReader(data)
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		log.Println("[DEACTIVATE/ACTIVATE ERROR]: ",err)
		return nil
	}
	authtoken := "Bearer "+c.authToken
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", authtoken)
	_, err = c.httpClient.Do(req)
	if err != nil {
		log.Println("[DEACTIVATE/ACTIVATE ERROR]: ",err)
		return nil
	}
	return nil
}

