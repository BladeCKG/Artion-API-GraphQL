// Package handlers holds HTTP/WS handlers chain along with separate middleware implementations.
package handlers

import (
	"artion-api-graphql/internal/config"
	"artion-api-graphql/internal/graphql/resolvers"
	"artion-api-graphql/internal/graphql/schema"
	"artion-api-graphql/internal/logger"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"github.com/rs/cors"
	"net/http"
)

// Api constructs and return the API HTTP handlers chain for serving GraphQL API calls.
func Api(cfg *config.Config, log logger.Logger, rs *resolvers.RootResolver) http.Handler {
	// Create new CORS handler and attach the logger into it so we get information on Debug level if needed
	corsHandler := cors.New(corsOptions(cfg))
	corsHandler.Log = log

	// create new parsed GraphQL schema
	sch, err := schema.Schema()
	if err != nil {
		log.Criticalf("can not get the GraphQL schema; %s", err)
		return nil
	}

	// return the constructed API handler chain
	return &LoggingHandler{
		logger:  log,
		handler: corsHandler.Handler(graphqlws.NewHandlerFunc(sch, &relay.Handler{Schema: sch})),
	}
}

// corsOptions constructs new set of options for the CORS handler based on provided configuration.
func corsOptions(cfg *config.Config) cors.Options {
	return cors.Options{
		AllowedOrigins: cfg.Server.CorsOrigin,
		AllowedMethods: []string{"HEAD", "GET", "POST"},
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With"},
		MaxAge:         300,
	}
}