package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	v1 "github.com/obynonwane/my_blockchain_prototype/cmd/app/handlers/V1"
	appConfig "github.com/obynonwane/my_blockchain_prototype/cmd/config"
	"github.com/obynonwane/my_blockchain_prototype/cmd/database"
	"github.com/obynonwane/my_blockchain_prototype/cmd/genesis"
	"github.com/obynonwane/my_blockchain_prototype/cmd/state"
	"go.uber.org/zap"

	"github.com/obynonwane/my_blockchain_prototype/cmd/logger"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "main"

var counts int64

func main() {

	// Construct the application logger.
	log, err := logger.New("NODE")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1) // exit node if error
	}

}

func run(log *zap.SugaredLogger) error {

	log.Infow("starting service", "version", build)

	//===========================Connect to DB==========================================================================
	conn, err := connectToDB(log)
	if conn == nil {
		log.Errorw("connection to postgres DB", "status", "error connecting", "ERROR", err)
		log.Panic("can't connect to Postgres")
	}
	//====================================================================================================================

	//================================load the genesis file===============================================================
	genesis, err := genesis.Load()
	if err != nil {
		return err
	}
	//====================================================================================================================

	//=============================================define a function that runs like printf ===============================
	ev := func(v string, args ...any) {
		s := fmt.Sprintf(v, args...)                                    // Step 1: format string like printf
		log.Infow(s, "traceid", "00000000-0000-0000-0000-000000000000") // Step 2: structured log
	}
	//=====================================================================================================================

	//===================load the private key file of the onfigured beneficiary:node========================================
	path := fmt.Sprintf("%s%s.ecdsa", "cmd/zblock/accounts/", "miner1")
	privateKey, err := crypto.LoadECDSA(path)

	if err != nil {
		return fmt.Errorf("unable to load private key for node: %w", err)
	}
	//======================================================================================================================
	//================setup config==========================================================================================
	app := &appConfig.Config{
		DB:              conn,
		Models:          database.New(conn),
		Genesis:         genesis,
		SelectStrategy:  "Tip",                                               // add default of tip
		BeneficiaryID:   database.PublicKeyToAccountID(privateKey.PublicKey), // publick key of the node operator/beneficiary
		EvHandler:       ev,
		ReadTimeout:     5 * time.Second,  // 5 seconds read timeout
		WriteTimeout:    10 * time.Second, // 10 seconds write timeout
		IdleTimeout:     10 * time.Second, // 10 seconds idle timeout
		ShutdownTimeout: 20 * time.Second, // 20 seconds shutdown timeout
	}
	//========================================================================================================================

	//=================create new instance of state and inject depencies/config depencies=====================================
	state, err := state.New(*app)
	if err != nil {
		return err
	}
	defer state.Shutdown()
	//========================================================================================================================

	//========================== Inject config/depencies into routes file=====================================================
	routes := v1.NewRoutes(app, state)
	//========================================================================================================================

	//======create a buffered channel to hold the error from listening========================================================
	serverErrors := make(chan error, 1)
	//========================================================================================================================

	//=====================Start the service listening on public port.========================================================

	log.Infow("public port", "PORT", os.Getenv("PUB_PORT"))
	// define http server
	pub := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8080"),
		Handler: routes.PublicRoutes(),
	}
	go func() {

		// start the server
		err := pub.ListenAndServe()
		serverErrors <- err

	}()
	//=======================================================================================================================

	//====================== Start the service listening on private port.====================================================

	// start second server port
	// define http server
	prv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PRV_PORT")),
		Handler: routes.PrivateRoutes(),
	}
	go func() {

		// start the server
		err := prv.ListenAndServe()
		serverErrors <- err

	}()
	//=======================================================================================================================

	// =======================Start the service listening on web port.=======================================================

	// start second server port
	// define http server
	// web := &http.Server{
	// 	Addr:    fmt.Sprintf(":%s", os.Getenv("WEB_PORT")),
	// 	Handler: routes.WebRoutes(),
	// }
	// go func() {

	// 	// start the server
	// 	err := web.ListenAndServe()
	// 	serverErrors <- err

	// }()
	//=========================================================================================================================

	err = state.SeedGenesisAccount(&genesis)
	if err != nil {
		return fmt.Errorf("error creating default and genesis account to DB: %w", err)
	}

	//=======================Blocking main from exiting and accepting request - waiting for shutdown===========================
	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancelPub := context.WithTimeout(context.Background(), app.ShutdownTimeout)
		defer cancelPub()

		// Asking listener to shut down and shed load.
		log.Infow("shutdown", "status", "shutdown private API started")
		if err := prv.Shutdown(ctx); err != nil {
			prv.Close()
			return fmt.Errorf("could not stop private service gracefully: %w", err)
		}

		// Give outstanding requests a deadline for completion.
		ctx, cancelPri := context.WithTimeout(context.Background(), app.ShutdownTimeout)
		defer cancelPri()

		// Asking listener to shut down and shed load.
		log.Infow("shutdown", "status", "shutdown public API started")
		if err := pub.Shutdown(ctx); err != nil {
			pub.Close()
			return fmt.Errorf("could not stop public service gracefully: %w", err)
		}
	}

	return nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB(log *zap.SugaredLogger) (*sql.DB, error) {

	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")

	// Construct the DSN string
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Infow("Postgres not yet ready ...", "Status", "trying to connect to postgres DB")
			counts++
		} else {
			log.Infow("Connected to Postgres ...", "Status", "Connected DB")
			return connection, nil
		}

		if counts > 10 {
			log.Errorw("Error waiting to connect to DB", "Error", err)
			return nil, err
		}

		log.Infow("backing off for 2 seconds", "Status", "trying to connect to postgres DB")

		time.Sleep(2 * time.Second)
		continue
	}
}
