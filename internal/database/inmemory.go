package database

import (
	log "github.com/sirupsen/logrus"
)

var CAMPAIGN_NOT_FOUND = CampaignNotFound{}

type CampaignNotFound struct{}

func (e *CampaignNotFound) Error() string {
	return "Campaign not found"
}

// In memory database. Uses 2 hashmaps to store data
type InMemoryDb struct {
	// Mapping from a target word to the list of matching campaigns
	// Those campaigns are stored in a sorted fachion according to the Campaign.Less specification
	wordsToCampigns map[string][]*Campaign
	// Simple mapping from campaign id to campaign itself
	campaigns map[string]*Campaign
	log       *log.Logger
}

// returns list for campaigns which match requested keyword
func (db *InMemoryDb) GetListOfCampaigns(key string) ([]*Campaign, error) {
	if value, exists := db.wordsToCampigns[key]; !exists {
		return []*Campaign{}, nil
	} else {
		return value, nil
	}
}

// add campaign to the list of campaigns for a specified word
func (db *InMemoryDb) AppendToListOfCampaigns(key string, campaign *Campaign) error {
	if _, exists := db.wordsToCampigns[key]; !exists {
		db.wordsToCampigns[key] = []*Campaign{}
	}

	listOfCampaigns := db.wordsToCampigns[key]

	db.log.Infof("AppendToListOfCampaigns: start %v %v\n", listOfCampaigns, campaign)

	// because we can control our data, I do binary search in order to insert campaign in an appropriate position
	position := searchPosition(listOfCampaigns, campaign)
	if position == len(listOfCampaigns) {
		// if end of the list simply add it
		db.wordsToCampigns[key] = append(listOfCampaigns, campaign)
	} else {
		// if we want to insert in the middle, increase list by one and put random element at needed position
		listOfCampaigns = append(listOfCampaigns[:position+1], listOfCampaigns[position:]...)
		// replace random object in the middle with our object
		listOfCampaigns[position] = campaign
		db.log.Infof("AppendToListOfCampaigns: done %v\n", listOfCampaigns)
		db.wordsToCampigns[key] = listOfCampaigns
	}

	db.campaigns[campaign.CampaignId] = campaign

	return nil
}

// Increment campaign impressions counter by campaign id
func (db *InMemoryDb) UpdateCampaignImpression(key string) error {
	db.log.Infof("UpdateCampaignImpression: db.campaigns %v %v\n", db.campaigns, key)
	campaign, exists := db.campaigns[key]
	if !exists {
		return &CAMPAIGN_NOT_FOUND
	}
	campaign.Impressions += 1
	return nil
}

// Get number of impressions by campaign id
func (db *InMemoryDb) GetCampaignImpressions(key string) (int64, error) {
	campaign, exists := db.campaigns[key]
	if !exists {
		return 0, &CAMPAIGN_NOT_FOUND
	}
	return campaign.Impressions, nil
}

// Binary search position in the array of Campaign where campaign should be inserted
func searchPosition(campaigns []*Campaign, campaign *Campaign) int {
	start := 0
	end := len(campaigns) - 1

	for start <= end {
		median := (start + end) / 2

		if campaigns[median].Less(campaign) {
			start = median + 1
		} else {
			end = median - 1
		}
	}

	return start
}

// Creates new in memory database
func NewInMemoryDb(log *log.Logger) (*InMemoryDb, error) {
	return &InMemoryDb{
		wordsToCampigns: make(map[string][]*Campaign),
		campaigns:       make(map[string]*Campaign),
		log:             log,
	}, nil
}
