package app

import (
	"net/http"

	"github.com/LocTr/my-blog-be/database"
	"github.com/LocTr/my-blog-be/logging"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type ctxKey int

const (
	ctxPost ctxKey = iota
)

// API provides application resources and handlers.
type API struct {
	Post *PostResource
}

// NewAPI configures and returns application API.
func NewAPI(db *bun.DB) (*API, error) {
	PostStore := database.NewPostStore(db)

	api := &API{
		Post: NewPostResource(PostStore),
	}
	return api, nil
}

// Router provides application routes
func (api *API) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Mount("/", api.Post.router())

	return router
}

func log(request *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(request)
}
