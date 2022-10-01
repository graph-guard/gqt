```graphql
type Query {u: U!}
type Foo {name: String!}
type Bar {name: String!}
union U = Foo | Bar
```

```graphql
query { u { ... on U { ... on U { __typename } } } }
```

```yaml
Operation[1:1](Query):
  - SelectionField[1:9](u):
    selections:
      - SelectionInlineFrag[1:13](U):
        selections:
          - SelectionInlineFrag[1:24](U):
            selections:
              - SelectionField[1:35](__typename)

```