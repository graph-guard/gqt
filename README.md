<a href="https://github.com/graph-guard/gqt/actions?query=workflow%3ACI">
    <img src="https://github.com/graph-guard/gqt/workflows/CI/badge.svg" alt="GitHub Actions: CI">
</a>
<a href="https://coveralls.io/github/graph-guard/gqt">
    <img src="https://coveralls.io/repos/github/graph-guard/gqt/badge.svg" alt="Coverage Status" />
</a>
<a href="https://goreportcard.com/report/github.com/graph-guard/gqt">
    <img src="https://goreportcard.com/badge/github.com/graph-guard/gqt" alt="GoReportCard">
</a>
<a href="https://pkg.go.dev/github.com/graph-guard/gqt/v4">
    <img src="https://godoc.org/github.com/graph-guard/gqt/v4?status.svg" alt="GoDoc">
</a>

# GraphQL Query Template Language

This Go package provides a parser for **GQT**, the **G**raphQL **Q**uery **T**emplate language which is used by [Graph Guard](https://graphguard.io) for defining [GraphQL](https://graphql.org) query templates in a flexible and human-friendly way.

A GQT template declaratively defines the constraints of both the structure and input values of a GraphQL operation.

```graphql
# This is a GQT query template example.
query {
  # Allow selecting a user with any id.
  user(id: *) {
    id
    name
    max 2 { # Allow a maximum of 2 selections out of this set.
      birthdate
      email
      address {country city street}
    }
    # Allow a maximum limit of 100 friends.
    friends(after: *, limit=$friendsLimit: < 100) {
      id name birthdate email
      # Allow a maximum limit of 100 friends in total with a maximum depth of 2.
      # limit is restricted to a maximum value of 100 divided by the limit above.
      friends(after: *, limit: < 100 / $friendsLimit) {
        id name birthdate email
      }
    }
    pictures(
      # Allow a string with a maximum byte length of 64.
      prefix: len <= 64,
      # Allow an array with a maximum of 8 items,
      # where each item string is not longer than 64 bytes.
      tags: len < 8 && [...len < 64],
      # Allow only a subset of possible enum values in the category argument.
      category: PUBLIC || FRIENDSONLY,
      # Allow a rating value between 10 and 20
      rating: > 10 && < 20, 
    ) { url }
  }
}
```

GraphQL Schema:

```graphqls
# This is a GQT query template example.
scalar Date
enum PictureCategory {
    PUBLIC
    FRIENDSONLY
    PRIVATE
}
type Query { user(id: ID!): User }
type User {
    id:        ID!
    name:      String
    birthdate: Date
    email:     String
    address:   Address
    friends(
        after: ID,
        limit: Int!,
    ): [User!]!
    pictures(
        prefix:   String,
        tags:     [String!],
        category: PictureCategory,
        rating:   Int,
    ): [Picture!]
}
type Address {
    country: String!
    city:    String!
    street:  String!
}
type Picture { url: String! }
```

## Features
- Flexible restriction of the structure of a GraphQL request including max depth.
- Intuitive GraphQL-like syntax.
- Schema-aware mode with full type checking and validation.
- Schemaless mode (no validation against a GraphQL schema).
- Arithmetic and boolean expressions in input value constraints.
- Restriction of the maximum number of selections inside a `max` set.

Full documentation is available at [docs.graphguard.io/gqt](https://docs.graphguard.io/gqt.html).
