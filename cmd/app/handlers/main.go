package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
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

const publicPort = "8080"
const privatePort = "8081"
const webPort = "8082"

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
	conn := connectToDB()
	if conn == nil {
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
		DB:             conn,
		Models:         database.New(conn),
		Genesis:        genesis,
		SelectStrategy: "Tip",                                               // add default of tip
		BeneficiaryID:  database.PublicKeyToAccountID(privateKey.PublicKey), // publick key of the node operator/beneficiary
		EvHandler:      ev,
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
	//=====================Start the service listening on public port.========================================================
	go func() {
		// define http server
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", publicPort),
			Handler: routes.PublicRoutes(),
		}

		// start the server
		err := srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}

	}()
	//=======================================================================================================================

	//====================== Start the service listening on private port.====================================================
	go func() {
		// start second server port
		// define http server
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", privatePort),
			Handler: routes.PrivateRoutes(),
		}

		// start the server
		err := srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}

	}()
	//=======================================================================================================================

	// =======================Start the service listening on web port.=======================================================
	go func() {
		// start second server port
		// define http server
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", webPort),
			Handler: routes.WebRoutes(),
		}

		// start the server
		err := srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}

	}()
	//=========================================================================================================================

	//=======================Prevent main from exiting and accepting request===================================================
	select {}
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

func connectToDB() *sql.DB {

	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")
	log.Println("db user", dbUser)

	// Construct the DSN string
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	log.Printf("%s", dsn)
	// dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres ...")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("backing off for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
