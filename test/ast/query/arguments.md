```graphql
query {
  a(a)
  b(a1, a2)
  c(a1, a2, a3)
  d(
    a1,
    a2,
    a3
  )
  e(
    a1,
    a2,
    a3,
  )
  f(a=$fa)
  g(a=$ga , b = $gb)
}
```

```yaml
Operation[1:1](query):
  - SelectionField[2:3](a):
    arguments:
      - Argument[2:5](a)
  - SelectionField[3:3](b):
    arguments:
      - Argument[3:5](a1)
      - Argument[3:9](a2)
  - SelectionField[4:3](c):
    arguments:
      - Argument[4:5](a1)
      - Argument[4:9](a2)
      - Argument[4:13](a3)
  - SelectionField[5:3](d):
    arguments:
      - Argument[6:5](a1)
      - Argument[7:5](a2)
      - Argument[8:5](a3)
  - SelectionField[10:3](e):
    arguments:
      - Argument[11:5](a1)
      - Argument[12:5](a2)
      - Argument[13:5](a3)
  - SelectionField[15:3](f):
    arguments:
      - Argument[15:5](a=$fa)
  - SelectionField[16:3](g):
    arguments:
      - Argument[16:5](a=$ga)
      - Argument[16:13](b=$gb)

```