```graphql
type Query {f(a: Input!):String!}
input Input {name: String!}
```

```graphql
query { f(a:{name: "okay", inexistent: "not okay"}) }
```

```
1:28: field "inexistent" is undefined in type Input
```