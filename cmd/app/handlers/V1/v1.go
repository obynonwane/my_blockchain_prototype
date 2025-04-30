package v1

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1/private"
	"github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1/public"
	"github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1/web"
	"github.com/obynonwane/my_blockchain_prototype/cmd/config"
)

type AppConfig struct {
	App *config.Config
}

func NewRoutes(app *config.Config) *AppConfig {
	return &AppConfig{App: app}
}
func (app *AppConfig) PublicRoutes() http.Handler {

	pbl := public.Handlers{}
	log.Println("starting public server")
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://*", "http://*",
			"https://liquidity.algoralign.com",
			"https://admin.liquidity.algoralign.com",
			"https://pos.algoralign.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"X-PINGOTHER", "X-Requested-With", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))

	//other users link
	mux.Get("/v1/genesis/list", pbl.Genesis)

	return mux
}

func (app *AppConfig) PrivateRoutes() http.Handler {

	prv := private.Handlers{}
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://*", "http://*",
			"https://liquidity.algoralign.com",
			"https://admin.liquidity.algoralign.com",
			"https://pos.algoralign.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"X-PINGOTHER", "X-Requested-With", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))

	//routes
	mux.Get("/v1/genesis/list", prv.Genesis)

	return mux
}

func (app *AppConfig) WebRoutes() http.Handler {
	web := web.Handlers{}
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://*", "http://*",
			"https://liquidity.algoralign.com",
			"https://admin.liquidity.algoralign.com",
			"https://pos.algoralign.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"X-PINGOTHER", "X-Requested-With", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}))

	//routes
	mux.Get("/v1/genesis/list", web.Genesis)

	return mux
}
