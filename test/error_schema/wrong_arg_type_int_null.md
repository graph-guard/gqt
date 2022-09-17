```graphql
type Query {string(a: Int!): String!}
```

```graphql
query { string(a: null) }
```

```
1:19: wrong type for argument "a"
```