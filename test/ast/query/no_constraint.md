```graphql
query {
  foo {
    bar
    baz
  }
  fazz
}
```

```yaml
Operation[1:1](Query):
  - SelectionField[2:3](foo):
    selections:
      - SelectionField[3:5](bar)
      - SelectionField[4:5](baz)
  - SelectionField[6:3](fazz)

```