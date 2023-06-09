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
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the API for lnd client"))
	})
	return StandardClient{router: r}
}

func (sc *StandardClient) RegisterWalletBalanceHandler(handler func(w http.ResponseWriter, r *http.Request)) {
	sc.router.Get("/walletbalance", handler)
}

func (sc *StandardClient) RegisterNewAddressHandler(handler func(w http.ResponseWriter, r *http.Request)) {
	sc.router.Post("/newaddress", handler)
}

func (sc *StandardClient) RegisterPubKeyHandler(handler func(w http.ResponseWriter, r *http.Request)) {
	sc.router.Get("/pubkey", handler)
}

func (sc *StandardClient) RegisterConnectPeerHandler(handler func(w http.ResponseWriter, r *http.Request)) {
	sc.router.Post("/connectpeer", handler)
}

func (sc *StandardClient) RegisterOpenChannelHandler(handler func(w http.ResponseWriter, r *http.Request)) {
	sc.router.Post("/openchannel", handler)
}

func (sc *StandardClient) Start(port int) error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), sc.router)
	if err != nil {
		return errors.Wrap(err, "Starting StandardClient API")
	}
	return nil
}
