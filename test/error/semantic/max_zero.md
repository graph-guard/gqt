```graphql
query {
  max 0 {
    foo
    bar
  }
}
```

```yaml
2:7: maximum number of options must be an unsigned integer greater 0
```