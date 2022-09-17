```graphql
type Query {u: BazzOrTazzOrFoo!}
union BazzOrTazzOrFoo = Bazz | Tazz | Foo
type Bazz {name: String!}
type Tazz {name: String!}
type Foo implements InterfaceFooBar {name: String!}
type Bar implements InterfaceFooBar {name: String!}
interface InterfaceFooBar {name: String!}
```

```graphql
query { u { ... on InterfaceFooBar { name } } }
```

```yaml
Operation[1:1](query):
  - SelectionField[1:9](u):
    selections:
      - SelectionInlineFrag[1:13](InterfaceFooBar):
        selections:
          - SelectionField[1:38](name)

```