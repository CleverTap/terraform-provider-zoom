package client

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func init() {
	os.Setenv("ZOOM_TOKEN", "DEMOVALUE")
}

func TestClient_GetItem(t *testing.T) {
	testCases := []struct {
		testName     string
		itemName     string
		seedData     map[string]User
		expectErr    bool
		expectedResp *User
	}{
		{
			testName: "user exists",
			itemName: "user@gmail.com",
			seedData: map[string]User{
				"user@gmail.com": {
					Id:         "oJ8qBrheQ4KJ6qozaa4QhA",
					Email:      "user@gmail.com",
					FirstName:  "FirstName",
					LastName:   "LastName",
					Type:       1,
					Pmi:        6730446034,
					RoleName:   "Member",
					Status:     "active",
					Department: "",
					JobTitle:   "",
					Location:   "",
				},
			},
			expectErr: false,
			expectedResp: &User{
				Id:         "oJ8qBrheQ4KJ6qozaa4QhA",
				Email:      "user@gmail.com",
				FirstName:  "FirstName",
				LastName:   "LastName",
				Type:       1,
				Pmi:        6730446034,
				RoleName:   "Member",
				Status:     "active",
				Department: "",
				JobTitle:   "",
				Location:   "",
			},
		},

		{
			testName:     "user does not exist",
			itemName:     "user@gmail.com",
			seedData:     nil,
			expectErr:    true,
			expectedResp: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient(os.Getenv("ZOOM_TOKEN"))
			item, err := client.GetItem(tc.itemName)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResp, item)
		})
	}
}

func TestClient_NewItem(t *testing.T) {
	testCases := []struct {
		testName  string
		newItem   *User
		seedData  map[string]User
		expectErr bool
	}{
		{
			testName: "user creation successful",
			newItem: &User{
				Id:         "oJ8qBrheQ4KJ6qozaa4QhA",
				Email:      "user@gmail.com",
				FirstName:  "FirstName",
				LastName:   "LastName",
				Type:       1,
				Pmi:        6730446034,
				RoleName:   "Member",
				Status:     "active",
				Department: "",
				JobTitle:   "",
				Location:   "",
			},
			seedData:  nil,
			expectErr: false,
		},
		{
			testName: "item already exists",
			newItem: &User{
				Id:         "oJ8qBrheQ4KJ6qozaa4QhA",
				Email:      "user@gmail.com",
				FirstName:  "FirstName",
				LastName:   "LastName",
				Type:       1,
				Pmi:        6730446034,
				RoleName:   "Member",
				Status:     "active",
				Department: "",
				JobTitle:   "",
				Location:   "",
			},
			seedData: map[string]User{
				"item1": {
					Id:         "oJ8qBrheQ4KJ6qozaa4QhA",
					Email:      "user@gmail.com",
					FirstName:  "FirstName",
					LastName:   "LastName",
					Type:       1,
					Pmi:        6730446034,
					RoleName:   "Member",
					Status:     "active",
					Department: "",
					JobTitle:   "",
					Location:   "",
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient(os.Getenv("ZOOM_TOKEN"))
			err := client.NewItem(tc.newItem)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			item, err := client.GetItem(tc.newItem.Email)
			assert.NoError(t, err)
			assert.Equal(t, tc.newItem, item)
		})
	}
}

func TestClient_UpdateItem(t *testing.T) {
	testCases := []struct {
		testName    string
		updatedItem *User
		seedData    map[string]User
		expectErr   bool
	}{
		{
			testName: "item exists",
			updatedItem: &User{
				Id:         "oJ8qBrheQ4KJ6qozaa4QhA",
				Email:      "user@gmail.com",
				FirstName:  "FirstName",
				LastName:   "LastName",
				Type:       1,
				Pmi:        6730446034,
				RoleName:   "Member",
				Status:     "active",
				Department: "",
				JobTitle:   "",
				Location:   "",
			},
			seedData: map[string]User{
				"item1": {
					Id:         "oJ8qBrheQ4KJ6qozaa4QhA",
					Email:      "user@gmail.com",
					FirstName:  "FirstName",
					LastName:   "LastName",
					Type:       1,
					Pmi:        6730446034,
					RoleName:   "Member",
					Status:     "active",
					Department: "",
					JobTitle:   "",
					Location:   "",
				},
			},
			expectErr: false,
		},
		{
			testName: "item does not exist",
			updatedItem: &User{
				Id:         "dfhjjddfjsd",
				Email:      "user@gmail.com",
				FirstName:  "FirstName",
				LastName:   "LastName",
				Type:       1,
				RoleName:   "Member",
				Status:     "active",
				Department: "devops",
				JobTitle:   "Engineer",
				Location:   "Delhi",
			},
			seedData:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient(os.Getenv("ZOOM_TOKEN"))
			err := client.UpdateItem(tc.updatedItem)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			item, err := client.GetItem(tc.updatedItem.Email)
			assert.NoError(t, err)
			assert.Equal(t, tc.updatedItem, item)
		})
	}
}

func TestClient_DeleteItem(t *testing.T) {
	testCases := []struct {
		testName  string
		itemName  string
		seedData  map[string]User
		expectErr bool
	}{
		{
			testName: "user exists",
			itemName: "user@gmail.com",
			seedData: map[string]User{
				"user1": {
					Id:         "t2OUx6lvTMedrAiW2ffURA",
					Email:      "user@gmail.com",
					FirstName:  "FirstName",
					LastName:   "LastName",
					Type:       1,
					RoleName:   "Member",
					Status:     "active",
					Department: "devops",
					JobTitle:   "Engineer",
					Location:   "Delhi",
				},
			},
			expectErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient(os.Getenv("ZOOM_TOKEN"))
			err := client.DeleteItem(tc.itemName)
			log.Println(err)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			_, err = client.GetItem(tc.itemName)
			assert.Error(t, err)
		})
	}
}
