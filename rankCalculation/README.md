# Rank Calculator in Go

The `RankCalculator` is a thread-safe implementation of a rank calculator. It is used to calculate the rank of entries based on either their number of hits or the rank can be calculated based on the position of the entry in the sorted list or based on the value of the entry.  A cusomized Ranker can be passed into the Rank Calculator as an option to calculate the rank based on the value of the entry.

## Usage

First, import the `rankCalculation` package:

```go
import "github.com/rbell/rankCalculation"
```

### Creating a new Rank Calculator

You can create a new `RankCalculator` using the `NewRankCalculator` function:

```go
calculator := rankCalculation.NewRankCalculator[int]()
```

### Accumulating Entries

To add entries to the `RankCalculator`, use the `Accumulate` method:

```go
calculator.Accumulate(1)
calculator.Accumulate(2)
calculator.Accumulate(3)
```

### Resetting the Rank Calculator

To reset the `RankCalculator` and clear all entries, use the `Reset` method:

```go
calculator.Reset()
```

### Calculating Ranks

To calculate the ranks of the accumulated entries, use the `Calculate` method:

```go
ranks, err := calculator.Calculate()
if err != nil {
    log.Fatal(err)
}
fmt.Println(ranks)
```

This will output a map where the keys are the entries and the values are their respective ranks.

## Examples

Here is a complete example of how to use the `RankCalculator`:

```go
package main

import (
    "fmt"
    "log"
    "github.com/rbell/rankCalculation"
)

func main() {
    // Create a new RankCalculator
    calculator := rankCalculation.NewRankCalculator[int]()

    // Accumulate some entries
    calculator.Accumulate(1)
    calculator.Accumulate(2)
    calculator.Accumulate(3)

    // Calculate the ranks
    ranks, err := calculator.Calculate()
    if err != nil {
        log.Fatal(err)
    }

    // Print the ranks
    fmt.Println(ranks)

    // Reset the RankCalculator
    calculator.Reset()
}
```

In this example, the `RankCalculator` is used to calculate the ranks of the numbers 1, 2, and 3. The ranks are then printed to the console. Finally, the `RankCalculator` is reset, clearing all entries.