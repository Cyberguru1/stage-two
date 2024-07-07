////////////////////////////
// Code by Cyb3rGuru      //
// APi server             //
//                        //
///////////////////////////

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/cyberguru1/stage-two/config"
	"github.com/cyberguru1/stage-two/ent"
	"github.com/cyberguru1/stage-two/ent/migrate"
	"github.com/cyberguru1/stage-two/handlers"
	"github.com/cyberguru1/stage-two/middleware"
	"github.com/cyberguru1/stage-two/routes"
	"github.com/cyberguru1/stage-two/utils"
)


func getInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// json.NewEncoder(w).Encode()
	return
}

func main() {

	// load .env file from given path

	err := godotenv.Load(".env")

	if err != nil {
		log.Print("Error loading .env file")
	}

	// Setup a postgres connection

	conf := config.New()

	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Name, conf.Database.Password))

	if err != nil {
		utils.Fatalf("Database connection failed : ", err)
	}

	defer client.Close()

	ctx := context.Background()

	err = client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)

	if err != nil {
		utils.Fatalf("Migration Fail: ", err)
	}



	// Create a server using fiber
	app := fiber.New()
	middleware.SetMiddleware(app) //setup middleware



	// create a new handler
	handler := handlers.NewHandlers(client, conf)
	
	routes.SetupApiV1(app, handler)

	port := "8080"

	addr := flag.String("addr", port, "http service address")
	flag.Parse()
	log.Fatal(app.Listen(":" + *addr))
}
