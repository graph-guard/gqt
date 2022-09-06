```graphql
query {
  x(a: < 10)
  y(
    a: > 3,
    b: <= 10,
    c: >= 3.14,
  )
}
```

```yaml
Operation[1:1](query):
  - SelectionField[2:3](x):
    arguments:
      - Argument[2:5](a):
        - ConstrLess[2:8]:
          - Int[2:10](10)
  - SelectionField[3:3](y):
    arguments:
      - Argument[4:5](a):
        - ConstrGreater[4:8]:
          - Int[4:10](3)
      - Argument[5:5](b):
        - ConstrLessOrEqual[5:8]:
          - Int[5:11](10)
      - Argument[6:5](c):
        - ConstrGreaterOrEqual[6:8]:
          - Float[6:11](3.14)

```