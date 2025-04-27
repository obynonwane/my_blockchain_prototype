package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	// "github.com/rs/cors"
)

func (app *Config) publicRoutes() http.Handler {
	log.Println("starting public server")
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			// "http://localhost:3000",
			// "http://localhost:3001",
			// "http://localhost:3002",
			// "http://localhost:3003",
			"https://*", "http://*",
			"https://liquidity.algoralign.com",
			"https://admin.liquidity.algoralign.com",
			"https://pos.algoralign.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"X-PINGOTHER", "X-Requested-With", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))

	//other users link
	// mux.Post("/api/v1/auth/liquidity-signup-organisation", app.handleCreateLiquidityUserOrganisation)
	// mux.Post("/api/v1/auth/liquidity-signup-individual", app.handleCreateLiquidityUserIndividual)
	// mux.Post("/api/v1/auth/client-signup-remittance", app.handleCreateClientUserRemittance)
	// mux.Post("/api/v1/auth/client-signup-pos", app.handleCreateClientUserPos)

	return mux
}

func (app *Config) privateRoutes() http.Handler {
	log.Println("starting private server")
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			// "http://localhost:3000",
			// "http://localhost:3001",
			// "http://localhost:3002",
			// "http://localhost:3003",
			"https://*", "http://*",
			"https://liquidity.algoralign.com",
			"https://admin.liquidity.algoralign.com",
			"https://pos.algoralign.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"X-PINGOTHER", "X-Requested-With", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))

	//other users link
	// mux.Post("/api/v1/auth/liquidity-signup-organisation", app.handleCreateLiquidityUserOrganisation)
	// mux.Post("/api/v1/auth/liquidity-signup-individual", app.handleCreateLiquidityUserIndividual)
	// mux.Post("/api/v1/auth/client-signup-remittance", app.handleCreateClientUserRemittance)
	// mux.Post("/api/v1/auth/client-signup-pos", app.handleCreateClientUserPos)

	return mux
}
