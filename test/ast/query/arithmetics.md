```graphql
query {
  a(
    a: > 1 + 4,
    b: > 1 + 4 * 2,
    c: > ( 1 + 4 ) * 2,
  )
}
```

```yaml
Operation[1:1](query):
  - SelectionField[2:3](a):
    arguments:
      - Argument[3:5](a):
        - ConstrGreater[3:8]:
          - Addition[3:10]:
            - Int[3:10](1)
            - Int[3:14](4)
      - Argument[4:5](b):
        - ConstrGreater[4:8]:
          - Addition[4:10]:
            - Int[4:10](1)
            - Multiplication[4:14]:
              - Int[4:14](4)
              - Int[4:18](2)
      - Argument[5:5](c):
        - ConstrGreater[5:8]:
          - Multiplication[5:10]:
            - Parentheses[5:10]:
              - ConstrEquals[5:12]:
                - Addition[5:12]:
                  - Int[5:12](1)
                  - Int[5:16](4)
            - Int[5:22](2)

```