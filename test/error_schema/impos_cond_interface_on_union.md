```graphql
type Query {u: BazzOrTazz!}
union BazzOrTazz = Bazz | Tazz
type Bazz {name: String!}
type Tazz {name: String!}
type Foo implements InterfaceFooBar {name: String!}
type Bar implements InterfaceFooBar {name: String!}
interface InterfaceFooBar {name: String!}
```

```graphql
query { u { ... on InterfaceFooBar { name } } }
```

```
1:20: type BazzOrTazz can never be of type InterfaceFooBar
```