```graphql
type Query { f(a: Int b: Int!):Int! }
```


```graphql
query { f(a=$a: *, b: $a) }
```

```
1:23: expected type "Int!" but received "Int"
```