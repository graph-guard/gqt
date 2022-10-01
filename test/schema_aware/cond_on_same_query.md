```graphql
type Query {name: String!}
```

```graphql
query { ... on Query { ... on Query { name } } }
```

```yaml
Operation[1:1](Query):
  - SelectionInlineFrag[1:9](Query):
    selections:
      - SelectionInlineFrag[1:24](Query):
        selections:
          - SelectionField[1:39](name)

```