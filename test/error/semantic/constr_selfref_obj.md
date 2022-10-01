```graphql
query { f(a:{foo:{bar=$bar:$bar}}) }
```

```
1:28: object field self reference in constraint
```