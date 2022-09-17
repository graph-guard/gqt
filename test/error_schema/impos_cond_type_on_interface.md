```graphql
type Query {i: Interface!}
interface Interface {name: String!}
type Foo implements Interface {name: String!}
type Bar {name: String!}
```

```graphql
query { i { ... on Bar { name } } }
```

```
1:20: type "Interface" can never be of type "Bar"
```