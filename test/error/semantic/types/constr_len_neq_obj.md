```graphql
query { f(a: len != {f:2}) }
```

```
1:21: expected type Int! but received {f:Int}
```