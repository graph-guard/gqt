```graphql
type Query {u: FooOrBar!}
union FooOrBar = Foo | Bar
type Foo implements Interface {name: String!}
type Bar implements Interface {name: String!}
interface Interface {name: String!}
```

```graphql
query { u { ... on Interface { name } } }
```

```yaml
Operation[1:1](query):
  - SelectionField[1:9](u):
    selections:
      - SelectionInlineFrag[1:13](Interface):
        selections:
          - SelectionField[1:32](name)

```