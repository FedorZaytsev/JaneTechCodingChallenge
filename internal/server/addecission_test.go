package server

import (
	"fmt"
	"net/url"
	"testing"
)

func TestGenerateImpressionUrl(t *testing.T) {
	impressionUrl, _ := url.Parse("http://localhost:8000/impression/record")

	var testCases = []struct {
		server     *Server
		campaignId string
		url        string
	}{
		{
			server:     &Server{impressionUrl: impressionUrl},
			campaignId: "qqq",
			url:        "http://localhost:8000/impression/record/qqq",
		},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			if testCase.server.generateImpressionUrl(testCase.campaignId) != testCase.url {
				t.Fatalf("Incorrect generateImpressionUrl %v %v", testCase.server.generateImpressionUrl(testCase.campaignId), testCase.url)
			}
		})
	}
}
