package database

import (
	"fmt"
	"testing"
	"time"
)

func TestCampaignIsActive(t *testing.T) {
	var testCases = []struct {
		campaign *Campaign
		isActive bool
	}{
		{
			campaign: &Campaign{StartTimestamp: 0, EndTimestamp: 100, Cpm: 1, CampaignId: "3"},
			isActive: false,
		},
		{
			campaign: &Campaign{StartTimestamp: 0, EndTimestamp: time.Date(2036, 1, 1, 1, 1, 1, 1, time.Local).Unix(), Cpm: 1, CampaignId: "3"},
			isActive: true,
		},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			if testCase.campaign.isActive() != testCase.isActive {
				t.Fatalf("incorrect isActive %v %v", testCase.campaign, testCase.isActive)
			}
		})
	}
}
