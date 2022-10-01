```graphql
query {
  x ( a : [ ... [...10] ] )
}
```

```yaml
Operation[1:1](Query):
  - SelectionField[2:3](x):
    arguments:
      - Argument[2:7](a):
        - ConstrMap[2:11]:
          - ConstrMap[2:17]:
            - ConstrEquals[2:21]:
              - Int[2:21](10)

```