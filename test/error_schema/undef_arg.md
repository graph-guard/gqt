```graphql
type Query {f(a: Int): Int!}
```

```graphql
query { f(b: *) }
```

```
1:11: argument "b" is undefined on field "f" in type Query
```