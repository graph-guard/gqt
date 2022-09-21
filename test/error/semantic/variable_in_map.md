```graphql
query { f(a: $f, b: [...{f = $f:*}]) }
```

```
1:30: the use of variables inside map constraints is prohibited
```