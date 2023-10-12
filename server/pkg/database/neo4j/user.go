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
		createCypher: `CREATE (u:User {username: $username, password: $password})`,
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

type userGetter struct {
	getUserByUsernameCypher string
}

func NewUserGetter() *userGetter {
	return &userGetter{
		getUserByUsernameCypher: `MATCH (u:User {username: $username}) RETURN u.password as password`,
	}
}

func (g userGetter) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"username": username,
	}

	result, err := session.Run(ctx, g.getUserByUsernameCypher, params)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get the user from the database")
	}

	record, err := result.Single(ctx)
	if record == nil || err != nil {
		return domain.User{}, fmt.Errorf("user wasn't found")
	}

	password, passwordFound := record.Get("password")
	if !passwordFound {
		return domain.User{}, fmt.Errorf("failed to get the user from the database")
	}

	return domain.User{
		Username: username,
		Password: password.(string),
	}, nil
}

type userUpdater struct {
	updateUsernameCypher string
	updatePasswordCypher string
}

func NewUserUpdater() *userUpdater {
	return &userUpdater{
		updateUsernameCypher: `MATCH (u:User {username: $username}) SET u.username = $new_username RETURN COUNT(u)`,
		updatePasswordCypher: `MATCH (u:User {username: $username}) SET u.password = $new_password RETURN COUNT(u)`,
	}
}

func (u userUpdater) UpdateUsername(ctx context.Context, username, newUsername string) error {
	sessions := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer sessions.Close(ctx)

	params := map[string]any{
		"username":     username,
		"new_username": newUsername,
	}

	// TODO: Use "result" variable to check if any node was updated
	_, err := sessions.Run(ctx, u.updateUsernameCypher, params)
	if err != nil && err.(*neo4j.Neo4jError).Code == ConstraintValidationFailed {
		return fmt.Errorf("username already taken")
	} else if err != nil {
		return fmt.Errorf("failed to update username")
	}
	return nil
}

func (u userUpdater) UpdatePassword(ctx context.Context, username, newPassword string) error {
	sessions := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer sessions.Close(ctx)

	params := map[string]any{
		"username":     username,
		"new_password": newPassword,
	}

	// TODO: Use "result" variable to check if any node was updated
	_, err := sessions.Run(ctx, u.updateUsernameCypher, params)
	if err != nil {
		return fmt.Errorf("failed to update password")
	}
	return nil
}

type userDeleter struct {
	deleteCypher string
}

func NewUserDeleter() *userDeleter {
	return &userDeleter{
		deleteCypher: `MATCH (u:User {username: $username}) DELETE u RETURN COUNT(u)`,
	}
}

func (d userDeleter) Delete(ctx context.Context, username string) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"username": username,
	}

	// TODO: Use "result" variable to check if any node was deleted
	_, err := session.Run(ctx, d.deleteCypher, params)
	if err != nil {
		return fmt.Errorf("failed to delete the user")
	}
	return nil
}
