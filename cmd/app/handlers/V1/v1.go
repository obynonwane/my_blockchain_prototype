package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1/private"
	"github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1/public"
	"github.com/obynonwane/my_blockchain_prototype/cmd/config"
	"github.com/obynonwane/my_blockchain_prototype/cmd/state"
)

type AppConfig struct {
	App   *config.Config
	State *state.State
}

func NewRoutes(app *config.Config, st *state.State) *AppConfig {
	return &AppConfig{App: app, State: st}
}
func (app *AppConfig) PublicRoutes() http.Handler {

	pbl := public.Handlers{
		State: app.State,
	}

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
	mux.Get("/v1/accounts/list", pbl.Accounts)
	mux.Get("/v1/accounts/list/{account}", pbl.Accounts)
	mux.Get("/v1/tx/uncommitted/list", pbl.Mempool)
	mux.Post("/v1/create/user", pbl.CreateUser)

	return mux
}

func (app *AppConfig) PrivateRoutes() http.Handler {

	prv := private.Handlers{
		Model: &app.App.Models,
	}
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

// func (app *AppConfig) WebRoutes() http.Handler {
// 	web := web.Handlers{
// 		Model: &app.App.Models,
// 	}
// 	mux := chi.NewRouter()
// 	mux.Use(cors.Handler(cors.Options{
// 		AllowedOrigins: []string{
// 			"https://*", "http://*",
// 			"https://liquidity.algoralign.com",
// 			"https://admin.liquidity.algoralign.com",
// 			"https://pos.algoralign.com"},
// 		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowedHeaders: []string{"X-PINGOTHER", "X-Requested-With", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
// 		ExposedHeaders: []string{"Link"},
// 	}))

// 	//routes
// 	mux.Get("/v1/genesis/list", web.Genesis)

// 	return mux
// }
