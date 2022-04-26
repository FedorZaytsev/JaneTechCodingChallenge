package server

import (
	"encoding/json"
	"jane_tech/internal/database"
	"jane_tech/internal/logger"
	"net/http"
)

// handler for POST /campaign request
func (s *Server) createCampaign(w http.ResponseWriter, req *http.Request) {
	log := logger.Entry(req.Context())

	campaign := database.Campaign{}

	// parse request into campaign structure
	err := json.NewDecoder(req.Body).Decode(&campaign)
	if err != nil {
		log.WithError(err).Errorln("Cannot decode body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate that the object is correct
	err = campaign.Validate()
	if err != nil {
		log.WithError(err).Errorln("Invalid validation")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// add to the storage
	err = s.connector.AddCampaign(&campaign)
	if err != nil {
		log.WithError(err).Errorln("Cannot process campaign")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		CampaignId string `json:"campaign-id"`
	}{
		CampaignId: campaign.CampaignId,
	})

}
