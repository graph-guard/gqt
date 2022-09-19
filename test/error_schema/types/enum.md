```graphql
type Query {foo(a: Color!): Color!}
enum Color {
    red
    green
    blue
    bananayellow
}
enum Fruit {
    banana
    orange
}
```

```graphql
query { foo(a: banana) }
```

```
1:16: expected type "Color!" but received "Fruit"
```