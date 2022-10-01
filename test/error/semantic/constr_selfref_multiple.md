```graphql
query {
    f(
        a=$a:{
            foo=$foo:{
                foo:$foo,
                a:$a,
                bar=$bar:$bar
            }
        }
    )
}
```

```
5:21: object field self reference in constraint
6:19: argument self reference in constraint
7:26: object field self reference in constraint
```