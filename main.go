package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

var auth = neo4j.BasicAuth("neo4j", "secretpassword", "")

func testNeoDriver() error {
	var (
		driver  neo4j.Driver
		session neo4j.Session
		result  neo4j.Result
		err     error
	)

	if driver, err = neo4j.NewDriver("bolt://neo4j:7687", auth); err != nil {
		return err // handle error
	}
	// handle driver lifetime based on your application lifetime requirements
	// driver's lifetime is usually bound by the application lifetime, which usually implies one driver instance per application
	defer driver.Close()

	if session, err = driver.Session(neo4j.AccessModeWrite); err != nil {
		return err
	}
	defer session.Close()

	result, err = session.Run("CREATE (n:Item { id: $id, name: $name }) RETURN n.id, n.name", map[string]interface{}{
		"id":   1,
		"name": "Item 1",
	})
	if err != nil {
		return err // handle error
	}

	for result.Next() {
		fmt.Printf("Created Item with Id = '%d' and Name = '%s'\n", result.Record().GetByIndex(0).(int64), result.Record().GetByIndex(1).(string))
	}
	if err = result.Err(); err != nil {
		return err // handle error
	}

	return nil
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	if _, err := rw.Write([]byte("index")); err != nil {
		panic(err)
	}
}

func main() {
	if err := testNeoDriver(); err != nil {
		panic(err)
	}

	fmt.Println("initializing graphql")
	if err := InitializeSchema(); err != nil {
		panic(err)
	}

	fmt.Println("starting proxy")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/graphql", Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
