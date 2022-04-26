package server

import (
	"fmt"
	"jane_tech/internal/database"
	"net/http"

	"github.com/go-chi/chi"
)

// handler for GET [impression-url]
func (s *Server) getCampaignImpressions(w http.ResponseWriter, req *http.Request) {
	campaignId := chi.URLParam(req, "campaignId")

	impressions, err := s.connector.GetCampaignImpressions(campaignId)
	switch err.(type) {
	case *database.CampaignNotFound:
		// in case of campaign not found return bad request
		w.WriteHeader(http.StatusBadRequest)
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d", impressions)))
	default:
		// in case of other errors, throw internal server error
		w.WriteHeader(http.StatusInternalServerError)
	}
}
