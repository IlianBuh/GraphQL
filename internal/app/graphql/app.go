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
	"github.com/IlianBuh/GraphQL/internal/graph"
	n "github.com/IlianBuh/GraphQL/internal/lib/net"
	e "github.com/IlianBuh/GraphQL/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"
)

type GraphQLApp struct {
	log     *slog.Logger
	port    int
	host    string
	server  http.Server
	timeout time.Duration
	wg      sync.WaitGroup
}

func New(
	log *slog.Logger,
	port int,
	host string,
) *GraphQLApp {
	return &GraphQLApp{
		log:  log,
		port: port,
		host: host,
	}
}

func (g *GraphQLApp) Start() {
	const op = "graphql-app.Start"
	log := g.log.With(slog.String("op", op))
	log.Info("starting graphql application")

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

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

		g.server = http.Server{
			Handler:      serverMux,
			ReadTimeout:  g.timeout,
			WriteTimeout: g.timeout,
			IdleTimeout:  g.timeout,
			Addr:         n.Join(g.host, g.port),
		}

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
