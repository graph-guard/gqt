```graphql
query {
  max -1 {
    foo
    bar
  }
}
```

```yaml
2:7: maximum number of options must be an unsigned integer greater 0
```