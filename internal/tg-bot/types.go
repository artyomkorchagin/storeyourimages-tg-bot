package bot

import (
	"github.com/artyomkorchagin/storeyourimages/internal/services/content"
	"github.com/artyomkorchagin/storeyourimages/internal/services/users"
)

type AllServices struct {
	Users   *users.Service
	Content *content.Service
}

func NewAllServices(users *users.Service, content *content.Service) *AllServices {
	return &AllServices{
		Users:   users,
		Content: content,
	}
}
