```graphql
type Query {i: Interface!}
union FooOrBar = Foo | Bar
type Foo implements Interface {name: String!}
type Bar implements Interface {name: String!}
interface Interface {name: String!}
```

```graphql
query { i { ... on Interface { ... on Interface { name } } } }
```

```yaml
Operation[1:1](Query):
  - SelectionField[1:9](i):
    selections:
      - SelectionInlineFrag[1:13](Interface):
        selections:
          - SelectionInlineFrag[1:32](Interface):
            selections:
              - SelectionField[1:51](name)

```