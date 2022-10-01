```graphql
query { f(a: "string" != {foo: "bar"}) }
```

```
1:14: mismatching types String and {foo:String}
```