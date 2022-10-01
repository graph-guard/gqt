```graphql
type Query {name: String!}
enum Enum { foo bar }
```

```graphql
query { ... on Enum { name } }
```

```
1:16: fragment can't condition on enum type Enum
```