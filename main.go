package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Krylphi/rockspoon-cart/internal/repository/mongo"
	"github.com/Krylphi/rockspoon-cart/internal/routing"

	"github.com/urfave/cli"
)

const (
	EnvConnectionString = "CONNSTR"
	EnvCartDatabase = "CARTDB"
	EnvCartNamespace = "CARTNAMESPACE"
)

func getEnv(env, def string) string {
	value, exists := os.LookupEnv(env)
	if !exists {
		value = def
	}
	return value
}



func run(c *cli.Context) error {
	repo, err := mongo.NewMongoRepository(&mongo.Config{
		ConnStr:         getEnv(EnvConnectionString, "mongodb://localhost:27017"),
		Database:        getEnv(EnvCartDatabase, "rockspoon-cart-0"),
		UsersCollection: getEnv(EnvCartNamespace, "carts"),
	})

	if err != nil {
		return err
	}

	srv := fmt.Sprint(getEnv("HOST", "0.0.0.0"),":", getEnv("PORT", "8080"))
	log.Printf("service address: %v", srv)
	return http.ListenAndServe(srv, routing.RouterInit(repo, repo))
}

func main() {
	var (
		configPath = "./config/.env"
		version    = "0.0.1"
	)

	var flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config, c",
			Value:       configPath,
			Usage:       "path to .env config file",
			Destination: &configPath,
		},
	}

	app := cli.NewApp()
	app.Name = "User Management API"
	app.Usage = "provide capabilities of creating updating, deleting shopping carts"
	app.UsageText = "rockspoon-cart [global options]"
	app.Version = version
	app.Flags = flags
	app.Action = run

	log.Fatal(app.Run(os.Args))
}
