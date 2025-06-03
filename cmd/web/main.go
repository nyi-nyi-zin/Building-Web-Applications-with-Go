package main

import ( // လိုအပ်တဲ့ package တွေကို import လုပ်ထားတာပါ
	"flag"
	"fmt"
	"html/template"
	"log"
	"myapp/internal/driver"
	"myapp/internal/models"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"
const cssVersion = "1"

type config struct { // configuration အတွက် struct ဖွဲ့ထားတာပါ
	port int
	env  string
	api  string
	db   struct { // database configuration
		dsn string // database connection string
	}
	stripe struct { // Stripe payment configuration
		secret string // Stripe secret key
		key    string // Stripe public key
	}
}

type application struct { // main application struct
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	DB            models.DBModel
}

func (app *application) serve() error { // HTTP server ကို start လုပ်တဲ့ function
	srv := &http.Server{ // HTTP server configuration
		Addr:              fmt.Sprintf(":%d", app.config.port), // server port
		Handler:           app.routes(),                        // route handler
		IdleTimeout:       30 * time.Second,                    // connection idle timeout
		ReadTimeout:       10 * time.Second,                    // read timeout
		ReadHeaderTimeout: 5 * time.Second,                     // header read timeout
		WriteTimeout:      5 * time.Second,                     // write timeout
	}
	app.infoLog.Printf(fmt.Sprintf("Starting HTTP server in %s mode on port %d", app.config.env, app.config.port))

	return srv.ListenAndServe() // server ကို start လုပ်ပါ
}

func main() { // main function
	var cfg config // configuration variable

	// command line flags တွေကို setup လုပ်ထားတာပါ
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development | production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")
	flag.StringVar(&cfg.db.dsn, "dsn", "nyinyizin:secret@tcp(localhost:3306)/widgets?parseTime=true&tls=false", "DSN")

	flag.Parse() // flags တွေကို ခေါ်သုံးထားတာ။

	// Stripe keys တွေကို environment variables ကနေ ယူထားတာပါ
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRECT")

	// loggers တွေကို create လုပ်ထားတာပါ
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	tc := make(map[string]*template.Template) // template cache ကို create လုပ်ထားတာပါ

	// application struct ကို create လုပ်ပြီး initialize လုပ်ထားတာပါ
	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
		DB:            models.DBModel{DB: conn},
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}
