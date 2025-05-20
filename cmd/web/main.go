package main // ဒါက main package ဖြစ်ကြောင်း ကြေငြာထားတာပါ

import ( // လိုအပ်တဲ့ package တွေကို import လုပ်ထားတာပါ
	"flag"          // command line flags တွေကို handle လုပ်ဖို့
	"fmt"           // string formatting အတွက်
	"html/template" // HTML templates တွေကို handle လုပ်ဖို့
	"log"           // logging အတွက်
	"net/http"      // HTTP server အတွက်
	"os"            // operating system functions တွေအတွက်
	"time"          // time related functions တွေအတွက်
)

const version = "1.0.0"  // application version
const cssVersion = "1"   // CSS version

type config struct {  // configuration အတွက် struct ဖွဲ့ထားတာပါ
	port int    // server port number
	env  string // environment (development/production)
	api  string // API URL
	db   struct {  // database configuration
		dsn string  // database connection string
	}
	stripe struct {  // Stripe payment configuration
		secrect string  // Stripe secret key
		key     string  // Stripe public key
	}
}

type application struct {  // main application struct
	config        config  // configuration
	infoLog       *log.Logger  // information logger
	errorLog      *log.Logger  // error logger
	templateCache map[string]*template.Template  // HTML templates cache
	version       string  // application version
}

func (app *application) serve() error {  // HTTP server ကို start လုပ်တဲ့ function
	srv := &http.Server{  // HTTP server configuration
		Addr: fmt.Sprintf(":%d",app.config.port),  // server port
		Handler:app.routes(),  // route handler
		IdleTimeout:30 * time.Second,  // connection idle timeout
		ReadTimeout:10 * time.Second,  // read timeout
		ReadHeaderTimeout: 5 *time.Second,  // header read timeout
		WriteTimeout: 5 *time.Second,  // write timeout
	}
	app.infoLog.Println(fmt.Sprintf("Starting HTTP server in %s mode on port %d", app.config.env, app.config.port))

	return srv.ListenAndServe()  // server ကို start လုပ်ပါ
}

func main() { // main function
	var cfg config // configuration variable

	// command line flags တွေကို set လုပ်ထားတာပါ
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development | production}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")

	flag.Parse() // flags တွေကို parse လုပ်ပါ

	// Stripe keys တွေကို environment variables ကနေ ယူထားတာပါ
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secrect = os.Getenv("STRIPE_SECRECT")

	// loggers တွေကို create လုပ်ထားတာပါ
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	tc := make(map[string]*template.Template) // template cache ကို create လုပ်ထားတာပါ

	// application struct ကို create လုပ်ပြီး initialize လုပ်ထားတာပါ
	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
	}

	// server ကို start လုပ်ပါ
	err := app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
}
