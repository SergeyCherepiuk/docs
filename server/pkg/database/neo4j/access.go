package neo4j

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/docs/domain"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type accessGranter struct {
	grantReadCypher      string
	grantReadWriteCypher string
}

func NewAccessGranter() *accessGranter {
	return &accessGranter{
		grantReadCypher:      `MATCH (u:User {username: $receiver}), (f:File {id: $id}) CREATE (u)-[:CAN_ACCESS {level: "R", grantedBy: $granter}]->(f)`,
		grantReadWriteCypher: `MATCH (u:User {username: $receiver}), (f:File {id: $id}) CREATE (u)-[:CAN_ACCESS {level: "RW", grantedBy: $granter}]->(f)`,
	}
}

func (ag accessGranter) Grant(ctx context.Context, file domain.File, access domain.Access) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	var cypher string
	if access.Level == domain.RWAccess {
		cypher = ag.grantReadWriteCypher
	} else if access.Level == domain.RAcess {
		cypher = ag.grantReadCypher
	} else {
		return fmt.Errorf("unknown access level value: %d", access.Level)
	}

	params := map[string]any{
		"receiver": access.Receiver,
		"id":       file.Id,
		"granter":  access.Granter,
	}

	if _, err := session.Run(ctx, cypher, params); err != nil {
		return fmt.Errorf("failed to grand an access")
	}

	return nil
}
