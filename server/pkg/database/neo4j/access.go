package neo4j

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/docs/domain"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j/internal"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type accessGranter struct {
	grantReadCypher      string
	grantReadWriteCypher string
}

type accessGetter struct {
	getAccessorsCypher string
}

func NewAccessGranter() *accessGranter {
	return &accessGranter{
		grantReadCypher:      `MATCH (u:User {username: $receiver}), (f:File {id: $id}) CREATE (u)-[:CAN_ACCESS {level: "R", grantedBy: $granter}]->(f)`,
		grantReadWriteCypher: `MATCH (u:User {username: $receiver}), (f:File {id: $id}) CREATE (u)-[:CAN_ACCESS {level: "RW", grantedBy: $granter}]->(f)`,
	}
}

func NewAccessGetter() *accessGetter {
	return &accessGetter{
		getAccessorsCypher: `MATCH (u:User)-[a:CAN_ACCESS]->(f:File {id: $id}) RETURN {granter: a.grantedBy, receiver: u.username, level: a.level} as accesses`,
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
		return fmt.Errorf("unknown access level value: %s", access.Level)
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

func (ag accessGetter) GetAccesses(ctx context.Context, file domain.File) ([]domain.Access, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"id": file.Id,
	}

	result, err := session.Run(ctx, ag.getAccessorsCypher, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get accessors")
	}

	return internal.GetMultiple[domain.Access](ctx, result, "accesses")
}
