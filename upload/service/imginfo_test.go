package upload_svc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testValidCase = []struct {
	req      ImgInfo
	expected bool
}{
	{
		req:      ImgInfo{},
		expected: false,
	},
	{
		req: ImgInfo{
			User:   "bwbwchen",
			ImgID:  "aaaaa",
			ImgURL: "aaaaa",
		},
		expected: true,
	},
	{
		req: ImgInfo{
			User:   "bwbwchen",
			ImgURL: "aaaaa",
		},
		expected: false,
	},
	{
		req: ImgInfo{
			ImgID:  "aaaaa",
			ImgURL: "aaaaa",
		},
		expected: false,
	},
	{
		req: ImgInfo{
			User:  "bwbwchen",
			ImgID: "aaaaa",
		},
		expected: false,
	},
}

func TestValid(t *testing.T) {

	for _, testcase := range testValidCase {
		resp := testcase.req.Valid()
		assert.Equal(t, testcase.expected, resp)
	}
}
