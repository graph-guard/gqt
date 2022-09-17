```graphql
query {
    x {
        __typename
        ... on Foo { foo __typename }
        ... on Bar { bar __typename }
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
          - SelectionField[4:26](__typename)
      - SelectionInlineFrag[5:9](Bar):
        selections:
          - SelectionField[5:22](bar)
          - SelectionField[5:26](__typename)

```