package server

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Create router for our server
func NewRouter(impressionUrl *url.URL, server *Server) *chi.Mux {

	router := chi.NewRouter()
	// For each request set individual UUID for simplier debugging
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log := log.WithField("uuid", uuid.New().String())
			ctx := context.WithValue(req.Context(), "log", log)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	})

	router.Post("/campaign", server.createCampaign)
	router.Post("/addecision", server.addecission)
	router.Get("/campaign/{campaignId}", server.getCampaignImpressions)
	router.Get("/impression/record/{campaignId}", server.recordImpression)

	return router
}
