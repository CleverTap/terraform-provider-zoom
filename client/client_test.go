package client

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"terraform-provider-zoom/token"
	"testing"
)

func TestClient_GetUser(t *testing.T) {
	testCases := []struct {
		testName     string
		email        string
		expectErr    bool
		expectedResp *User
	}{
		{
			testName:  "user exists",
			email:     "user@gmail.com",
			expectErr: false,
			expectedResp: &User{
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
			email:        "user@gmail.com",
			expectErr:    true,
			expectedResp: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			apiToken := token.GenerateToken(os.Getenv("ZOOM_API_SECRET"), os.Getenv(os.Getenv("ZOOM_API_KEY")))
			client := NewClient(apiToken, 2)
			user, err := client.GetUser(tc.email)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResp, user)
		})
	}
}

func TestClient_NewItem(t *testing.T) {
	testCases := []struct {
		testName  string
		user      *User
		expectErr bool
	}{
		{
			testName: "user creation successful",
			user: &User{
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
			expectErr: false,
		},
		{
			testName: "user already exists",
			user: &User{
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
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			apiToken := token.GenerateToken(os.Getenv("ZOOM_API_SECRET"), os.Getenv(os.Getenv("ZOOM_API_KEY")))
			client := NewClient(apiToken, 2)
			err := client.NewUser(tc.user)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			user, err := client.GetUser(tc.user.Email)
			assert.NoError(t, err)
			assert.Equal(t, tc.user, user)
		})
	}
}

func TestClient_UpdateUser(t *testing.T) {
	testCases := []struct {
		testName  string
		user      *User
		expectErr bool
	}{
		{
			testName: "user exists",
			user: &User{
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
			expectErr: false,
		},
		{
			testName: "user does not exist",
			user: &User{
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
			expectErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			apiToken := token.GenerateToken(os.Getenv("ZOOM_API_SECRET"), os.Getenv(os.Getenv("ZOOM_API_KEY")))
			client := NewClient(apiToken, 2)
			err := client.UpdateUser(tc.user.Email, tc.user)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			user, err := client.GetUser(tc.user.Email)
			assert.NoError(t, err)
			assert.Equal(t, tc.user, user)
		})
	}
}

func TestClient_DeleteUser(t *testing.T) {
	testCases := []struct {
		testName  string
		email     string
		expectErr bool
	}{
		{
			testName:  "user exists",
			email:     "user@gmail.com",
			expectErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			apiToken := token.GenerateToken(os.Getenv("ZOOM_API_SECRET"), os.Getenv(os.Getenv("ZOOM_API_KEY")))
			client := NewClient(apiToken, 2)
			err := client.DeleteUser(tc.email, "pending")
			log.Println(err)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			_, err = client.GetUser(tc.email)
			assert.Error(t, err)
		})
	}
}
