```graphql
query {
  x(a: true)
  y(
    a: false
    b: != true
  )
}
```

```yaml
Operation[1:1](query):
  - SelectionField[2:3](x):
    arguments:
      - Argument[2:5](a):
        - ConstrEquals[2:8]:
          - True[2:8]
  - SelectionField[3:3](y):
    arguments:
      - Argument[4:5](a):
        - ConstrEquals[4:8]:
          - False[4:8]
      - Argument[5:5](b):
        - ConstrNotEquals[5:8]:
          - True[5:11]

```