package app

import (
	"log/slog"

	"github.com/IlianBuh/GraphQL/internal/app/graphql"
	"github.com/IlianBuh/GraphQL/internal/config"
)

const emptyHost = ""

type App struct {
	GraphQL *graphql.GraphQLApp
	// SSO     SSOApp
}

func New(
	log *slog.Logger,
	graphqlPort int,
	ssoConfig config.SSOConfig,
) *App {
	graphqlapp := graphql.New(log, graphqlPort, emptyHost) // TODO: init ssoApp

	// TODO: init graphQLapp

	return &App{
		GraphQL: graphqlapp,
	}
}

// TODO: implement common stop
func Stop() {
	panic("implement me")
}
