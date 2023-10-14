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
