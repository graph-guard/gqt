```graphql
type Query {e: String!}
```

```graphql
query { e { foo } }
```

```
1:13: field "foo" is undefined in type "String"
```