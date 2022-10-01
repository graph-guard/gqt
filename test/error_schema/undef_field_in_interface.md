```graphql
type Query {
    i: Interface!
}
interface Interface {
    string: String!
}
```

```graphql
query { i { foo } }
```

```
1:13: field "foo" is undefined in type Interface
```