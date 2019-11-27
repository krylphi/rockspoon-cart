package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Krylphi/rockspoon-cart/internal/repository/mongo"
	"github.com/Krylphi/rockspoon-cart/internal/routing"

	"github.com/gorilla/mux"
	"github.com/urfave/cli"
)

var (
	configPath = "./config/.env"
	version    = "0.0.1"
)

func getEnv(env, def string) string {
	value, exists := os.LookupEnv(env)
	if !exists {
		value = def
	}
	return value
}

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "config, c",
		Value:       configPath,
		Usage:       "path to .env config file",
		Destination: &configPath,
	},
}

func run(c *cli.Context) error {
	r := mux.NewRouter()

	repo, err := mongo.NewMongoRepository(&mongo.Config{
		ConnStr:         getEnv("CONNSTR", "mongodb://localhost:27017"),
		Database:        getEnv("USERDB", "rockspoon-cart-0"),
		UsersCollection: getEnv("USERNAMESPACE", "carts"),
	})

	if err != nil {
		return err
	}
	return http.ListenAndServe(getEnv("HOST", "0.0.0.0")+":"+getEnv("PORT", "8080"), routing.RouterInit(r, repo, repo))
}

func main() {
	app := cli.NewApp()
	app.Name = "User Management API"
	app.Usage = "provide capabilities of creating updating, deleting users"
	app.UsageText = "user-manager-api [global options]"
	app.Version = version
	app.Flags = flags
	app.Action = run

	log.Fatal(app.Run(os.Args))
}
