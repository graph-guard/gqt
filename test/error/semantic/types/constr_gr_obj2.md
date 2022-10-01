```graphql
query { f(a: > {i:2,a:[]}) }
```

```
1:16: expected number but received {i:Int,a:array}
```