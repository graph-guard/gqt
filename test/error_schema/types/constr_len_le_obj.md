```graphql
type Query { f(a: [Int!]):Int! }
input Input { x: Int! }
```


```graphql
query { f(a: len < {f:2}) }
```

```
1:20: mismatching types: can't use object as number
```