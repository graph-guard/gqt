```graphql
type Query { f(a: Int! b: Boolean!):Int! }
```


```graphql
query { f(a=$a: *, b: true == $a) }
```

```
1:31: expected type "Boolean!" but received "Int!"
```