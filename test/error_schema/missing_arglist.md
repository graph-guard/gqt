```graphql
type Query {foo(a: Int! b: Int): Int!}
```

```graphql
query { foo }
```

```
1:9: argument "a" of type "Int!" is required but missing
```