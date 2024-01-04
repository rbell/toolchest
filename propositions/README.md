# Propositions Package

The `propositions` package provides a set of proposition functions in Go. These functions allow evaluations of various conditions on various types, each function returning either true or false.

## Map Propositions

### MapContainsKey

Checks if a map contains a specific key.

```go
m := map[string]int{"apple": 1, "banana": 2}
fmt.Println(propositions.MapContainsKey(m, "apple"))  // Output: true
fmt.Println(propositions.MapContainsKey(m, "orange")) // Output: false
```

### MapContainsValue

Checks if a map contains a specific value.

```go
m := map[string]int{"apple": 1, "banana": 2}
fmt.Println(propositions.MapContainsValue(m, 1)) // Output: true
fmt.Println(propositions.MapContainsValue(m, 3)) // Output: false
```

### MapKeyPropositionAny

Checks if any of the keys in a map satisfy a given proposition.

```go
m := map[int]string{1: "apple", 2: "banana"}
p := func(k int) bool { return k > 1 }
fmt.Println(propositions.MapKeyPropositionAny(m, p)) // Output: true
```

### MapKeyPropositionAll

Checks if all of the keys in a map satisfy a given proposition.

```go
m := map[int]string{1: "apple", 2: "banana"}
p := func(k int) bool { return k > 0 }
fmt.Println(propositions.MapKeyPropositionAll(m, p)) // Output: true
```

### MapKeyPropositionNone

Checks if none of the keys in a map satisfy a given proposition.

```go
m := map[int]string{1: "apple", 2: "banana"}
p := func(k int) bool { return k > 2 }
fmt.Println(propositions.MapKeyPropositionNone(m, p)) // Output: true
```

### MapValuePropositionAny

Checks if any of the values in a map satisfy a given proposition.

```go
m := map[string]int{"apple": 1, "banana": 2}
p := func(v int) bool { return v == 2 }
fmt.Println(propositions.MapValuePropositionAny(m, p)) // Output: true
```

### MapValuePropositionAll

Checks if all of the values in a map satisfy a given proposition.

```go
m := map[string]int{"apple": 1, "banana": 2}
p := func(v int) bool { return v > 0 }
fmt.Println(propositions.MapValuePropositionAll(m, p)) // Output: true
```

### MapValuePropositionNone

Checks if none of the values in a map satisfy a given proposition.

```go
m := map[string]int{"apple": 1, "banana": 2}
p := func(v int) bool { return v > 2 }
fmt.Println(propositions.MapValuePropositionNone(m, p)) // Output: true
```

## Slice Propositions

### SliceContains

Checks if a slice contains a specific value.

```go
s := []int{1, 2, 3}
fmt.Println(propositions.SliceContains(s, 2))  // Output: true
fmt.Println(propositions.SliceContains(s, 4)) // Output: false
```

### SliceContainsAll

Checks if a slice contains all the specified values.

```go
s := []int{1, 2, 3}
v := []int{1, 2}
fmt.Println(propositions.SliceContainsAll(s, v)) // Output: true
```

### SliceContainsAny

Checks if a slice contains any of the specified values.

```go
s := []int{1, 2, 3}
v := []int{4, 2}
fmt.Println(propositions.SliceContainsAny(s, v)) // Output: true
```

### SliceContainsNone

Checks if a slice contains none of the specified values.

```go
s := []int{1, 2, 3}
v := []int{4, 5}
fmt.Println(propositions.SliceContainsNone(s, v)) // Output: true
```

### SliceAllLessThan

Checks if all the values in the slice are less than the specified value.

```go
s := []int{1, 2, 3}
fmt.Println(propositions.SliceAllLessThan(s, 4)) // Output: true
```

### SliceAllLessThanOrEqualTo

Checks if all the values in the slice are less than or equal to the specified value.

```go
s := []int{1, 2, 3}
fmt.Println(propositions.SliceAllLessThanOrEqualTo(s, 3)) // Output: true
```

### SliceAllGreaterThan

Checks if all the values in the slice are greater than the specified value.

```go
s := []int{4, 5, 6}
fmt.Println(propositions.SliceAllGreaterThan(s, 3)) // Output: true
```

### SliceAllGreaterThanOrEqualTo

Checks if all the values in the slice are greater than or equal to the specified value.

```go
s := []int{4, 5, 6}
fmt.Println(propositions.SliceAllGreaterThanOrEqualTo(s, 4)) // Output: true
```

### SliceAnyLessThan

Checks if any of the values in the slice are less than the specified value.

```go
s := []int{4, 5, 6}
fmt.Println(propositions.SliceAnyLessThan(s, 5)) // Output: true
```

### SliceAnyLessThanOrEqualTo

Checks if any of the values in the slice are less than or equal to the specified value.

```go
s := []int{4, 5, 6}
fmt.Println(propositions.SliceAnyLessThanOrEqualTo(s, 6)) // Output: true
```

### SliceAnyGreaterThan

Checks if any of the values in the slice are greater than the specified value.

```go
s := []int{4, 5, 6}
fmt.Println(propositions.SliceAnyGreaterThan(s, 5)) // Output: true
```

### SliceAnyGreaterThanOrEqualTo

Checks if any of the values in the slice are greater than or equal to the specified value.

```go
s := []int{4, 5, 6}
fmt.Println(propositions.SliceAnyGreaterThanOrEqualTo(s, 4)) // Output: true
```

### SlicePropositionNone

Checks if none of the values in the slice satisfy the given proposition.

```go
s := []int{4, 5, 6}
p := func(v int) bool { return v > 6 }
fmt.Println(propositions.SlicePropositionNone(s, p)) // Output: true
```

### SlicePropositionAny

Checks if any of the values in the slice satisfy the given proposition.

```go
s := []int{4, 5, 6}
p := func(v int) bool { return v == 5 }
fmt.Println(propositions.SlicePropositionAny(s, p)) // Output: true
```

### SlicePropositionAll

Checks if all of the values in the slice satisfy the given proposition.

```go
s := []int{4, 5, 6}
p := func(v int) bool { return v > 3 }
fmt.Println(propositions.SlicePropositionAll(s, p)) // Output: true
```

## License

This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
