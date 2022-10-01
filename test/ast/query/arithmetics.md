```graphql
query {
  a(
    a: > 1 + 4,
    b: > 1 + 4 * 2,
    c: > ( 1 + 4 ) * 2 / 10.4,
    d: > 11 % 2 - -10,
  )
}
```

```yaml
Operation[1:1](Query):
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
          - Division[5:10]:
            - Multiplication[5:10]:
              - Parentheses[5:10]:
                - Addition[5:12]:
                  - Int[5:12](1)
                  - Int[5:16](4)
              - Int[5:22](2)
            - Float[5:26](10.4)
      - Argument[6:5](d):
        - ConstrGreater[6:8]:
          - Subtraction[6:10]:
            - Modulo[6:10]:
              - Int[6:10](11)
              - Int[6:15](2)
            - Int[6:19](-10)

```