package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const (
	webPort  = "3000"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

type Config struct {
	//	Session *scs.SessionManager
	Models data.Models
	// Etcd    *clientv3.Client
}

func main() {
	// Connect to Mongo and get a client.
	log.Println("LOGGER_SERVICE----- START")

	mongoClient, err := connectToMongo()
	client = mongoClient

	// We'll use this context to disconnect from mongo, since it needs one.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection to Mongo when application exits
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// session := scs.New()
	// session.Lifetime = 24 * time.Hour
	// session.Cookie.Persist = true
	// session.Cookie.SameSite = http.SameSiteLaxMode
	// session.Cookie.Secure = false

	app := Config{
		//Session: session,
		Models: data.New(client),
	}
	app.serve()

	// go app.gRPCListen()
	log.Println("LOGGER_SERVICE----- END")

}

// serve starts the web server.
func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	fmt.Println("--------------------------------------")
	fmt.Println("Starting logging web service on port", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// Connect opens a connection to the Mongo database and returns a client.
func connectToMongo() (*mongo.Client, error) {
	// // get env
	mongoURL, err := utils.ViperEnvVariable("MONGODB_URL")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// create connect options
	clientOptions := options.Client().ApplyURI(mongoURL)
	// clientOptions.SetAuth(options.Credential{
	// 	Username: "admin",
	// 	Password: "password",
	// })

	// Connect to the MongoDB and return Client instance
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("mongo.Connect() ERROR:", err)
		return nil, err
	}

	return c, nil
}

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Mount("/", app.webRouter())
	//mux.Mount("/api",app.apiRouter)
	return mux
}

func (app *Config) webRouter() http.Handler {
	mux := chi.NewRouter()
	//mux.Use(app.SessionLoad)

	//mux.Get("/", app.LoginPage)
	//mux.Get("login", app.LoginPage)
	//mux.Post("/login", app.LoginPagePost)
	//mux.Get("logout", app.LogoutPage)
	mux.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{panic: PANIC}"))
		log.Println("PANIC")
	})
	// mux.Route("/admin", func(mux chi.Router) {
	// 	mux.Use(app.Auth)
	// 	//mux.Get("/dasboard",app.Dasboard)
	// })

	return mux

}
