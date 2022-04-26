package database

import (
	"fmt"

	"jane_tech/internal/config"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	IN_MEMORY_DB = "inmemory"
)

// Underlying database interface. For now I use in memory, but this interface could be applied to something like a BigTable or other SQL/NoSQL DB
type Database interface {
	GetListOfCampaigns(key string) ([]*Campaign, error)
	AppendToListOfCampaigns(key string, campaigns *Campaign) error
	UpdateCampaignImpression(key string) error
	GetCampaignImpressions(key string) (int64, error)
}

type DatabaseConnector struct {
	db  Database
	log *log.Logger
}

// Add campaign to the database. Sets campaign id and for each keyword add campaign to the list of appropriate campaigns
func (connector *DatabaseConnector) AddCampaign(campaign *Campaign) error {
	campaign.CampaignId = uuid.New().String()
	for _, keyword := range campaign.TargetKeywords {
		connector.addCampaignKeyword(keyword, campaign)
	}

	return nil
}

// Add campaign to the list of campaigns matching provided keyword
func (connector *DatabaseConnector) addCampaignKeyword(word string, campaign *Campaign) error {
	return connector.db.AppendToListOfCampaigns(word, campaign)
}

// Fetching matching campaigns for a list of keywords
func (connector *DatabaseConnector) FetchMatchingCampaign(keywords []string) (*Campaign, error) {
	var bestCampaign *Campaign
	for _, keyword := range keywords {
		// TODO: place for optimization here. We fetch all the campaigns here because we have no info about underlying database and which optimizations we can use.
		// I used in memory hash map which can do fancy things like keep lists sorted and fast removal.
		// Some DB support this (BigTable always keep keys sorted), some DB not.
		// So I do no assumption here and fetch all the keys because there could be Inactive (this would be cleaned up by an external periodical job)
		// The only assumption here is that all the key would fit machine memory
		campaigns, err := connector.db.GetListOfCampaigns(keyword)
		if err != nil {
			return nil, err
		}
		connector.log.Infof("Fetched %v", campaigns)

		for _, campaign := range campaigns {
			if campaign.isActive() {
				connector.log.Infof("Active found %v", campaigns)
				if bestCampaign == nil || campaign.Less(bestCampaign) {
					bestCampaign = campaign
					connector.log.Infof("updated best campaign %v", bestCampaign)
				}
				break
			}
		}
	}

	return bestCampaign, nil
}

// Record impression for a specified campaign id
func (connector *DatabaseConnector) RecordImpression(campaignId string) error {
	return connector.db.UpdateCampaignImpression(campaignId)
}

// Returns a campaign impressions by campaign id
func (connector *DatabaseConnector) GetCampaignImpressions(campaignId string) (int64, error) {
	return connector.db.GetCampaignImpressions(campaignId)
}

func NewDatabaseConnector(cfg *config.Config, log *log.Logger) (*DatabaseConnector, error) {
	var db Database
	var err error
	switch cfg.Database.Type {
	case IN_MEMORY_DB:
		db, err = NewInMemoryDb(log)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unsupported DB type %s", cfg.Database.Type)
	}

	return &DatabaseConnector{
		log: log,
		db:  db,
	}, nil
}
