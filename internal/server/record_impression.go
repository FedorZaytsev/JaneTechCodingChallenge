package server

import (
	"jane_tech/internal/database"
	"net/http"

	"jane_tech/internal/logger"

	"github.com/go-chi/chi"
)

// handler for GET /campaign/[campaign-id]
func (s *Server) recordImpression(w http.ResponseWriter, req *http.Request) {
	log := logger.Entry(req.Context())

	campaignId := chi.URLParam(req, "campaignId")

	err := s.connector.RecordImpression(campaignId)
	log.Infof("err %v", err)
	switch err.(type) {
	case *database.CampaignNotFound:
		// in case of campaign not found return bad request
		w.WriteHeader(http.StatusBadRequest)
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		// in case of other errors, throw internal server error
		w.WriteHeader(http.StatusInternalServerError)
	}
}
