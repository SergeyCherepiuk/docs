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
	getByIdCypher string
}

type fileUpdater struct {
	updateNameCypher string
}

func NewFileCreator() *fileCreator {
	return &fileCreator{
		createCypher: `CREATE (f:File {id: $id, name: $name})`,
	}
}

func NewFileGetter() *fileGetter {
	return &fileGetter{
		getByIdCypher: `MATCH (f:File {id: $id}) RETURN f.id as id, f.name as name`,
	}
}

func NewFileUpdater() *fileUpdater {
	return &fileUpdater{
		updateNameCypher: `MATCH (f:File {id: $id}) SET f.name = $new_name RETURN COUNT(f) as count`,
	}
}

func (fc fileCreator) Create(ctx context.Context, file domain.File) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
		"id":   file.Id,
		"name": file.Name,
	}

	_, err := session.Run(ctx, fc.createCypher, params)
	if err != nil {
		if neo4jErr, ok := err.(*neo4j.Neo4jError); ok && neo4jErr.Code == ConstraintValidationFailed {
			return fmt.Errorf("file with this id already exists")
		} else {
			fmt.Println(err.Error())
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

	return internal.GetSingle[domain.File](ctx, result)
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

	if count, err := internal.GetSingle[int64](ctx, result); count <= 0 || err != nil {
		return fmt.Errorf("file wasn't found")
	}

	return nil
}
