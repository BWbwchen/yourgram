package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testVerifyJWTCase = []struct {
	req            AuthorizationRequest
	expectedStatus int
}{
	{
		req: AuthorizationRequest{
			UserData: "",
			JWTToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRpbS5jaGVuYndAZ21haWwuY29tIiwiaXNzIjoiQldid2NoZW4ifQ.AO1T_FPVVTjZlhTPjiiRlKY0mYZLRCkziq5OlvFxx1I",
		},
		expectedStatus: http.StatusOK,
	},
	{
		req: AuthorizationRequest{
			UserData: "",
			JWTToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRpbS5jaGVuYndAZ21haWwuY29tIiwiaXNzIjoiQldid2NoZW4ifQ.sAO1T_FPVVTjZlhTPjiiRlKY0mYZLRCkziq5OlvFxx1I",
		},
		expectedStatus: http.StatusForbidden,
	},
}

func TestVerifyJWT(t *testing.T) {
	auth := NewService()
	for _, testcase := range testVerifyJWTCase {
		resp := auth.VerifyJWT(context.TODO(), testcase.req)
		assert.Equal(t, testcase.expectedStatus, resp.StatusCode, resp)
	}
}

var testCreateJWTCase = []struct {
	req            AuthorizationRequest
	expectedStatus int
}{
	{
		req: AuthorizationRequest{
			UserData: "bowei",
			JWTToken: "",
		},
		expectedStatus: http.StatusOK,
	},
}

func TestCreateJWT(t *testing.T) {
	auth := NewService()
	for _, testcase := range testCreateJWTCase {
		resp := auth.CreateJWT(context.TODO(), testcase.req)
		assert.Equal(t, testcase.expectedStatus, resp.StatusCode)
	}
}
