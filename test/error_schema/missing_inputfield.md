```graphql
type Query {f(a: Input!):Int!}
input Input {
    req: Int!
    opt: Int
}
```

```graphql
query { f(a: {opt: 42}) }
```

```
1:14: field "req" of type "Int!" is required but missing
```