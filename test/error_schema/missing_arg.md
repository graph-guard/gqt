```graphql
type Query {foo(a: Int! b: Int): Int!}
```

```graphql
query { foo(b: *) }
```

```
1:12: argument "a" is required by schema but missing
```