package database

import (
	"fmt"
	"jane_tech/internal/config"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestAddWord1(t *testing.T) {
	var testCases = []struct {
		campaigns []*Campaign
		expected  []string
	}{
		{
			campaigns: []*Campaign{
				{StartTimestamp: 0, EndTimestamp: 100, Cpm: 1, CampaignId: "3"},
				{StartTimestamp: 0, EndTimestamp: 100, Cpm: 10, CampaignId: "1"},
				{StartTimestamp: 0, EndTimestamp: 100, Cpm: 5, CampaignId: "2"}},
			expected: []string{"1", "2", "3"},
		},
		{
			campaigns: []*Campaign{
				{StartTimestamp: 0, EndTimestamp: 1, Cpm: 1, CampaignId: "1"},
				{StartTimestamp: 0, EndTimestamp: 2000, Cpm: 1, CampaignId: "4"},
				{StartTimestamp: 0, EndTimestamp: 100, Cpm: 1, CampaignId: "3"},
				{StartTimestamp: 0, EndTimestamp: 50, Cpm: 1, CampaignId: "2"}},
			expected: []string{"1", "2", "3", "4"},
		},
		{
			campaigns: []*Campaign{
				{StartTimestamp: 0, EndTimestamp: 1, Cpm: 1, CampaignId: "1"},
				{StartTimestamp: 0, EndTimestamp: 1, Cpm: 1, CampaignId: "4"},
				{StartTimestamp: 0, EndTimestamp: 1, Cpm: 1, CampaignId: "3"},
				{StartTimestamp: 0, EndTimestamp: 1, Cpm: 1, CampaignId: "2"}},
			expected: []string{"1", "2", "3", "4"},
		},
		{
			campaigns: []*Campaign{
				{StartTimestamp: 0, EndTimestamp: 1, Cpm: 1, CampaignId: "1"},
				{StartTimestamp: 0, EndTimestamp: 1, Cpm: 1, CampaignId: "4"},
				{StartTimestamp: 0, EndTimestamp: 1, Cpm: 2, CampaignId: "3"},
				{StartTimestamp: 0, EndTimestamp: 1, Cpm: 2, CampaignId: "2"}},
			expected: []string{"2", "3", "1", "4"},
		},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {

			inmem, _ := NewDatabaseConnector(&config.Config{
				Database: struct {
					Type   string
					Config config.DbConfig
				}{
					Type: IN_MEMORY_DB,
				},
			},
				log.New(),
			)
			for _, campaign := range testCase.campaigns {
				inmem.addCampaignKeyword("some_word", campaign)
			}

			campaigns, _ := inmem.db.GetListOfCampaigns("some_word")
			for idx, campaign := range campaigns {
				if campaign.CampaignId != testCase.expected[idx] {
					t.Fatalf("incorrect binary insert %v, expected %v", campaigns, testCase.expected)
				}
			}
		})
	}
}

func TestAddFetchCampaign(t *testing.T) {

	endTime := time.Date(2036, 1, 1, 1, 1, 1, 1, time.Local).Unix()
	var testCases = []struct {
		campaigns []*Campaign
		keywords  []string
		expected  int
	}{
		{
			campaigns: []*Campaign{
				{StartTimestamp: 1, EndTimestamp: endTime, Cpm: 1, TargetKeywords: []string{"iphone"}},
				{StartTimestamp: 2, EndTimestamp: endTime, Cpm: 2, TargetKeywords: []string{"android"}},
				{StartTimestamp: 3, EndTimestamp: endTime, Cpm: 1, TargetKeywords: []string{"android", "5G"}},
			},
			keywords: []string{"iphone"},
			expected: 1,
		},
		{
			campaigns: []*Campaign{
				{StartTimestamp: 1, EndTimestamp: endTime, Cpm: 1, CampaignId: "1", TargetKeywords: []string{"iphone"}},
				{StartTimestamp: 2, EndTimestamp: endTime, Cpm: 2, CampaignId: "2", TargetKeywords: []string{"android"}},
				{StartTimestamp: 3, EndTimestamp: endTime, Cpm: 1, CampaignId: "3", TargetKeywords: []string{"android", "5G"}},
			},
			keywords: []string{"5G"},
			expected: 3,
		},
		{
			campaigns: []*Campaign{
				{StartTimestamp: 1, EndTimestamp: endTime, Cpm: 1, CampaignId: "1", TargetKeywords: []string{"iphone"}},
				{StartTimestamp: 2, EndTimestamp: endTime, Cpm: 2, CampaignId: "2", TargetKeywords: []string{"android"}},
				{StartTimestamp: 3, EndTimestamp: endTime, Cpm: 1, CampaignId: "3", TargetKeywords: []string{"android", "5G"}},
			},
			keywords: []string{"android"},
			expected: 2,
		},
		{
			campaigns: []*Campaign{
				{StartTimestamp: 1, EndTimestamp: endTime, Cpm: 1, CampaignId: "1", TargetKeywords: []string{"iphone"}},
				{StartTimestamp: 2, EndTimestamp: 100, Cpm: 2, CampaignId: "2", TargetKeywords: []string{"android"}},
				{StartTimestamp: 3, EndTimestamp: endTime, Cpm: 1, CampaignId: "3", TargetKeywords: []string{"android", "5G"}},
			},
			keywords: []string{"android"},
			expected: 3,
		},
	}

	for idx, testCase := range testCases {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {

			db, _ := NewDatabaseConnector(&config.Config{
				Database: struct {
					Type   string
					Config config.DbConfig
				}{
					Type: IN_MEMORY_DB,
				},
			},
				log.New(),
			)
			for _, campaign := range testCase.campaigns {
				db.AddCampaign(campaign)
			}
			result, _ := db.FetchMatchingCampaign(testCase.keywords)
			if result == nil {
				t.Fatalf("Cannot find matching campaign for testsCase %v", testCase)
			}

			if int(result.StartTimestamp) != testCase.expected {
				t.Fatalf("Incorrect fetch result for testCase %v, got %v", testCase, result.StartTimestamp)
			}
		})
	}
}

func TestRecordImpressions(t *testing.T) {
	endTime := time.Date(2036, 1, 1, 1, 1, 1, 1, time.Local).Unix()
	db, _ := NewDatabaseConnector(&config.Config{
		Database: struct {
			Type   string
			Config config.DbConfig
		}{
			Type: IN_MEMORY_DB,
		},
	},
		log.New(),
	)

	campaign1 := &Campaign{StartTimestamp: 1, EndTimestamp: endTime, Cpm: 1, TargetKeywords: []string{"iphone"}}
	campaign2 := &Campaign{StartTimestamp: 2, EndTimestamp: endTime, Cpm: 1, TargetKeywords: []string{"android"}}

	db.AddCampaign(campaign1)
	db.AddCampaign(campaign2)

	impressions, err := db.GetCampaignImpressions(campaign1.CampaignId)
	if impressions != 0 || err != nil {
		t.Fatalf("Incorrect number of impressions %v expected %v", impressions, 0)
	}
	db.RecordImpression(campaign1.CampaignId)
	db.RecordImpression(campaign1.CampaignId)
	db.RecordImpression(campaign1.CampaignId)

	impressions, err = db.GetCampaignImpressions(campaign1.CampaignId)
	if impressions != 3 || err != nil {
		t.Fatalf("Incorrect number of impressions %v expected %v", impressions, 3)
	}
	db.RecordImpression(campaign2.CampaignId)

	impressions, err = db.GetCampaignImpressions(campaign2.CampaignId)
	if impressions != 1 || err != nil {
		t.Fatalf("Incorrect number of impressions %v expected %v", impressions, 1)
	}
	impressions, err = db.GetCampaignImpressions(campaign1.CampaignId)
	if impressions != 3 || err != nil {
		t.Fatalf("Incorrect number of impressions %v expected %v", impressions, 3)
	}
}
