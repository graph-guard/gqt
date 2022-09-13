```graphql
query {
  ... on bar {
    fazz
  }
  bar
  max 1 {
    foo
    bar
  }
}
```

```yaml
8:5: redeclared field
```