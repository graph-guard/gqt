```graphql
type Query {f(a: [Int]):Int!}
```

```graphql
query { f(a: null) }
```

```yaml
Operation[1:1](Query):
  - SelectionField[1:9](f):
    arguments:
      - Argument[1:11](a):
        - ConstrEquals[1:14]:
          - Null[1:14]

```