```graphql
query { ...on A{__typename} ...on B{__typename} ...on A{__typename} }
```

```
1:49: redeclared type condition
```