```graphql
query {
  x(a: "")
  y(
    a: "bar",
    b: != "bazz",
  )
}
```

```yaml
Operation[1:1](Query):
  - SelectionField[2:3](x):
    arguments:
      - Argument[2:5](a):
        - ConstrEquals[2:8]:
          - String[2:8]("")
  - SelectionField[3:3](y):
    arguments:
      - Argument[4:5](a):
        - ConstrEquals[4:8]:
          - String[4:8]("bar")
      - Argument[5:5](b):
        - ConstrNotEquals[5:8]:
          - String[5:11]("bazz")

```