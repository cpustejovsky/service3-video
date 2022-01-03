package testgrp

import (
	"context"
	"net/http"

	"github.com/cpustejovsky/service3-video/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of check enpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK2",
	}
	statusCode := http.StatusOK
	return web.Respond(ctx, w, status, statusCode)
}
