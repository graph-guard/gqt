```graphql
query {
  a(a:*)
  b(a1:*, a2:*)
  c(a1:*, a2:*, a3:*)
  d(
    a1:*,
    a2:*,
    a3:*
  )
  e(
    a1:*,
    a2:*,
    a3:*,
  )
  f(a=$fa:*)
  g(a=$ga:* , b = $gb:*)
}
```

```yaml
Operation[1:1](Query):
  - SelectionField[2:3](a):
    arguments:
      - Argument[2:5](a):
        - Any[2:7]
  - SelectionField[3:3](b):
    arguments:
      - Argument[3:5](a1):
        - Any[3:8]
      - Argument[3:11](a2):
        - Any[3:14]
  - SelectionField[4:3](c):
    arguments:
      - Argument[4:5](a1):
        - Any[4:8]
      - Argument[4:11](a2):
        - Any[4:14]
      - Argument[4:17](a3):
        - Any[4:20]
  - SelectionField[5:3](d):
    arguments:
      - Argument[6:5](a1):
        - Any[6:8]
      - Argument[7:5](a2):
        - Any[7:8]
      - Argument[8:5](a3):
        - Any[8:8]
  - SelectionField[10:3](e):
    arguments:
      - Argument[11:5](a1):
        - Any[11:8]
      - Argument[12:5](a2):
        - Any[12:8]
      - Argument[13:5](a3):
        - Any[13:8]
  - SelectionField[15:3](f):
    arguments:
      - Argument[15:5](a=$fa):
        - Any[15:11]
  - SelectionField[16:3](g):
    arguments:
      - Argument[16:5](a=$ga):
        - Any[16:11]
      - Argument[16:15](b=$gb):
        - Any[16:23]

```