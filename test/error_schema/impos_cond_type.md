```graphql
type Query {name: String!}
type Foo {name: String!}
```

```graphql
query { ... on Foo { name } }
```

```
1:16: type "Query" can never be of type "Foo"
```