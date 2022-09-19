```graphql
type Query {foo(a: Color!): Color!}
enum Color {
    red
    green
    blue
}
```

```graphql
query { foo(a: yellow) }
```

```
1:16: undefined enum value "yellow"
```