package neo4j

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/docs/pkg/database/models"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j/internal"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type fileService struct {
	createCypher string

	getByIdCypher        string
	getOwnerCypher       string
	getAllForOwnerCypher string

	updateNameCypher string

	deleteCypher            string
	deleteAllForOwnerCypher string
}

func NewFileService() *fileService {
	return &fileService{
		createCypher: `MATCH (u:User {username: $username}) CREATE (u)-[:OWNS]->(f:File {id: $id, name: $name})`,

		getByIdCypher:        `MATCH (f:File {id: $id}) RETURN f`,
		getOwnerCypher:       `MATCH (u:User)-[:OWNS]->(f:File {id: $id}) RETURN u`,
		getAllForOwnerCypher: `MATCH (u:User {username: $username})-[:OWNS]->(f:File) RETURN f`,

		updateNameCypher: `MATCH (f:File {id: $id}) SET f.name = $new_name RETURN COUNT(f) as c`,

		deleteCypher:            `MATCH (f:File {id: $id}) DETACH DELETE f`,
		deleteAllForOwnerCypher: `MATCH (u:User {username: $username})-[:OWNS]->(f:File) DETACH DELETE f`,
	}
}

var FileService = NewFileService()

func (s fileService) Create(ctx context.Context, runner runner, file models.File, owner models.User) error {
	params := map[string]any{
		"username": owner.Username,
		"id":       file.Id,
		"name":     file.Name,
	}

	_, err := runner.Run(ctx, s.createCypher, params)
	if err != nil {
		if neo4jErr, ok := err.(*neo4j.Neo4jError); ok && neo4jErr.Code == ConstraintValidationFailed {
			return fmt.Errorf("file with this id already exists")
		} else {
			return fmt.Errorf("failed to store the file in the database")
		}
	}

	return nil
}

func (s fileService) GetById(ctx context.Context, runner runner, id uuid.UUID) (models.File, error) {
	params := map[string]any{
		"id": id.String(),
	}

	result, err := runner.Run(ctx, s.getByIdCypher, params)
	if err != nil {
		return models.File{}, fmt.Errorf("file to get the file from the database")
	}

	return internal.GetSingle[models.File](ctx, result, "f")
}

func (s fileService) GetOwner(ctx context.Context, runner runner, file models.File) (models.User, error) {
	params := map[string]any{
		"id": file.Id,
	}

	result, err := runner.Run(ctx, s.getOwnerCypher, params)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get the file's owner from the database")
	}

	return internal.GetSingle[models.User](ctx, result, "u")
}

func (s fileService) GetAllForOwner(ctx context.Context, runner runner, owner models.User) ([]models.File, error) {
	params := map[string]any{
		"username": owner.Username,
	}

	result, err := runner.Run(ctx, s.getAllForOwnerCypher, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get all files for owner from the database")
	}

	return internal.GetMultiple[models.File](ctx, result, "f")
}

func (s fileService) UpdateName(ctx context.Context, runner runner, file models.File, name string) error {
	params := map[string]any{
		"id":       file.Id,
		"new_name": name,
	}

	result, err := runner.Run(ctx, s.updateNameCypher, params)
	if err != nil {
		return fmt.Errorf("failed to update file's name")
	}

	if count, err := internal.GetSingle[int64](ctx, result, "c"); count <= 0 || err != nil {
		return fmt.Errorf("file wasn't found")
	}

	return nil
}

func (s fileService) Delete(ctx context.Context, runner runner, file models.File) error {
	params := map[string]any{
		"id": file.Id,
	}

	if _, err := runner.Run(ctx, s.deleteCypher, params); err != nil {
		return fmt.Errorf("failed to delete the file")
	}

	return nil
}

func (s fileService) DeleteAllForOwner(ctx context.Context, runner runner, owner models.User) error {
	params := map[string]any{
		"username": owner.Username,
	}

	if _, err := runner.Run(ctx, s.deleteAllForOwnerCypher, params); err != nil {
		return fmt.Errorf("failed to delete all files for owner")
	}

	return nil
}
