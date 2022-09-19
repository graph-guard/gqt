```graphql
type Query {foo(a: Int!): Int!}
```

```graphql
query { foo(b: *) }
```

```
1:13: argument is undefined in schema
```