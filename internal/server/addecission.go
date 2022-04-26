package server

import (
	"encoding/json"
	"fmt"
	"jane_tech/internal/logger"
	"net/http"
	"path"
	"strings"
)

// request structure for POST /addecision
type addecissionReqest struct {
	Keywords []string `json:"keywords"`
}

func (r *addecissionReqest) Validate() error {
	for _, word := range r.Keywords {
		if strings.Contains(word, " ") {
			return fmt.Errorf("Incorrect target words")
		}
	}
	if len(r.Keywords) == 0 {
		return fmt.Errorf("Incorrect target words")
	}

	return nil
}

// Handler for POST /addecision
func (s *Server) addecission(w http.ResponseWriter, req *http.Request) {
	log := logger.Entry(req.Context())

	addecissionRequest := addecissionReqest{}

	// decode request
	err := json.NewDecoder(req.Body).Decode(&addecissionRequest)
	if err != nil {
		log.WithError(err).Errorln("Cannot decode body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate request
	err = addecissionRequest.Validate()
	if err != nil {
		log.WithError(err).Errorln("Invalid validation")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	campaign, err := s.connector.FetchMatchingCampaign(addecissionRequest.Keywords)
	if err != nil {
		log.WithError(err).Errorln("Cannot fetch mathing campaign")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	if campaign != nil {
		json.NewEncoder(w).Encode(struct {
			CampaignId    string `json:"campaign_id"`
			ImpressionUrl string `json:"impression_url"`
		}{
			CampaignId:    campaign.CampaignId,
			ImpressionUrl: s.generateImpressionUrl(campaign.CampaignId),
		})
	}
}

// generate impression url based on the campaign Id
func (s *Server) generateImpressionUrl(campaignId string) string {
	// get pre parsed impression url template
	impressionUrl := *s.impressionUrl
	// add campaign id to the end of the url
	impressionUrl.Path = path.Join(impressionUrl.Path, campaignId)
	return impressionUrl.String()
}
