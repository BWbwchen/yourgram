package account_svc

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCreateAccountCase = []struct {
	req            AuthRequest
	expectedStatus int
}{
	{
		req: AuthRequest{
			Email:    "",
			Name:     "testUser",
			Password: "",
		},
		expectedStatus: http.StatusBadRequest,
	},
	{
		req: AuthRequest{
			Email:    "test@gmail.com",
			Name:     "testUser",
			Password: "",
		},
		expectedStatus: http.StatusBadRequest,
	},
	{
		req: AuthRequest{
			Email:    "test@test.com",
			Name:     "testUser",
			Password: "asdf",
		},
		expectedStatus: http.StatusOK,
	},
}

func TestCreateAccount(t *testing.T) {
	auth := NewService()
	for _, testcase := range testCreateAccountCase {
		resp := auth.CreateAccount(context.TODO(), testcase.req)
		assert.Equal(t, testcase.expectedStatus, resp.StatusCode)
	}
}

var testUserLoginCase = []struct {
	req            AuthRequest
	expectedStatus int
}{
	{
		req: AuthRequest{
			Email:    "test@test.com",
			Name:     "testUser",
			Password: "asdf",
		},
		expectedStatus: http.StatusOK,
	},
	{
		req: AuthRequest{
			Email:    "",
			Name:     "testUser",
			Password: "asdf",
		},
		expectedStatus: http.StatusOK,
	},
	{
		req: AuthRequest{
			Email:    "test@test.com",
			Name:     "",
			Password: "asdf",
		},
		expectedStatus: http.StatusOK,
	},
	{
		req: AuthRequest{
			Email:    "test@test.com",
			Name:     "",
			Password: "asdfa",
		},
		expectedStatus: http.StatusUnauthorized,
	},
}

func TestUserLogin(t *testing.T) {
	auth := NewService()
	for _, testcase := range testUserLoginCase {
		resp := auth.UserLogin(context.TODO(), testcase.req)
		assert.Equal(t, testcase.expectedStatus, resp.StatusCode)
	}
}
