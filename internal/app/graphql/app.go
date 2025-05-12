package graphql

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/IlianBuh/GraphQL/internal/domain/models"
	"github.com/IlianBuh/GraphQL/internal/graph"
	n "github.com/IlianBuh/GraphQL/internal/lib/net"
	e "github.com/IlianBuh/GraphQL/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"
)

type GraphQLApp struct {
	log  *slog.Logger
	port int
	host string

	client  SSOApi
	server  *http.Server
	timeout time.Duration

	wg sync.WaitGroup
}

type SSOApi interface {
	SignUp(
		ctx context.Context,
		login string,
		email string,
		password string,
	) (token string, _ error)
	LogIn(
		ctx context.Context,
		login string,
		password string,
	) (token string, _ error)
	FollowersList(
		ctx context.Context,
		userID int32,
	) ([]*models.User, error)
	FolloweesList(
		ctx context.Context,
		userID int32,
	) ([]*models.User, error)
	Follow(
		ctx context.Context,
		srcId int,
		targetId int,
	) error
	Unfollow(
		ctx context.Context,
		srcId int,
		targetId int,
	) error

	User(
		ctx context.Context,
		uuid int,
	) (*models.User, error)
	Users(
		ctx context.Context,
		uuid []int,
	) ([]*models.User, error)
	UsersExist(
		ctx context.Context,
		uuid []int,
	) (bool, error)
	UsersByLogin(
		ctx context.Context,
		login string,
	) ([]*models.User, error)
}

func New(
	log *slog.Logger,
	port int,
	host string,
	timeout time.Duration,
	client SSOApi,
) *GraphQLApp {
	srv := newServer(host, port, timeout, client)

	return &GraphQLApp{
		log:     log,
		port:    port,
		host:    host,
		timeout: timeout,
		client:  client,
		server:  srv,
	}
}

func newServer(
	host string,
	port int,
	timeout time.Duration,
	client SSOApi,
) *http.Server {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{SSO: client}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	serverMux := http.NewServeMux()
	serverMux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	serverMux.Handle("/query", srv)

	server := &http.Server{
		Handler:      serverMux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  timeout,
		Addr:         n.Join(host, port),
	}

	return server
}

func (g *GraphQLApp) Start() {
	const op = "graphql-app.Start"
	log := g.log.With(slog.String("op", op))
	log.Info("starting graphql application")

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		if err := g.server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
		}
	}()

	log.Info("graphql application started")
}

func (g *GraphQLApp) Stop() error {
	const op = "graphql-app.Stop"
	log := g.log.With(slog.String("op", op))
	log.Info("stopping graphql application")

	ctx, cncl := context.WithTimeout(context.Background(), g.timeout)
	defer cncl()

	err := g.server.Shutdown(ctx)
	if err != nil {
		return e.Fail(op, err)
	}
	g.wg.Wait()

	log.Info("graphql application stopped")
	return nil
}
