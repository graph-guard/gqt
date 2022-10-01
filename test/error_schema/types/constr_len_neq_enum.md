```graphql
type Query { f(a: [Int!]):Int! }
enum Color { red green blue }
```

```graphql
query { f(a: len != green) }
```

```
1:21: expected type Int! but received Color
```