package lightning

import (
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type StandardClient struct {
	router chi.Router
}

func NewStandardClient() StandardClient {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	return StandardClient{router: r}
}

func (sc *StandardClient) RegisterWalletBalanceHandler(handler func(w http.ResponseWriter, r *http.Request)) {
	sc.router.Get("/walletbalance", handler)
}

func (sc *StandardClient) Start(port int) error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), sc.router)
	if err != nil {
		return errors.Wrap(err, "Starting StandardClient API")
	}
	return nil
}
