package authentication_service

import (
	"context"
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
		expectedStatus: 400,
	},
	{
		req: AuthRequest{
			Email:    "test@gmail.com",
			Name:     "testUser",
			Password: "",
		},
		expectedStatus: 400,
	},
	{
		req: AuthRequest{
			Email:    "test@test.com",
			Name:     "testUser",
			Password: "asdf",
		},
		expectedStatus: 200,
	},
}

func TestCreateAccount(t *testing.T) {
	auth := NewService()
	for _, testcase := range testCreateAccountCase {
		resp := auth.CreateAccount(context.TODO(), testcase.req)
		assert.Equal(t, testcase.expectedStatus, resp.StatusCode)
	}
}
