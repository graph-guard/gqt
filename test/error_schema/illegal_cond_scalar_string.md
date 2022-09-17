```graphql
type Query {name: String!}
```

```graphql
query { ... on String { name } }
```

```
1:16: fragment can't condition on scalar type "String"
```