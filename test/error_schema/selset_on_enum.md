```graphql
type Query {e: Enum!}
enum Enum {foo bar baz}
```

```graphql
query { e { foo } }
```

```
1:13: field "foo" is undefined in type Enum
```