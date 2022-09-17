```graphql
type Query {foo: Foo!}
type Foo {bar: String!}
```

```graphql
query { foo }
```

```
1:13: missing selection set for selection "foo" of type Foo
```