package app

import (
	"log/slog"
	"time"

	"github.com/IlianBuh/GraphQL/internal/app/graphql"
	sso "github.com/IlianBuh/GraphQL/internal/clients/sso/grpc"
	"github.com/IlianBuh/GraphQL/internal/config"
	"github.com/IlianBuh/GraphQL/internal/lib/sl"
	e "github.com/IlianBuh/GraphQL/pkg/errors"
)

const emptyHost = ""

type App struct {
	log       *slog.Logger
	GraphQL   *graphql.GraphQLApp
	SSOClient *sso.Client
}

func New(
	log *slog.Logger,
	graphqlPort int,
	timeout time.Duration,
	ssoConfig config.SSOConfig,
) *App {
	client, err := sso.NewClient(log, ssoConfig.Port, ssoConfig.Host)
	if err != nil {
		panic(err)
	}

	graphqlapp := graphql.New(
		log,
		graphqlPort,
		emptyHost,
		timeout,
		client,
	)

	return &App{
		log:       log,
		GraphQL:   graphqlapp,
		SSOClient: client,
	}
}

func (a *App) Stop() error {
	const op = "app.Stop"
	log := a.log.With(slog.String("op", op))
	log.Info("stopping application")

	err := a.GraphQL.Stop()
	if err != nil {
		log.Error("failed to stop graphql", sl.Err(err))

		return e.Fail(op, err)
	}

	err = a.SSOClient.Stop()
	if err != nil {
		log.Error("failed to stop sso client", sl.Err(err))

		return e.Fail(op, err)
	}

	return nil
}
