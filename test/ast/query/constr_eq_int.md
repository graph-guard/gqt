```graphql
query {
  x(a: 0)
  y(
    a: 2
    b: != 123
  )
}
```

```yaml
Operation[1:1](query):
  - SelectionField[2:3](x):
    arguments:
      - Argument[2:5](a):
        - ConstrEquals[2:8]:
          - Int[2:8](0)
  - SelectionField[3:3](y):
    arguments:
      - Argument[4:5](a):
        - ConstrEquals[4:8]:
          - Int[4:8](2)
      - Argument[5:5](b):
        - ConstrNotEquals[5:8]:
          - Int[5:11](123)

```