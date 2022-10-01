```graphql
type Query {f(a: [Int]):Int!}
```

```graphql
query { f(a: []) }
```

```yaml
Operation[1:1](Query):
  - SelectionField[1:9](f):
    arguments:
      - Argument[1:11](a):
        - ConstrEquals[1:14]:
          - Array[1:14](0 items)

```