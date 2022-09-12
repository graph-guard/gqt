```graphql
query {
  ... on bar {
    fazz
  }
  bar
  max 1 {
    foo
    ... on bar {
      traz
    }
  }
}
```

```yaml
8:5: redeclared type condition
```