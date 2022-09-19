```graphql
type Query {string(a: Int!): String!}
```

```graphql
query { string(a: "not okay") }
```

```
1:19: expected type "Int!" but received "String"
```