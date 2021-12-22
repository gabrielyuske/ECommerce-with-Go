package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const version = "1.0.0"
const cssVersion = "1"

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		Secret string
		Key    string
	}
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Println(fmt.Sprintf("Starting HTTP server in %s mode on port %d", app.config.env, app.config.port))

	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Printf("Could not load .env file")
		os.Exit(1)
	}

	return srv.ListenAndServe()
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application enviornment {development|production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")

	flag.Parse()

	// cfg.stripe.Key = os.Getenv("STRIPE_KEY")
	cfg.stripe.Key = "sk_test_51K3kSFA4X8RRJipI3Cszt1Ed9sUu811UUS3yPz9CIHDMpzBezp8IDdUYhsNZafAKvRFOLfETaeYPmY6q8zynZf7D00FsAzRCjl"
	cfg.stripe.Secret = os.Getenv("STRIPE_SECRET")

	envMap, mapErr := godotenv.Read(".env")
	if mapErr != nil {
		fmt.Printf("Error loading .env into map[string]string\n")
		os.Exit(1)
	}
	fmt.Printf("API_KEY %s\n", envMap["TESTE"])

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
	}

	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}

}
