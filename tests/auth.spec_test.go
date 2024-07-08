package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cyberguru1/stage-two/config"
	"github.com/cyberguru1/stage-two/ent"
	"github.com/cyberguru1/stage-two/ent/migrate"
	"github.com/cyberguru1/stage-two/ent/user"
	"github.com/cyberguru1/stage-two/handlers"
	"github.com/cyberguru1/stage-two/middleware"
	"github.com/cyberguru1/stage-two/routes"
)

// let's mock handlers

// MockConfig is a mock of the config struct
type MockConfig struct {
	mock.Mock
}

var Handle *handlers.Handlers

func TestHandlers_UserRegister(t *testing.T) {

	app := SetupTestApp()

	type args map[string]interface{}

	tokenRes := ""

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		expectedCode int
	}{
		{
			name: "Successful Registration",
			args: args{
				"method": "POST",
				"route":  "/auth/register",
				"body": handlers.RegisterReq{
					Firstname: "John",
					Lastname:  "Doe",
					Email:     "john.doe@example.com",
					Password:  "password123",
					Phone:     "1234567890",
				},
			},
			wantErr:      false,
			expectedCode: 201,
		},
		{
			name: "Default Organisation",
			args: args{
				"method": "GET",
				"route":  "/api/organisations",
			},
			wantErr:      false,
			expectedCode: 200,
		},
		{
			name: "Missing Required Fields",
			args: args{
				"method": "POST",
				"route":  "/auth/register",
				"body": handlers.RegisterReq{
					Firstname: "John",
					Lastname:  "Doe",
					Email:     "",
					Password:  "password123",
					Phone:     "1234567890",
				},
			},
			wantErr:      true,
			expectedCode: 422,
		},
		{
			name: "Duplicate Email",
			args: args{
				"method": "POST",
				"route":  "/auth/register",
				"body": handlers.RegisterReq{
					Firstname: "John",
					Lastname:  "Doe",
					Email:     "john.doe@example.com",
					Password:  "password123",
					Phone:     "1234567890",
				},
			},
			wantErr:      true,
			expectedCode: 422,
		},
	}

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case
		method := test.args["method"].(string)
		var req *http.Request

		if method == "POST" {
			reqBodyBytes, err := json.Marshal(test.args["body"])

			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req = httptest.NewRequest(method, test.args["route"].(string), bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")

		} else {
			req = httptest.NewRequest(method, test.args["route"].(string), nil)
			req.Header.Set("Authorization", "Bearer "+tokenRes)

		}

		resp, _ := app.Test(req, -1)

		var response map[string]interface{}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Println(err)
		}

		if data, ok := response["data"].(map[string]interface{}); ok {
			if _, ok = data["accessToken"].(string); ok {
				tokenRes = data["accessToken"].(string)
			}
		}
		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.name)
		resp.Body.Close()
	}

}

func TestHandlers_UserLogin(t *testing.T) {
	app := SetupTestApp()

	type args map[string]interface{}

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "Successful Login",
			args: args{
				"method": "POST",
				"route":  "/auth/login",
				"body": map[string]string{
					"email":    "john.doe@example.com",
					"password": "password123",
				},
			},
			wantErr:      false,
			expectedCode: 200,
		},
		{
			name: "Invalid Credentials",
			args: args{
				"method": "POST",
				"route":  "/auth/login",
				"body": map[string]string{
					"email":    "john.doe@example.com",
					"password": "wrongpassword",
				},
			},
			wantErr:      true,
			expectedCode: 401,
		},
	}

	for _, test := range tests {
		method := test.args["method"].(string)
		var req *http.Request

		tokenRes := ""

		if method == "POST" {
			reqBodyBytes, err := json.Marshal(test.args["body"])
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req = httptest.NewRequest(method, test.args["route"].(string), bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req = httptest.NewRequest(method, test.args["route"].(string), nil)
		}

		resp, _ := app.Test(req, -1)

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Println(err)
		}

		if data, ok := response["data"].(map[string]interface{}); ok {
			if _, ok = data["accessToken"].(string); ok {
				tokenRes = data["accessToken"].(string)
				if tokenRes == "" {
					assert.Fail(t, test.name)
				}
			}
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.name)
		resp.Body.Close()

	}
}

func TestHandle_Lastdelete(t *testing.T) {

	email := "john.doe@example.com"
	deleteUser(email) // delete's the default creted user
}

func deleteUser(email string) {

	err := godotenv.Load("../.env")

	if err != nil {
		log.Print("Error loading .env file")
	}

	// Setup a postgres connection

	conf := config.New()

	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Name, conf.Database.Password))

	if err != nil {
		fmt.Println("Database connection failed : ", err)
	}

	ctx := context.Background()

	defer client.Close()

	// Run the auto migration tool to create the schema for the User entity.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Find the user by email
	u, err := client.User.Query().WithOrganisations().Where(user.Email(email)).Only(ctx)
	if err != nil {
		fmt.Println("failed to find user: ", err)
	}

	if err := client.Organisation.DeleteOne(u.Edges.Organisations[0]).Exec(ctx); err != nil {
		fmt.Errorf("failed to delete organization: %w", err)
	}

	// Delete the user
	if err := client.User.DeleteOne(u).Exec(ctx); err != nil {

		fmt.Errorf("failed to delete user: %w", err)
	}

}
func SetupTestApp() *fiber.App {
	// load .env file from given path

	err := godotenv.Load("../.env")

	if err != nil {
		log.Print("Error loading .env file")
	}

	// Setup a postgres connection

	conf := config.New()

	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Name, conf.Database.Password))

	if err != nil {
		fmt.Println("Database connection failed : ", err)
	}

	ctx := context.Background()

	err = client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)

	if err != nil {
		fmt.Println("Migration Fail: ", err)
	}

	// Create a server using fiber
	app := fiber.New()
	middleware.SetMiddleware(app) //setup middleware

	defer app.Shutdown()

	// create a new handler
	Handle = handlers.NewHandlers(client, conf)

	routes.SetupApiV1(app, Handle)

	return app
}
