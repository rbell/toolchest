# Storage Package

The `storage` package provides concurrency-safe map and cache implementations in Go.

## SafeMap

`SafeMap` is a struct that wraps a map with a simple `RWMutex` to facilitate concurrency. It supports generic types for keys and values.

## FifoMapCache

`FifoMapCache` is a struct that implements a First-In-First-Out (FIFO) cache with a maximum size. When the cache is full, the oldest entries are evicted. It supports generic types for keys and values.

## GenericStack

`GenericStack` is a struct that implements a generic stack data structure. It supports any type of values.

## Testing

Unit tests for the `SafeMap` and `FifoMapCache` are located in the `storage/safeMap_test.go` and `storage/fifoMapCache_test.go` files respectively. They cover all methods and some edge cases, including concurrent operations.

# OrderedBTree

The `OrderedBTree` is a thread-safe implementation of a BTree that supports ordered types. It is a wrapper around the google/btree package.

## Usage

First, import the `storage` package:

```go
import "your_project/storage"
```

### Initialization

Create a new instance of `OrderedBTree`:

```go
btree := storage.NewOrderedBTree[int, string]()
```

### Set

To set a value for a key in the BTree:

```go
btree.Set(1, "value1")
```

### Get

To get a value for a key from the BTree:

```go
value, ok := btree.Get(1)
if ok {
    fmt.Println(value)
}
```

### Delete

To delete a key from the BTree:

```go
value, ok := btree.Delete(1)
if ok {
    fmt.Println("Deleted value:", value)
}
```

### Has

To check if a key exists in the BTree:

```go
exists := btree.Has(1)
fmt.Println("Key exists:", exists)
```

### Len

To get the length of the BTree:

```go
length := btree.Len()
fmt.Println("Length of BTree:", length)
```

### Min and Max

To get the minimum and maximum key and value in the BTree:

```go
minKey, minValue := btree.Min()
maxKey, maxValue := btree.Max()
fmt.Println("Min:", minKey, minValue)
fmt.Println("Max:", maxKey, maxValue)
```

### Ascend and Descend

To iterate over the BTree in ascending or descending order:

```go
btree.Ascend(func(key int, value *string) bool {
    fmt.Println(key, *value)
    return true
})

btree.Descend(func(key int, value *string) bool {
    fmt.Println(key, *value)
    return true
})
```

### DeleteMin and DeleteMax

To delete the minimum and maximum key and value in the BTree:

```go
minKey, minValue := btree.DeleteMin()
maxKey, maxValue := btree.DeleteMax()
fmt.Println("Deleted Min:", minKey, minValue)
fmt.Println("Deleted Max:", maxKey, maxValue)
```

### Clone

To clone the BTree:

```go
clone := btree.Clone()
```

# Tree

Tree provides a generic tree data structure with methods for adding and walking through the tree.

## Usage

Here are some examples of how to use the `Tree` package:

### Creating a new tree

```go
factory := func(value string) storage.childAdderGetter[string] {
    return NewNode(value)
}
tree := storage.NewTree(factory)
```

### Adding an ancestry chain to the tree
```go
ancestry := []string{"top", "child1", "child1.1"}
    err := storage.AddAncestryChain(tree, ancestry...)
    if err != nil {
    log.Fatal(err)
}
```

### Walking through the tree
```go
tree.Walk(func(s string, level int) {
    fmt.Println(strings.Repeat("->", level), s)
})
```

# License
This project is licensed under the Apache Public License, version 2.0. See the LICENSE file for details.