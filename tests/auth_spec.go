package handlers

import (
	"testing"

	"github.com/cyberguru1/stage-two/config"
	"github.com/cyberguru1/stage-two/ent"
	"github.com/gofiber/fiber/v2"
)

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
			h := &Handlers{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			if err := h.UserLogin(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Handlers.UserLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			if err := h.UserRegister(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Handlers.UserRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
