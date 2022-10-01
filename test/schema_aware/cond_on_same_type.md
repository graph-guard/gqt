```graphql
type Query {foo: Foo!}
type Foo {name: String!}
```

```graphql
query { foo { ... on Foo { ... on Foo { name } } } }
```

```yaml
Operation[1:1](Query):
  - SelectionField[1:9](foo):
    selections:
      - SelectionInlineFrag[1:15](Foo):
        selections:
          - SelectionInlineFrag[1:28](Foo):
            selections:
              - SelectionField[1:41](name)

```