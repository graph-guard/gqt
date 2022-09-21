```graphql
query {
  x(a: {x:null})
  y(
    a: != {a=$a:1, b:$a},
    b: {a:*,b:*,c:*},
  )
}
```

```yaml
Operation[1:1](query):
  - SelectionField[2:3](x):
    arguments:
      - Argument[2:5](a):
        - ConstrEquals[2:8]:
          - Object[2:8](1 fields):
            - ObjectField[2:9](x):
              - ConstrEquals[2:11]:
                - Null[2:11]
  - SelectionField[3:3](y):
    arguments:
      - Argument[4:5](a):
        - ConstrNotEquals[4:8]:
          - Object[4:11](2 fields):
            - ObjectField[4:12](a=$a):
              - ConstrEquals[4:17]:
                - Int[4:17](1)
            - ObjectField[4:20](b):
              - ConstrEquals[4:22]:
                - Variable[4:22](a)
      - Argument[5:5](b):
        - ConstrEquals[5:8]:
          - Object[5:8](3 fields):
            - ObjectField[5:9](a):
              - Any[5:11]
            - ObjectField[5:13](b):
              - Any[5:15]
            - ObjectField[5:17](c):
              - Any[5:19]

```