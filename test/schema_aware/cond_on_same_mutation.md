```graphql
type Mutation {name: String!}
```

```graphql
mutation { ... on Mutation { ... on Mutation { name } } }
```

```yaml
Operation[1:1](Mutation):
  - SelectionInlineFrag[1:12](Mutation):
    selections:
      - SelectionInlineFrag[1:30](Mutation):
        selections:
          - SelectionField[1:48](name)

```