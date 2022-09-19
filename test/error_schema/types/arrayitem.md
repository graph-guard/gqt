```graphql
type Query {foo(a: [Int!]!): Int!}
```

```graphql
query { foo(a: [ null ]) }
```

```
1:18: expected type "Int!" but received "null"
```