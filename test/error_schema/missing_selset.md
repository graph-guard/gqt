```graphql
type Query {foo: Foo!}
type Foo {bar: String!}
```

```graphql
query { foo }
```

```
1:9: missing selection set for field "foo" of type Foo!
```