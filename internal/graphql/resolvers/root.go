// Package resolvers implements GraphQL resolvers to incoming API requests.
package resolvers

import (
	"artion-api-graphql/cmd/artionapi/build"
	"artion-api-graphql/internal/config"
	"artion-api-graphql/internal/logger"
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
)

// RootResolver is GraphQL resolver of root namespace.
type RootResolver struct {
	wg      sync.WaitGroup
	cg      singleflight.Group
	sigStop chan bool
}

// log represents the logger to be used by the repository.
var log logger.Logger

// config represents the configuration setup used by the repository
// to establish and maintain required connectivity to external services
// as needed.
var cfg *config.Config

// SetLogger sets the repository logger to be used to collect logging info.
func SetLogger(l logger.Logger) {
	log = l
}

// SetConfig sets the repository configuration to be used to establish
// and maintain external repository connections.
func SetConfig(c *config.Config) {
	cfg = c
}

// New creates a new root resolver instance and initializes its internal structure.
func New() *RootResolver {
	if cfg == nil {
		panic(fmt.Errorf("missing configuration"))
	}
	if log == nil {
		panic(fmt.Errorf("missing logger"))
	}

	// create new resolver
	rs := RootResolver{
		// create terminator
		sigStop: make(chan bool, 1),
	}

	// handle broadcast and subscriptions in a separate routine
	rs.wg.Add(1)
	go rs.run()

	return &rs
}

// Close terminates resolver's broadcast service.
func (rs *RootResolver) Close() {
	// log
	log.Notice("GraphQL resolver is closing")

	// send the signal
	rs.sigStop <- true
	rs.wg.Wait()
}

// run monitors and handles subscriptions and broadcasts incoming events to their subscribers.
func (rs *RootResolver) run() {
	// sign off on leaving
	defer func() {
		// terminate
		log.Notice("GraphQL resolver done")
		rs.wg.Done()
	}()

	// log action
	log.Notice("GraphQL resolver started")

	// main loop waits for data on any channel and act upon it
	for {
		select {
		case <-rs.sigStop:
			return
		}
	}
}

// Version resolves the current version of the API server.
func (rs *RootResolver) Version() string {
	return build.Short(cfg)
}