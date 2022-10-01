```graphql
type Query {name: String!}
input Input {name: String!}
```

```graphql
query { ... on Input { name } }
```

```
1:16: fragment can't condition on input type Input
```