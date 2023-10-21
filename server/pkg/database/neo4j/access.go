package neo4j

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/docs/pkg/database/models"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j/internal"
)

type accessService struct {
	grantReadCypher      string
	grantReadWriteCypher string

	getCypher          string
	getAccessorsCypher string

	updateLevelCypher string

	revokeCypher string
}

func NewAccessService() *accessService {
	return &accessService{
		grantReadCypher:      `MATCH (u:User {username: $receiver}), (f:File {id: $id}) CREATE (u)-[:CAN_ACCESS {level: "R", grantedBy: $granter}]->(f)`,
		grantReadWriteCypher: `MATCH (u:User {username: $receiver}), (f:File {id: $id}) CREATE (u)-[:CAN_ACCESS {level: "RW", grantedBy: $granter}]->(f)`,

		getCypher:          `MATCH (u:User {username: $username})-[a:CAN_ACCESS]->(f:File {id: $id}) RETURN {granter: a.grantedBy, receiver: u.username, level: a.level} as a`,
		getAccessorsCypher: `MATCH (u:User)-[a:CAN_ACCESS]->(f:File {id: $id}) RETURN {granter: a.grantedBy, receiver: u.username, level: a.level} as a`,

		updateLevelCypher: `MATCH (u:User {username: $receiver})-[a:CAN_ACCESS {grantedBy: $granter}]->(f:File {id: $id}) SET a.level = $new_level RETURN COUNT(a) as c`,

		revokeCypher: `MATCH (u:User {username: $receiver})-[a:CAN_ACCESS {grantedBy: $granter}]->(f:File {id: $id}) DELETE a`,
	}
}

var AccessService = NewAccessService()

func (s accessService) Grant(ctx context.Context, runner runner, file models.File, access models.Access) error {
	var cypher string
	switch access.Level {
	case models.RWAccess:
		cypher = s.grantReadWriteCypher
	case models.RAcess:
		cypher = s.grantReadCypher
	default:
		return fmt.Errorf("unknown access level value: %s", access.Level)
	}

	params := map[string]any{
		"receiver": access.Receiver,
		"id":       file.Id,
		"granter":  access.Granter,
	}

	if _, err := runner.Run(ctx, cypher, params); err != nil {
		return fmt.Errorf("failed to grand an access")
	}

	return nil
}

func (s accessService) Get(ctx context.Context, runner runner, file models.File, user models.User) (models.Access, error) {
	params := map[string]any{
		"username": user.Username,
		"id":       file.Id,
	}

	result, err := runner.Run(ctx, s.getCypher, params)
	if err != nil {
		return models.Access{}, fmt.Errorf("failed to get get the file access")
	}

	return internal.GetSingle[models.Access](ctx, result, "a")
}

func (s accessService) GetAccesses(ctx context.Context, runner runner, file models.File) ([]models.Access, error) {
	params := map[string]any{
		"id": file.Id,
	}

	result, err := runner.Run(ctx, s.getAccessorsCypher, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get accessors")
	}

	return internal.GetMultiple[models.Access](ctx, result, "a")
}

func (s accessService) UpdateLevel(ctx context.Context, runner runner, file models.File, access models.Access, newLevel string) error {
	params := map[string]any{
		"receiver":  access.Receiver,
		"granter":   access.Granter,
		"id":        file.Id,
		"new_level": newLevel,
	}

	result, err := runner.Run(ctx, s.updateLevelCypher, params)
	if err != nil {
		return fmt.Errorf("failed to update access level")
	}

	if count, err := internal.GetSingle[int64](ctx, result, "c"); count <= 0 || err != nil {
		return fmt.Errorf("access record wasn't found")
	}

	return nil
}

func (s accessService) Revoke(ctx context.Context, runner runner, file models.File, access models.Access) error {
	params := map[string]any{
		"receiver": access.Receiver,
		"granter":  access.Granter,
		"id":       file.Id,
	}

	if _, err := runner.Run(ctx, s.revokeCypher, params); err != nil {
		return fmt.Errorf("failed to revoke the access")
	}

	return nil
}
