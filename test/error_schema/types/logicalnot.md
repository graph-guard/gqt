```graphql
type Query {foo(a: Int!): Int!}
```

```graphql
query { foo(a: !false) }
```

```
1:16: expected type Int! but received Boolean
```