package upload_svc

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCheckImageTypeCase = []struct {
	req      []byte
	expected bool
}{
	{
		req:      getJPG(),
		expected: true,
	},
	{
		req:      getPNG(),
		expected: true,
	},
	{
		req:      getErrorImg(),
		expected: false,
	},
}

func TestCheckImageType(t *testing.T) {
	for _, testcase := range testCheckImageTypeCase {
		resp := checkImageType(testcase.req)
		assert.Equal(t, testcase.expected, resp)
	}
}

func getJPG() []byte {
	jpgImg, _ := ioutil.ReadFile("../lena.jpg")
	return jpgImg
}

func getPNG() []byte {
	pngImg, _ := ioutil.ReadFile("../lena.png")
	return pngImg
}

func getErrorImg() []byte {
	errorImg, _ := ioutil.ReadFile("../main.go")
	return errorImg
}
