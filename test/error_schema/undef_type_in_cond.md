```graphql
type Query {i: Interface!}
interface Interface {name: String!}
type Foo {name: String!}
type Bar {name: String!}
```

```graphql
query { i { ... on Baz { name } } }
```

```
1:20: type Baz is undefined in schema
```