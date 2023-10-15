package neo4j

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/docs/domain"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j/internal"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type fileCreator struct {
	createCypher string
}

type fileGetter struct {
	getByIdCypher        string
	getOwnerCypher       string
	getAllForOwnerCypher string
}

type fileUpdater struct {
	updateNameCypher string
}

type fileDeleter struct {
	deleteCypher            string
	deleteAllForOwnerCypher string
}

func NewFileCreator() *fileCreator {
	return &fileCreator{
		createCypher: `MATCH (u:User {username: $username}) CREATE (u)-[:OWNS]->(f:File {id: $id, name: $name})`,
	}
}

func NewFileGetter() *fileGetter {
	return &fileGetter{
		getByIdCypher:        `MATCH (f:File {id: $id}) RETURN f`,
		getOwnerCypher:       `MATCH (u:User)-[:OWNS]->(f:File {id: $id}) RETURN u`,
		getAllForOwnerCypher: `MATCH (u:User {username: $username})-[:OWNS]->(f:File) RETURN f`,
	}
}

func NewFileUpdater() *fileUpdater {
	return &fileUpdater{
		updateNameCypher: `MATCH (f:File {id: $id}) SET f.name = $new_name RETURN COUNT(f) as c`,
	}
}

func NewFileDeleter() *fileDeleter {
	return &fileDeleter{
		deleteCypher:            `MATCH (f:File {id: $id}) DETACH DELETE f`,
		deleteAllForOwnerCypher: `MATCH (u:User {username: $username})-[:OWNS]->(f:File) DETACH DELETE f`,
	}
}

func (fc fileCreator) Create(ctx context.Context, file domain.File, owner domain.User) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"username": owner.Username,
		"id":       file.Id,
		"name":     file.Name,
	}

	_, err := session.Run(ctx, fc.createCypher, params)
	if err != nil {
		if neo4jErr, ok := err.(*neo4j.Neo4jError); ok && neo4jErr.Code == ConstraintValidationFailed {
			return fmt.Errorf("file with this id already exists")
		} else {
			return fmt.Errorf("failed to store the file in the database")
		}
	}

	return nil
}

func (fg fileGetter) GetById(ctx context.Context, id string) (domain.File, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"id": id,
	}

	result, err := session.Run(ctx, fg.getByIdCypher, params)
	if err != nil {
		return domain.File{}, fmt.Errorf("file to get the file from the database")
	}

	return internal.GetSingle[domain.File](ctx, result, "f")
}

func (fg fileGetter) GetOwner(ctx context.Context, file domain.File) (domain.User, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"id": file.Id,
	}

	result, err := session.Run(ctx, fg.getOwnerCypher, params)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get the file's owner from the database")
	}

	return internal.GetSingle[domain.User](ctx, result, "u")
}

func (fg fileGetter) GetAllForOwner(ctx context.Context, owner domain.User) ([]domain.File, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"username": owner.Username,
	}

	result, err := session.Run(ctx, fg.getAllForOwnerCypher, params)
	if err != nil {
		return []domain.File{}, fmt.Errorf("failed to get all files for owner from the database")
	}

	return internal.GetMultiple[domain.File](ctx, result, "f")
}

func (fu fileUpdater) UpdateName(ctx context.Context, file domain.File, name string) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"id":       file.Id,
		"new_name": name,
	}

	result, err := session.Run(ctx, fu.updateNameCypher, params)
	if err != nil {
		return fmt.Errorf("failed to update file's name")
	}

	if count, err := internal.GetSingle[int64](ctx, result, "c"); count <= 0 || err != nil {
		return fmt.Errorf("file wasn't found")
	}

	return nil
}

func (fd fileDeleter) Delete(ctx context.Context, file domain.File) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"id": file.Id,
	}

	if _, err := session.Run(ctx, fd.deleteCypher, params); err != nil {
		return fmt.Errorf("failed to delete the file")
	}

	return nil
}

func (fd fileDeleter) DeleteAllForOwner(ctx context.Context, owner domain.User) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"username": owner.Username,
	}

	if _, err := session.Run(ctx, fd.deleteAllForOwnerCypher, params); err != nil {
		return fmt.Errorf("failed to delete all files for owner")
	}

	return nil
}
