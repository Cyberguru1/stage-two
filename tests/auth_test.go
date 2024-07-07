package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"

	"github.com/cyberguru1/stage-two/config"
	"github.com/cyberguru1/stage-two/ent"
	"github.com/cyberguru1/stage-two/ent/migrate"
	"github.com/cyberguru1/stage-two/handlers"
)

func TestHandlers_UserRegister(t *testing.T) {
	type fields struct {
		Client *ent.Client
		Config *config.Config
	}
	type args struct {
		ctx *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Successful Registration",
			fields: fields{
				Client: setupTestClient(t),
				Config: &config.Config{},
			},
			args: args{
				ctx: setupTestContext("POST", "/auth/register", map[string]string{
					"firstName": "John",
					"lastName":  "Doe",
					"email":     "john.doe@example.com",
					"password":  "password123",
					"phone":     "1234567890",
				}),
			},
			wantErr: false,
		},
		{
			name: "Missing Required Fields",
			fields: fields{
				Client: setupTestClient(t),
				Config: &config.Config{},
			},
			args: args{
				ctx: setupTestContext("POST", "/auth/register", map[string]string{
					"firstName": "John",
					"lastName":  "Doe",
					"email":     "",
					"password":  "password123",
					"phone":     "1234567890",
				}),
			},
			wantErr: true,
		},
		{
			name: "Duplicate Email",
			fields: fields{
				Client: setupTestClientWithExistingUser(t),
				Config: &config.Config{},
			},
			args: args{
				ctx: setupTestContext("POST", "/auth/register", map[string]string{
					"firstName": "John",
					"lastName":  "Doe",
					"email":     "john.doe@example.com",
					"password":  "password123",
					"phone":     "1234567890",
				}),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewHandlers(tt.fields.Client, tt.fields.Config)
			if err := h.UserRegister(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Handlers.UserRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandlers_UserLogin(t *testing.T) {
	type fields struct {
		Client *ent.Client
		Config *config.Config
	}
	type args struct {
		ctx *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handlers.NewHandlers(tt.fields.Client, tt.fields.Config)
			if err := h.UserLogin(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Handlers.UserLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func setupTestClient(t *testing.T) *ent.Client {

	// Setup a postgres connection
	conf := config.New()

	client, err := ent.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			conf.Database.Host,
			conf.Database.Port,
			conf.Database.User,
			conf.Database.Name,
			conf.Database.Password,
		))

	if err != nil {
		// utils.Fatalf("Database connection failed : ", err)
		fmt.Println(err)
	}

	ctx := context.Background()

	err = client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)

	return client
}

func setupTestClientWithExistingUser(t *testing.T) *ent.Client {
	client := setupTestClient(t)
	ctx := context.Background()
	_, err := client.User.Create().
		SetEmail("john.doe@example.com").
		SetFirstName("John").
		SetLastName("Doe").
		SetPassword("password123").
		SetPhone("1234567890").
		Save(ctx)
	if err != nil {
		fmt.Println(err)
	}
	require.NoError(t, err)
	return client
}

func setupTestContext(method, path string, body map[string]string) *fiber.Ctx {
	app := fiber.New()
	jsonBody, _ := json.Marshal(body)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(path)
	req.Header.SetMethod(method)
	req.Header.SetContentType("application/json")
	req.SetBody(jsonBody)
	resp := fasthttp.AcquireResponse()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{
		Request:  *req,
		Response: *resp,
	})

	return ctx
}
