```graphql
type Query {name: String!}
```

```graphql
query { ... on Int { name } }
```

```
1:16: fragment can't condition on scalar type "Int"
```