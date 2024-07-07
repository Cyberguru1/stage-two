package handlers

import (
	"github.com/cyberguru1/stage-two/config"
	"github.com/cyberguru1/stage-two/ent"
)

type Handlers struct {
	Client *ent.Client
	Config *config.Config
}

func NewHandlers(client *ent.Client, config *config.Config) *Handlers {
	return &Handlers{
		Client: client,
		Config: config,
	}
}
