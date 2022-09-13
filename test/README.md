Tests can be grouped by directories recursively.
All directories must contain `.md` files exclusively.
Each `.md` test file must contain two code blocks,
one for the input and another one for the expected output, for example:

```graphql
query { f() }
```

```
1:10: empty argument list
```