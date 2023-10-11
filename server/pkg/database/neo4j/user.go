package neo4j

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/docs/domain"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type userCreator struct {
	createCypher string
}

func NewUserCreator() *userCreator {
	return &userCreator{
		createCypher: `CREATE (user:User {username: $username, password: $password})`,
	}
}

func (c userCreator) Create(ctx context.Context, user domain.User) error {
	params := map[string]any{
		"username": user.Username,
		"password": user.Password,
	}

	sessions := driver.NewSession(ctx, neo4j.SessionConfig{})
	_, err := sessions.Run(ctx, c.createCypher, params)

	if err != nil && err.(*neo4j.Neo4jError).Code == ConstraintValidationFailed {
		return fmt.Errorf("username already taken")
	} else if err != nil {
		return fmt.Errorf("failed to store user in the database")
	}
	return nil
}
