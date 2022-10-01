```graphql
type Query { f(a: Int! b: Boolean!):Int! }
```


```graphql
query { f(a=$a: *, b: true == $a) }
```

```
1:23: mismatching types Boolean and Int!
```