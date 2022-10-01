```graphql
query {
  x(a: null)
  y(
    a: null,
    b: != null,
  )
}
```

```yaml
Operation[1:1](Query):
  - SelectionField[2:3](x):
    arguments:
      - Argument[2:5](a):
        - ConstrEquals[2:8]:
          - Null[2:8]
  - SelectionField[3:3](y):
    arguments:
      - Argument[4:5](a):
        - ConstrEquals[4:8]:
          - Null[4:8]
      - Argument[5:5](b):
        - ConstrNotEquals[5:8]:
          - Null[5:11]

```