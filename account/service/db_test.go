package account_svc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testValidCase = []struct {
	input    UserInfo
	expected bool
}{
	{
		input: UserInfo{
			Email:    "",
			Name:     "testUser",
			Password: "",
		},
		expected: false,
	},
	{
		input: UserInfo{
			Email:    "test@gmail.com",
			Name:     "testUser",
			Password: "",
		},
		expected: false,
	},
	{
		input: UserInfo{
			Email:    "",
			Name:     "",
			Password: "",
		},
		expected: false,
	},
	{
		input: UserInfo{
			Email:    "test@test.com",
			Name:     "testUser",
			Password: "asdf",
		},
		expected: true,
	},
}

func TestValid(t *testing.T) {
	service := DBStruct{}
	for _, testcase := range testValidCase {
		resp := service.valid(testcase.input)
		assert.Equal(t, testcase.expected, resp)
	}
}
