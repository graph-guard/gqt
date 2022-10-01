```graphql
query {
  max max 1 {
    foo
    bar
  }
}
```

```yaml
Operation[1:1](Query):
  - SelectionField[2:3](max)
  - SelectionMax[2:7](1):
    options:
      - SelectionField[3:5](foo)
      - SelectionField[4:5](bar)

```