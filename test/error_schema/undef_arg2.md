```graphql
type Query {foo: Foo!}
type Foo {f(a: Int): Int!}
```

```graphql
query { foo { f(b: *) }}
```

```
1:17: argument "b" is undefined on field "f" in type Foo
```