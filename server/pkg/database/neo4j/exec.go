package neo4j

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type runner interface {
	Run(ctx context.Context, cypher string, params map[string]any) (neo4j.ResultWithContext, error)
}

type Session struct {
	neo4j.SessionWithContext
}

// NOTE: "Overriding" neo4j.SessionWithContext's Run method,
// to use it with default transaction parameters
func (s Session) Run(ctx context.Context, cypher string, params map[string]any) (neo4j.ResultWithContext, error) {
	return s.SessionWithContext.Run(ctx, cypher, params)
}

func NewSession(ctx context.Context) Session {
	return Session{
		SessionWithContext: driver.NewSession(ctx, neo4j.SessionConfig{}),
	}
}
