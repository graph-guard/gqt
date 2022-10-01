```graphql
type Query { f(a: [Int!]):Int! }
```


```graphql
query { f(a: len < {f:2}) }
```

```
1:20: expected type Int! but received {f:Int}
```