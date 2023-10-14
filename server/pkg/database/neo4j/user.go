package neo4j

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/docs/domain"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j/internal"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type userCreator struct {
	createCypher string
}

type userGetter struct {
	getByUsernameCypher string
}

type userUpdater struct {
	updateUsernameCypher string
	updatePasswordCypher string
}

type userDeleter struct {
	deleteCypher string
}

func NewUserCreator() *userCreator {
	return &userCreator{
		createCypher: `CREATE (u:User {username: $username, password: $password})`,
	}
}

func NewUserGetter() *userGetter {
	return &userGetter{
		getByUsernameCypher: `MATCH (u:User {username: $username}) RETURN u.username as username, u.password as password`,
	}
}

func NewUserUpdater() *userUpdater {
	return &userUpdater{
		updateUsernameCypher: `MATCH (u:User {username: $username}) SET u.username = $new_username RETURN COUNT(u) as count`,
		updatePasswordCypher: `MATCH (u:User {username: $username}) SET u.password = $new_password RETURN COUNT(u) as count`,
	}
}

func NewUserDeleter() *userDeleter {
	return &userDeleter{
		deleteCypher: `MATCH (u:User {username: $username}) DELETE u RETURN COUNT(u) as count`,
	}
}

func (c userCreator) Create(ctx context.Context, user domain.User) error {
	sessions := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer sessions.Close(ctx)

	params := map[string]any{
		"username": user.Username,
		"password": user.Password,
	}

	_, err := sessions.Run(ctx, c.createCypher, params)
	if err != nil && err.(*neo4j.Neo4jError).Code == ConstraintValidationFailed {
		return fmt.Errorf("username already taken")
	} else if err != nil {
		return fmt.Errorf("failed to store user in the database")
	}
	return nil
}

func (g userGetter) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"username": username,
	}

	result, err := session.Run(ctx, g.getByUsernameCypher, params)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get the user from the database")
	}

	return internal.GetSingle[domain.User](ctx, result)
}

func (u userUpdater) UpdateUsername(ctx context.Context, user domain.User, newUsername string) error {
	sessions := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer sessions.Close(ctx)

	params := map[string]any{
		"username":     user.Username,
		"new_username": newUsername,
	}

	result, err := sessions.Run(ctx, u.updateUsernameCypher, params)
	if err != nil && err.(*neo4j.Neo4jError).Code == ConstraintValidationFailed {
		return fmt.Errorf("username already taken")
	} else if err != nil {
		return fmt.Errorf("failed to update username")
	}

	count, err := internal.GetSingle[int64](ctx, result)
	if count <= 0 || err != nil {
		return fmt.Errorf("user wasn't found")
	}

	return nil
}

func (u userUpdater) UpdatePassword(ctx context.Context, user domain.User, newPassword string) error {
	sessions := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer sessions.Close(ctx)

	params := map[string]any{
		"username":     user.Username,
		"new_password": newPassword,
	}

	result, err := sessions.Run(ctx, u.updatePasswordCypher, params)
	if err != nil {
		return fmt.Errorf("failed to update password")
	}

	count, err := internal.GetSingle[int64](ctx, result)
	if count <= 0 || err != nil {
		return fmt.Errorf("user wasn't found")
	}

	return nil
}

func (d userDeleter) Delete(ctx context.Context, user domain.User) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"username": user.Username,
	}

	_, err := session.Run(ctx, d.deleteCypher, params)
	if err != nil {
		return fmt.Errorf("failed to delete the user")
	}
	return nil
}
