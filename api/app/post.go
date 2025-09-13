package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LocTr/my-blog-be/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// The list of error types return by resource
var (
	ErrPostFailed = errors.New("post operation failed")
)

// PostStore defines the interface for post storage operations
type PostStore interface {
	GetPost(id int) (*models.Post, error)
	GetPosts(page, size int) ([]*models.Post, error)
}

// PostResource implement the post management
type PostResource struct {
	Store PostStore
}

// NewPostResource creates and returns an instance of PostResource
func NewPostResource(store PostStore) *PostResource {
	return &PostResource{Store: store}
}

func (rs *PostResource) router() chi.Router {
	r := chi.NewRouter()
	r.Get("/{postID}", rs.get)
	return r
}

func (rs *PostResource) postCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("running postCtx\n")
		id, err := strconv.Atoi(chi.URLParam(r, "postID"))
		fmt.Printf("ctx post: %+v\n", id)

		if err != nil {
			log(r).WithField("postID", id).Error(err)
			render.Render(w, r, ErrBadRequest)
			return
		}
		post, err := rs.Store.GetPost(id)
		if err != nil {
			log(r).WithField("postID", id).Error(err)
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxPost, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type postResponse struct {
	*models.Post
}

func newPostResponse(post *models.Post) *postResponse {
	return &postResponse{Post: post}
}

func (resource *PostResource) get(writer http.ResponseWriter, request *http.Request) {
	// post := request.Context().Value(ctxPost).(*models.Post)

	id, err := strconv.Atoi(chi.URLParam(request, "postID"))

	if err != nil {
		log(request).WithField("postID", id).Error(err)
		render.Render(writer, request, ErrBadRequest)
		return
	}

	post, err := resource.Store.GetPost(id)

	if err != nil {
		log(request).WithField("postID", id).Error(err)
		render.Render(writer, request, ErrNotFound)
		return
	}

	fmt.Printf("ctx post: %+v\n", post.ID)
	render.Respond(writer, request, newPostResponse(post))
}
