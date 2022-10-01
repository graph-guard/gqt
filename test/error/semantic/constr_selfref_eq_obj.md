```graphql
query { f(a=$a: {foo: 42, bar: $a}) }
```

```
1:32: argument self reference in constraint
```