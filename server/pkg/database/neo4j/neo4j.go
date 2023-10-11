package neo4j

import (
	"context"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var ndb neo4j.DriverWithContext

func MustInitialize() {
	var (
		err error

		dsn      = os.Getenv("NEO4J_DSN")
		username = os.Getenv("NEO4J_USERNAME")
		password = os.Getenv("NEO4J_PASSWORD")
		realm    = os.Getenv("NEO4J_REALM")
	)

	ndb, err = neo4j.NewDriverWithContext(dsn, neo4j.BasicAuth(username, password, realm))
	if err != nil {
		log.Fatal(err)
	}

	if err := ndb.VerifyConnectivity(context.Background()); err != nil {
		log.Fatal(err)
	}
}
