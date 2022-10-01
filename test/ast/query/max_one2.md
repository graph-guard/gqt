```graphql
query {
  max 1 {
    foo
    bar
  }
  max
}
```

```yaml
Operation[1:1](Query):
  - SelectionMax[2:3](1):
    options:
      - SelectionField[3:5](foo)
      - SelectionField[4:5](bar)
  - SelectionField[6:3](max)

```