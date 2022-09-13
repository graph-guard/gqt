```graphql
query {
  bar
  max 1 {
    foo
    ... on bar {
      traz
    }
  }
  ... on bar {
    fazz
  }
}
```

```yaml
9:3: redeclared type condition
```