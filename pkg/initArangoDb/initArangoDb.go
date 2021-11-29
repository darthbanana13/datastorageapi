package initArangoDb

import (
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func InitDbConn(connStr, user, pass string) (driver.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{connStr},
	})
	if err != nil {
		return nil, err
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(user, pass),
	})
	if err != nil {
		return nil, err
	}
	return client, err
}

//TODO: Add a trace log in case logrus is present
func InitDb(client driver.Client, dbName string) (driver.Database, error) {
	var db driver.Database

	dbExists, err := client.DatabaseExists(nil, dbName)
	if dbExists {
		// fmt.Println("That db exists already")
		db, err = client.Database(nil, dbName)
		if err != nil {
			return nil, fmt.Errorf("Failed to open existing database: %v", err)
		}
	} else {
		db, err = client.CreateDatabase(nil, dbName, nil)
		if err != nil {
			return nil, fmt.Errorf("Failed to create database: %v", err)
		}
	}
	return db, nil
}

func InitColl(db driver.Database, collName string) error {
	collExists, err := db.CollectionExists(nil, collName)
	if !collExists {
		_, err = db.CreateCollection(nil, collName, nil)

		if err != nil {
			return fmt.Errorf("Failed to create collection: %v", err)
		}
	}
	return nil
}

func InitDbWith(connStr, user, pass, dbName string, colls []string) (driver.Database, error) {
	client, err := InitDbConn(connStr, user, pass)
	if err != nil {
		return nil, err
	}

	db, err := InitDb(client, dbName)
	if err != nil {
		return nil, err
	}

	for _, coll := range colls {
		err = InitColl(db, coll)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
