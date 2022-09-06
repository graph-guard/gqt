```graphql
query {
  x(a = $a: 1)
  y(
    a: != $a,
    b: > $a,
  )
}
```

```yaml
Operation[1:1](query):
  - SelectionField[2:3](x):
    arguments:
      - Argument[2:5](a=$a):
        - ConstrEquals[2:13]:
          - Int[2:13](1)
  - SelectionField[3:3](y):
    arguments:
      - Argument[4:5](a):
        - ConstrNotEquals[4:8]:
          - Variable[4:11](a)
      - Argument[5:5](b):
        - ConstrGreater[5:8]:
          - Variable[5:10](a)

```