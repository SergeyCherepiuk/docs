package neo4j

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/docs/pkg/database/models"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type sessionService struct {
	createCypher string
}

func NewSessionService() *sessionService {
	return &sessionService{
		createCypher: `MATCH (u:User {username: $username}) CREATE (u)-[:HAS]->(s:Session {id: $id, created_at: apoc.date.currentTimestamp(), expires_at: apoc.date.currentTimestamp() + $expires_in})`,
	}
}

var SessionService = NewSessionService()

func (s sessionService) Create(ctx context.Context, session models.Session) error {
	dbSession := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer dbSession.Close(ctx)

	params := map[string]any{
		"username":   session.User.Username,
		"id":         uuid.NewString(),
		"expires_in": 7 * 24 * 60 * 60,
	}

	if _, err := dbSession.Run(ctx, s.createCypher, params); err != nil {
		return fmt.Errorf("failed to create a session")
	}

	return nil
}
