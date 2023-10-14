package neo4j

import (
	"context"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	driver neo4j.DriverWithContext

	EquivalentSchemaRuleAlreadyExists = "Neo.ClientError.Schema.EquivalentSchemaRuleAlreadyExists"
	ConstraintValidationFailed        = "Neo.ClientError.Schema.ConstraintValidationFailed"
)

func MustInitialize() {
	var (
		err error

		dsn      = os.Getenv("NEO4J_DSN")
		username = os.Getenv("NEO4J_USERNAME")
		password = os.Getenv("NEO4J_PASSWORD")
		realm    = os.Getenv("NEO4J_REALM")
	)

	driver, err = neo4j.NewDriverWithContext(dsn, neo4j.BasicAuth(username, password, realm))
	if err != nil {
		log.Fatal(err)
	}

	if err := driver.VerifyConnectivity(context.Background()); err != nil {
		log.Fatal(err)
	}

	defineConstraints()
}

func defineConstraints() {
	constraints := []string{
		`CREATE CONSTRAINT constraint_user_name_unique FOR (u:User) REQUIRE u.username IS UNIQUE`,
		`CREATE CONSTRAINT constraint_file_id_unique FOR (f:File) REQUIRE f.id IS UNIQUE`,
	}

	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	for _, constraint := range constraints {
		_, err := session.Run(ctx, constraint, nil)
		if err != nil && err.(*neo4j.Neo4jError).Code != EquivalentSchemaRuleAlreadyExists {
			log.Fatal(err)
		}
	}
}
