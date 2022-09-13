```graphql
query {
  max 3 {
    foo
    bar
    baz
    fazz(a: 10)
  }
}
```

```yaml
Operation[1:1](query):
  - SelectionMax[2:3](3):
    options:
      - SelectionField[3:5](foo)
      - SelectionField[4:5](bar)
      - SelectionField[5:5](baz)
      - SelectionField[6:5](fazz):
        arguments:
          - Argument[6:10](a):
            - ConstrEquals[6:13]:
              - Int[6:13](10)

```