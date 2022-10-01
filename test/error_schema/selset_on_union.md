```graphql
type Query {
    u: UnionFooBar!
}
type Foo { string: String! }
type Bar { string: String! }
union UnionFooBar = Foo | Bar
```

```graphql
query { u { foo } }
```

```
1:13: field "foo" is undefined in type UnionFooBar
```