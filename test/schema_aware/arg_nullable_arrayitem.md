```graphql
type Query {f(a: [Int]!):Int!}
```

```graphql
query { f(a: [42, null]) }
```

```yaml
Operation[1:1](query):
  - SelectionField[1:9](f):
    arguments:
      - Argument[1:11](a):
        - ConstrEquals[1:14]:
          - Array[1:14](2 items):
            - ConstrEquals[1:15]:
              - Int[1:15](42)
            - ConstrEquals[1:19]:
              - Null[1:19]

```