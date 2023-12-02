// The http package is responsible for initializing the server, the router with handlers, and for processing requests.
package http

import (
	"billing-service/internal/domain/models"
	"billing-service/internal/ports"
	"billing-service/pkg/infra/logger"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/sync/errgroup"
)

type Adapter struct {
	s          *http.Server
	l          net.Listener
	billingSvc ports.BillingService
	authURL    string
	client     *http.Client
}

type AdapterOptions struct {
	HTTP_port   int
	Timeout     time.Duration
	IdleTimeout time.Duration
	AuthURL     string
}

var router *gin.Engine

func GetRouter() *gin.Engine {
	return router
}

// New instantiates the adapter.
func New(billingService ports.BillingService, opts AdapterOptions) (*Adapter, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.HTTP_port))
	if err != nil {
		return nil, fmt.Errorf("server start failed: %w", err)
	}

	router := gin.Default()
	server := http.Server{
		Handler:      router,
		ReadTimeout:  opts.Timeout,
		WriteTimeout: opts.Timeout,
		IdleTimeout:  opts.IdleTimeout, // client connection lifetime
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	a := Adapter{
		s:          &server,
		l:          l,
		billingSvc: billingService,
		authURL:    opts.AuthURL,
		client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Jar:       jar,
		},
	}
	err = initRouter(&a, router)
	return &a, err
}

// Start starts an http server that accepts incoming connections on the Listener.
func (a *Adapter) Start() error {
	logger.Get().Info("starting http server...")

	eg := &errgroup.Group{}
	eg.Go(func() error {
		return a.s.Serve(a.l)
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// Stop stops the http server.
func (a *Adapter) Stop(ctx context.Context) error {
	var (
		err  error
		once sync.Once
	)
	once.Do(func() {
		err = a.s.Shutdown(ctx)
	})
	return err
}

// Verify receives access and refresh tokens from the Authorization header, and makes a request to the auth-service to authenticate the user.
func (a *Adapter) Verify(ctx *gin.Context) error {
	// get the access and refresh tokens
	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader == "" {
		return models.ErrBadRequest
	}
	tokens := strings.Split(authorizationHeader, "Bearer ")[1]
	access, refresh := strings.Split(tokens, ";")[0], strings.Split(tokens, ";")[1]

	// form a request
	authR, err := http.NewRequestWithContext(ctx.Request.Context(), "POST", a.authURL+"verify", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	authR.Header.Set("Authorization", "Bearer "+fmt.Sprintf("%s;%s", access, refresh))

	// execute the request
	resp, err := a.client.Do(authR)
	if err != nil || resp.StatusCode != http.StatusOK {
		return models.ErrForbidden
	}
	defer resp.Body.Close()

	// read the response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	r := &models.VerifyResponse{}
	err = json.Unmarshal(data, r)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	ctx.Set("login", r.Login)
	log.Info().Msg("Received user's email address, login: " + r.Email + " " + r.Login)
	return nil
}
