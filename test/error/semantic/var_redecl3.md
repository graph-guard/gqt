```graphql
query { a(a=$a) b { bb(a: {f1, f2=$a}) } }
```

```
1:35: redeclared variable
```