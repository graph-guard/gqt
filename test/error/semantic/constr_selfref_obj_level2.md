```graphql
query { f(a:{foo=$foo:{bar:$foo}}) }
```

```
1:28: object field self reference in constraint
```