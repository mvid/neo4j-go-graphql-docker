A small project using Neo4J and Go, meant to be an example for a few different features.

* A Go based proxy in front of the Neo4J GraphQL endpoint, allowing for authorization/authentication against the GraphQL backend
* The Neo4J Go driver in Docker 

Special thanks to:
* [@jensneuse](https://github.com/jensneuse/) for designing and developing the [graphql tools](https://github.com/jensneuse/graphql-go-tools) used for the proxy
* [Neo4j community](https://github.com/neo4j-graphql) for developing the [GraphQL plugin](https://github.com/neo4j-graphql/neo4j-graphql)
* [Neo4j company](neo4j.com) for implementing a [Go driver](https://github.com/neo4j/neo4j-go-driver)