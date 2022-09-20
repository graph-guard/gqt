```graphql
type Query { f(a: Int! b: Boolean!):Int! }
```


```graphql
query { f(a=$a: *, b: true == $a) }
```

```
1:24: can't use "Int!" as "Boolean"
```