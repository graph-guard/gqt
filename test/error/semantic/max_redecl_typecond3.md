```graphql
query {
  max 1 {
    foo
    bar
  }
  max 1 {
    fizz
    buzz
  }
}
```

```yaml
6:3: multiple max blocks in one selection set
```