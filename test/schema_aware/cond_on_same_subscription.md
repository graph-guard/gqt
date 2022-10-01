```graphql
type Subscription {name: String!}
```

```graphql
subscription { ... on Subscription { ... on Subscription { name } } }
```

```yaml
Operation[1:1](Subscription):
  - SelectionInlineFrag[1:16](Subscription):
    selections:
      - SelectionInlineFrag[1:38](Subscription):
        selections:
          - SelectionField[1:60](name)

```