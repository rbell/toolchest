# SliceOps

This package provides a set of utility functions for working with slices in Go. These functions include operations such as cutting, removing, inserting, filtering, pushing, popping, and more.

## Functions

### Cut

The `Cut` function removes the elements `s[i:j]` from the slice and returns them, while preserving the order of `s`.

Example:
```go
s := []int{1, 2, 3, 4, 5}
cut := Cut(&s, 1, 3)
fmt.Println(cut) // Output: [2 3]
fmt.Println(s) // Output: [1 4 5]
```

### Remove

The `Remove` function removes the elements `s[i:j]` from the slice, while preserving the order of `s` and zeroing out the remaining elements to prevent memory leak.

Example:
```go
s := []int{1, 2, 3, 4, 5}
Remove(&s, 1, 3)
fmt.Println(s) // Output: [1 4 5]
```

### Insert

The `Insert` function inserts the elements `v` into the slice `s` at index `i`.

Example:
```go
s := []int{1, 2, 3, 4, 5}
s = Insert(s, 2, 7, 8)
fmt.Println(s) // Output: [1 2 7 8 3 4 5]
```

### FilterInPlace

The `FilterInPlace` function removes the elements of `s` that do not satisfy the `keep` predicate.

Example:
```go
s := []int{1, 2, 3, 4, 5}
FilterInPlace(&s, func(v int) bool { return v%2 == 0 })
fmt.Println(s) // Output: [2 4]
```

### Push

The `Push` function appends the elements `v` to the slice `s`.

Example:
```go
s := []int{1, 2, 3, 4, 5}
Push(&s, 6, 7)
fmt.Println(s) // Output: [6 7 1 2 3 4 5]
```

### Pop

The `Pop` function removes the first element of the slice `s` and returns it.

Example:
```go
s := []int{1, 2, 3, 4, 5}
v := Pop(&s)
fmt.Println(v) // Output: 1
fmt.Println(s) // Output: [2 3 4 5]
```

### Distinct

The `Distinct` function returns a new slice containing only the distinct elements of `s`.

Example:
```go
s := []int{1, 2, 2, 3, 3, 3}
s = Distinct(s)
fmt.Println(s) // Output: [1 2 3]
```

### Union

The `Union` function returns a new slice containing the distinct union of set of slices.

Example:
```go
s1 := []int{1, 2, 3}
s2 := []int{3, 4, 5}
s := Union(s1, s2)
fmt.Println(s) // Output: [1 2 3 4 5]
```

### Intersection

The `Intersection` function returns a new slice containing the distinct intersection of set of slices.

Example:
```go
s1 := []int{1, 2, 3}
s2 := []int{2, 3, 4}
s := Intersection(s1, s2)
fmt.Println(s) // Output: [2 3]
```

### Difference

The `Difference` function returns a new slice containing distinct elements in `s1` that are not in `s2`.

Example:
```go
s1 := []int{1, 2, 3}
s2 := []int{2, 3, 4}
s := Difference(s1, s2)
fmt.Println(s) // Output: [1]
```

### Disjoin

The `Disjoin` function returns a new slice containing distinct elements not in common between the slices provided.

Example:
```go
s1 := []int{1, 2, 3}
s2 := []int{2, 3, 4}
s3 := []int{3, 4, 5}
s := Disjoin(s1, s2, s3)
fmt.Println(s) // Output: [1 5]
```

## License

This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.