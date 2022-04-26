package database

import (
	"fmt"
	"strings"
	"time"
)

// Represent campaign object
type Campaign struct {
	StartTimestamp int64    `json:"start_timestamp"`
	EndTimestamp   int64    `json:"end_timestamp"`
	TargetKeywords []string `json:"target_keywords"`
	MaxImpression  int64    `json:"max_impression"`
	Impressions    int64    `json:"-"`
	Cpm            float64  `json:"cpm"`
	CampaignId     string   `json:"-"`
}

func (c *Campaign) String() string {
	return fmt.Sprintf("Campaign(%d %d %v %d %f %s)", c.StartTimestamp, c.EndTimestamp, c.TargetKeywords, c.MaxImpression, c.Cpm, c.CampaignId)
}

func (c *Campaign) Validate() error {
	if !(c.StartTimestamp < c.EndTimestamp) || c.StartTimestamp == 0 || c.EndTimestamp == 0 {
		return fmt.Errorf("Incorrect time")
	}
	if c.Cpm <= 0 {
		return fmt.Errorf("incorrect CPM")
	}
	for _, word := range c.TargetKeywords {
		if strings.Contains(word, " ") {
			return fmt.Errorf("Incorrect target words")
		}
	}

	return nil
}

func (c *Campaign) isActive() bool {
	return time.Now().Before(time.Unix(c.EndTimestamp, 0)) && c.Impressions < c.MaxImpression
}

func (c1 *Campaign) Less(c2 *Campaign) bool {
	if c1.Cpm != c2.Cpm {
		return c1.Cpm > c2.Cpm
	}
	if c1.EndTimestamp != c2.EndTimestamp {
		return c1.EndTimestamp < c2.EndTimestamp
	}
	return strings.Compare(c1.CampaignId, c2.CampaignId) == -1
}
