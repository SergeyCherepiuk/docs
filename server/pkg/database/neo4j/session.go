package neo4j

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/docs/pkg/database/models"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j/internal"
	"github.com/google/uuid"
)

type sessionService struct {
	createCypher string

	checkCypher string
}

func NewSessionService() *sessionService {
	return &sessionService{
		createCypher: `MATCH (u:User {username: $username}) CREATE (u)-[:HAS]->(s:Session {id: $id, created_at: $created_at, expires_at: $expires_at})`,

		checkCypher: `MATCH (s:Session {id: $id}) RETURN s.expires_at > datetime() as a`,
	}
}

var SessionService = NewSessionService()

func (s sessionService) Create(ctx context.Context, runner runner, session models.Session) error {
	params := map[string]any{
		"id":         session.Id,
		"username":   session.Username,
		"created_at": session.CreatedAt,
		"expires_at": session.ExpiresAt,
	}

	if _, err := runner.Run(ctx, s.createCypher, params); err != nil {
		return fmt.Errorf("failed to create a session")
	}

	return nil
}

// TODO: Return session's owner (models.User)
func (s sessionService) Check(ctx context.Context, runner runner, id uuid.UUID) (bool, error) {
	params := map[string]any{
		"id": id.String(),
	}

	result, err := runner.Run(ctx, s.checkCypher, params)
	if err != nil {
		return false, fmt.Errorf("failed to check the session")
	}

	return internal.GetSingle[bool](ctx, result, "a")
}
