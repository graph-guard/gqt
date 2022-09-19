```graphql
type Query {f(a: Int):Int!}
```

```graphql
query { f(a: 42) }
```

```yaml
Operation[1:1](query):
  - SelectionField[1:9](f):
    arguments:
      - Argument[1:11](a):
        - ConstrEquals[1:14]:
          - Int[1:14](42)

```