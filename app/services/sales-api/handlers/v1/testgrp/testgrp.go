package testgrp

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/cpustejovsky/service3-video/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of check enpoints.
type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		// return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
		panic("testing panic!")
	}
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	statusCode := http.StatusOK
	return web.Respond(ctx, w, status, statusCode)
}
