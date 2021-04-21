package storage

import (
	"fmt"
	"github.com/gocql/gocql"
	"os"
)

type SistStorage struct {
	Session *gocql.Session
}

func NewSistStorage() (*SistStorage, error){
	cluster := gocql.NewCluster(fmt.Sprintf("%s:%s", os.Getenv("CASSANDRA_HOSTNAME"), os.Getenv("CASSANDRA_HOSTNAME_PORT")))
	cluster.Keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &SistStorage{session}, nil
}
