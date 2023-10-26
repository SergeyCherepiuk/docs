package neo4j

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/docs/pkg/database/models"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j/internal"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type userService struct {
	createCypher string

	getByUsernameCypher string

	updateUsernameCypher string
	updatePasswordCypher string

	deleteCypher string
}

func NewUserService() *userService {
	return &userService{
		createCypher: `CREATE (u:User {username: $username, password: $password})`,

		getByUsernameCypher: `MATCH (u:User {username: $username}) RETURN u`,

		updateUsernameCypher: `MATCH (u:User {username: $username}) SET u.username = $new_username RETURN COUNT(u) as c`,
		updatePasswordCypher: `MATCH (u:User {username: $username}) SET u.password = $new_password RETURN COUNT(u) as c`,

		deleteCypher: `MATCH (u:User {username: $username}) OPTIONAL MATCH (u)-[r:OWNS]->(f:File) DETACH DELETE u, r, f`,
	}
}

var UserService = NewUserService()

func (s userService) Create(ctx context.Context, runner runner, user models.User) error {
	params := map[string]any{
		"username": user.Username,
		"password": user.Password,
	}

	_, err := runner.Run(ctx, s.createCypher, params)
	if err != nil {
		if neo4jErr, ok := err.(*neo4j.Neo4jError); ok && neo4jErr.Code == ConstraintValidationFailed {
			return fmt.Errorf("username already taken")
		} else {
			return fmt.Errorf("failed to store user in the database")
		}
	}
	return nil
}

func (s userService) GetByUsername(ctx context.Context, runner runner, username string) (models.User, error) {
	params := map[string]any{
		"username": username,
	}

	result, err := runner.Run(ctx, s.getByUsernameCypher, params)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get the user from the database")
	}

	user, err := internal.GetSingle[models.User](ctx, result, "u")
	if err != nil {
		switch err.(type) {
		case internal.ErrorNoRecords, internal.ErrorNilRecord:
			return models.User{}, fmt.Errorf("user wasn't found")
		default:
			return models.User{}, fmt.Errorf("failed to get the user from the database")
		}
	}

	return user, nil
}

func (s userService) UpdateUsername(ctx context.Context, runner runner, user models.User, newUsername string) error {
	params := map[string]any{
		"username":     user.Username,
		"new_username": newUsername,
	}

	result, err := runner.Run(ctx, s.updateUsernameCypher, params)
	if err != nil {
		if neo4jErr, ok := err.(*neo4j.Neo4jError); ok && neo4jErr.Code == ConstraintValidationFailed {
			return fmt.Errorf("username already taken")
		} else {
			return fmt.Errorf("failed to update user's username")
		}
	}

	count, err := internal.GetSingle[int64](ctx, result, "c")
	if err != nil {
		switch err.(type) {
		case internal.ErrorNoRecords, internal.ErrorNilRecord:
			return fmt.Errorf("user wasn't found")
		default:
			return fmt.Errorf("failed to update user's username")
		}
	}

	if count <= 0 {
		return fmt.Errorf("failed to update user's username")
	}

	return nil
}

func (s userService) UpdatePassword(ctx context.Context, runner runner, user models.User, newPassword string) error {
	params := map[string]any{
		"username":     user.Username,
		"new_password": newPassword,
	}

	result, err := runner.Run(ctx, s.updatePasswordCypher, params)
	if err != nil {
		return fmt.Errorf("failed to update user's password")
	}

	count, err := internal.GetSingle[int64](ctx, result, "c")
	if err != nil {
		switch err.(type) {
		case internal.ErrorNoRecords, internal.ErrorNilRecord:
			return fmt.Errorf("user wasn't found")
		default:
			return fmt.Errorf("failed to update user's password")
		}
	}

	if count <= 0 {
		return fmt.Errorf("failed to update user's password")
	}

	return nil
}

func (s userService) Delete(ctx context.Context, runner runner, user models.User) error {
	params := map[string]any{
		"username": user.Username,
	}

	if _, err := runner.Run(ctx, s.deleteCypher, params); err != nil {
		return fmt.Errorf("failed to delete the user")
	}
	return nil
}
