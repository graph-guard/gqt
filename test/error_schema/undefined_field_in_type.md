```graphql
type Query {foo: Foo!}
type Foo {bar: String!}
```

```graphql
query { foo { baz } }
```

```
1:15: field "baz" is undefined in type "Foo"
```