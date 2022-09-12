```graphql
query {
  max 1 {
    foo
    bar
  }
  ... on bar {
    fazz
  }
  bar
}
```

```yaml
9:3: redeclared field
```