```graphql
query {
  max 1 {
    foo
    bar
    max 1 {
      bazz
      fuzz
    }
  }
}
```

```yaml
5:5: nested max block
```