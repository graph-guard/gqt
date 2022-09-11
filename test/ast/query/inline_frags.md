```graphql
query {
    x {
        __typename
        ... on Foo { foo }
        ... on Bar { bar }
    }
}
```

```yaml
Operation[1:1](query):
  - SelectionField[2:5](x):
    selections:
      - SelectionField[3:9](__typename)
      - SelectionInlineFrag[4:9](Foo):
        selections:
          - SelectionField[4:22](foo)
      - SelectionInlineFrag[5:9](Bar):
        selections:
          - SelectionField[5:22](bar)

```