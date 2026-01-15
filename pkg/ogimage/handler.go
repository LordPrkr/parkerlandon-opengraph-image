package ogimage

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/ParkerGits/go-backend-starter/pkg/httplib"
	"github.com/ParkerGits/go-backend-starter/pkg/templates"
)

type Handler struct {
	generator *Generator
}

func NewHandler(generator *Generator) *Handler {
	return &Handler{
		generator: generator,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.URL.Query().Get("title"))
	subtitle := strings.TrimSpace(r.URL.Query().Get("subtitle"))

	if title == "" {
		httplib.WriteJSON(w, http.StatusBadRequest,
			httplib.NewBadRequestResponse("title query parameter is required"))
		return
	}

	params := templates.OGImageParams{
		Title:    title,
		Subtitle: subtitle,
		Handle:   "@lordprkr",
	}

	pngData, err := h.generator.Generate(r.Context(), params)
	if err != nil {
		slog.Error("failed to generate image", "error", err)
		httplib.WriteJSON(w, http.StatusInternalServerError,
			httplib.NewInternalErrorResponse("failed to generate image"))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400, s-maxage=604800, stale-while-revalidate=3600")
	w.WriteHeader(http.StatusOK)
	w.Write(pngData)
}
