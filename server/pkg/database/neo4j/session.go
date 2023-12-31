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

	deleteAllCypher string
}

func NewSessionService() *sessionService {
	return &sessionService{
		createCypher: `MATCH (u:User {username: $username}) CREATE (u)-[:HAS]->(s:Session {id: $id, created_at: $created_at, expires_at: $expires_at})`,

		checkCypher: `MATCH (u:User)-[:HAS]->(s:Session {id: $id}) WHERE s.expires_at > datetime() RETURN u`,

		deleteAllCypher: `MATCH (u:User {username: $username})-[:HAS]->(s:Session) DETACH DELETE s`,
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

func (s sessionService) Check(ctx context.Context, runner runner, id uuid.UUID) (models.User, error) {
	params := map[string]any{
		"id": id.String(),
	}

	result, err := runner.Run(ctx, s.checkCypher, params)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to check the session")
	}

	user, err := internal.GetSingle[models.User](ctx, result, "u")
	if err != nil {
		switch err.(type) {
		case internal.ErrorNoRecords, internal.ErrorNilRecord:
			return models.User{}, fmt.Errorf("user wasn't found")
		default:
			return models.User{}, fmt.Errorf("failed to check the session")
		}
	}

	return user, nil
}

func (s sessionService) DeleteAll(ctx context.Context, runner runner, user models.User) error {
	params := map[string]any{
		"username": user.Username,
	}

	if _, err := runner.Run(ctx, s.deleteAllCypher, params); err != nil {
		return fmt.Errorf("failed to delete all sessions")
	}

	return nil
}
