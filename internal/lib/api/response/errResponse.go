package response

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	HTTPStatusCode int    `json:"-"`
	Error          string `json:"error"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	render.SetContentType(render.ContentTypeJSON)
	return nil
}

func ErrInternal() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		Error:          http.StatusText(http.StatusInternalServerError),
	}
}

func ErrBadRequest(errorText string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		Error:          errorText,
	}
}

func ErrNotFound() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		Error:          http.StatusText(http.StatusNotFound),
	}
}
