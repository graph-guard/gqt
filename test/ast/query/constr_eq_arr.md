```graphql
query {
  x(a: [])
  y(
    a: != [1,2,3],
    b: [>5, <6],
  )
}
```

```yaml
Operation[1:1](query):
  - SelectionField[2:3](x):
    arguments:
      - Argument[2:5](a):
        - ConstrEquals[2:8]:
          - Array[2:8](0 items)
  - SelectionField[3:3](y):
    arguments:
      - Argument[4:5](a):
        - ConstrNotEquals[4:8]:
          - Array[4:11](3 items):
            - ConstrEquals[4:12]:
              - Int[4:12](1)
            - ConstrEquals[4:14]:
              - Int[4:14](2)
            - ConstrEquals[4:16]:
              - Int[4:16](3)
      - Argument[5:5](b):
        - ConstrEquals[5:8]:
          - Array[5:8](2 items):
            - ConstrGreater[5:9]:
              - Int[5:10](5)
            - ConstrLess[5:13]:
              - Int[5:14](6)

```