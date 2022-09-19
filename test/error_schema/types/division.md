```graphql
type Query {foo(a: Boolean!): Int!}
```

```graphql
query { foo(a: 3 / 2) }
```

```
1:16: expected type "Boolean!" but received "Int|Float"
```