```graphql
type Query {
    foo: String!
    bar: String!
}
```

```graphql
query { foo baz }
```

```
1:13: field "baz" is undefined in type Query
```