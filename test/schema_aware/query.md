```graphql
type Query {i: Interface!}
interface Interface {
  string: String!
  boolean: Boolean!
  union: UnionFooBar
}
type Foo implements Interface {
  string: String!
  boolean: Boolean!
  union: UnionFooBar
}
type Bar {
  string: String!
  boolean: Boolean!
  union: UnionFooBar
}
union UnionFooBar = Foo | Bar
```

```graphql
query {
  i {
    __typename
    string
    boolean
    union {
      ... on Interface {
        __typename
        string
        boolean
        union { __typename }
      }
      ... on Foo {
        __typename
        string
        boolean
        union { __typename }
      }
      ... on Bar {
        __typename
        string
        boolean
        union { __typename }
      }
      ... on UnionFooBar {
        ... on Bar { __typename }
        ... on Foo { __typename }
      }
    }
  }
}
```

```yaml
Operation[1:1](query):
  - SelectionField[2:3](i):
    selections:
      - SelectionField[3:5](__typename)
      - SelectionField[4:5](string)
      - SelectionField[5:5](boolean)
      - SelectionField[6:5](union):
        selections:
          - SelectionInlineFrag[7:7](Interface):
            selections:
              - SelectionField[8:9](__typename)
              - SelectionField[9:9](string)
              - SelectionField[10:9](boolean)
              - SelectionField[11:9](union):
                selections:
                  - SelectionField[11:17](__typename)
          - SelectionInlineFrag[13:7](Foo):
            selections:
              - SelectionField[14:9](__typename)
              - SelectionField[15:9](string)
              - SelectionField[16:9](boolean)
              - SelectionField[17:9](union):
                selections:
                  - SelectionField[17:17](__typename)
          - SelectionInlineFrag[19:7](Bar):
            selections:
              - SelectionField[20:9](__typename)
              - SelectionField[21:9](string)
              - SelectionField[22:9](boolean)
              - SelectionField[23:9](union):
                selections:
                  - SelectionField[23:17](__typename)
          - SelectionInlineFrag[25:7](UnionFooBar):
            selections:
              - SelectionInlineFrag[26:9](Bar):
                selections:
                  - SelectionField[26:22](__typename)
              - SelectionInlineFrag[27:9](Foo):
                selections:
                  - SelectionField[27:22](__typename)

```