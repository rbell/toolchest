package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderedBTree_New_ReturnsInitializedBTree(t *testing.T) {
	// setup

	// test
	tree := NewOrderedBTree[int, string]()

	// assert
	assert.NotNil(t, tree)
	assert.NotNil(t, tree.BTree)
	assert.NotNil(t, tree.writeMux)
}

func TestOrderedBTree_Set_SetsValueForGivenKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1
	value := "value"

	// test
	tree.Set(key, &value)

	// assert
	assert.Equal(t, &value, tree.BTree.Get(orderedItem[int, string]{key, &value}).(orderedItem[int, string]).value)
}

func TestOrderedBTree_Get_ReturnsValueForGivenKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1
	value := "value"
	tree.Set(key, &value)

	// test
	result, ok := tree.Get(key)

	// assert
	assert.True(t, ok)
	assert.Equal(t, &value, result)
}

func TestOrderedBTree_Get_ReturnsFalseIfKeyNotFound(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1

	// test
	result, ok := tree.Get(key)

	// assert
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestOrderedBTree_Delete_DeletesKeyAndReturnsValueForGivenKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1
	value := "value"
	tree.Set(key, &value)

	// test
	result, ok := tree.Delete(key)

	// assert
	assert.True(t, ok)
	assert.Equal(t, &value, result)
	assert.False(t, tree.Has(key))
}

func TestOrderedBTree_DeleteMin_DeletesMinimumKeyAndReturnsValueForGivenKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)

	// test
	resultKey, resultValue := tree.DeleteMin()

	// assert
	assert.Equal(t, key1, resultKey)
	assert.Equal(t, &value1, resultValue)
	assert.False(t, tree.Has(key1))
}

func TestOrderedBTree_DeleteMax_DeletesMaximumKeyAndReturnsValueForGivenKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)

	// test
	resultKey, resultValue := tree.DeleteMax()

	// assert
	assert.Equal(t, key2, resultKey)
	assert.Equal(t, &value2, resultValue)
	assert.False(t, tree.Has(key2))
}

func TestOrderedBTree_Delete_ReturnsFalseIfKeyNotFound(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1

	// test
	result, ok := tree.Delete(key)

	// assert
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestOrderedBTree_Has_ReturnsTrueIfKeyExists(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1
	value := "value"
	tree.Set(key, &value)

	// test
	result := tree.Has(key)

	// assert
	assert.True(t, result)
}

func TestOrderedBTree_Has_ReturnsFalseIfKeyDoesNotExist(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1

	// test
	result := tree.Has(key)

	// assert
	assert.False(t, result)
}

func TestOrderedBTree_Len_ReturnsLengthOfBTree(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1
	value := "value"
	tree.Set(key, &value)

	// test
	result := tree.Len()

	// assert
	assert.Equal(t, 1, result)
}

func TestOrderedBTree_Min_ReturnsMinimumKeyAndValue(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1
	value := "value"
	tree.Set(key, &value)

	// test
	resultKey, resultValue := tree.Min()

	// assert
	assert.Equal(t, key, resultKey)
	assert.Equal(t, &value, resultValue)
}

func TestOrderedBTree_Min_ReturnsZeroValuesIfBTreeIsEmpty(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()

	// test
	resultKey, resultValue := tree.Min()

	// assert
	assert.Zero(t, resultKey)
	assert.Nil(t, resultValue)
}

func TestOrderedBTree_Max_ReturnsMaximumKeyAndValue(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1
	value := "value"
	tree.Set(key, &value)

	// test
	resultKey, resultValue := tree.Max()

	// assert
	assert.Equal(t, key, resultKey)
	assert.Equal(t, &value, resultValue)
}

func TestOrderedBTree_Max_ReturnsZeroValuesIfBTreeIsEmpty(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()

	// test
	resultKey, resultValue := tree.Max()

	// assert
	assert.Zero(t, resultKey)
	assert.Nil(t, resultValue)
}

func TestOrderedBTree_Set_MultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"

	// test
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	// assert
	assert.Equal(t, &value1, tree.BTree.Get(orderedItem[int, string]{key1, &value1}).(orderedItem[int, string]).value)
	assert.Equal(t, &value2, tree.BTree.Get(orderedItem[int, string]{key2, &value2}).(orderedItem[int, string]).value)
	assert.Equal(t, &value3, tree.BTree.Get(orderedItem[int, string]{key3, &value3}).(orderedItem[int, string]).value)
}

func TestOrderedBTree_Get_MultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	// test
	result1, ok1 := tree.Get(key1)
	result2, ok2 := tree.Get(key2)
	result3, ok3 := tree.Get(key3)

	// assert
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, &value1, result1)
	assert.Equal(t, &value2, result2)
	assert.Equal(t, &value3, result3)
}

func TestOrderedBTree_Delete_MultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	// test
	result1, ok1 := tree.Delete(key1)
	result2, ok2 := tree.Delete(key2)
	result3, ok3 := tree.Delete(key3)

	// assert
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.True(t, ok3)
	assert.Equal(t, &value1, result1)
	assert.Equal(t, &value2, result2)
	assert.Equal(t, &value3, result3)
	assert.False(t, tree.Has(key1))
	assert.False(t, tree.Has(key2))
	assert.False(t, tree.Has(key3))
}

func TestOrderedBTree_Has_MultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	// test
	result1 := tree.Has(key1)
	result2 := tree.Has(key2)
	result3 := tree.Has(key3)

	// assert
	assert.True(t, result1)
	assert.True(t, result2)
	assert.True(t, result3)
}

func TestOrderedBTree_Len_MultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	// test
	result := tree.Len()

	// assert
	assert.Equal(t, 3, result)
}

func TestOrderedBTree_Min_MultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	// test
	resultKey, resultValue := tree.Min()

	// assert
	assert.Equal(t, key1, resultKey)
	assert.Equal(t, &value1, resultValue)
}

func TestOrderedBTree_Max_MultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	// test
	resultKey, resultValue := tree.Max()

	// assert
	assert.Equal(t, key3, resultKey)
	assert.Equal(t, &value3, resultValue)
}

func TestOrderedBTree_Set_OverwritesValueForGivenKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key := 1
	value1 := "value1"
	value2 := "value2"
	tree.Set(key, &value1)

	// test
	tree.Set(key, &value2)

	// assert
	assert.Equal(t, &value2, tree.BTree.Get(orderedItem[int, string]{key, &value2}).(orderedItem[int, string]).value)
}

func TestOrderedBTree_Ascend_IteratesOnValuesAscending(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)
	expected := []string{value1, value2, value3}
	result := []string{}

	// test
	tree.Ascend(func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_Ascend_IteratesOnValuesAscendingAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)
	expected := []string{value1}
	result := []string{}

	// test
	tree.Ascend(func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_Descend_IteratesOnValuesDescending(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)
	expected := []string{value3, value2, value1}
	result := []string{}

	// test
	tree.Descend(func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_Descend_IteratesOnValuesDescendingAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"
	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)
	expected := []string{value3}
	result := []string{}

	// test
	tree.Descend(func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DeleteMin_DeletesMinimumKeyAndValueAndReturnsFalseIfBTreeIsEmpty(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()

	// test
	resultKey, resultValue := tree.DeleteMin()

	// assert
	assert.Zero(t, resultKey)
	assert.Nil(t, resultValue)
}

func TestOrderedBTree_DeleteMax_DeletesMaximumKeyAndValueAndReturnsFalseIfBTreeIsEmpty(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()

	// test
	resultKey, resultValue := tree.DeleteMax()

	// assert
	assert.Zero(t, resultKey)
	assert.Nil(t, resultValue)
}

func TestOrderedBTree_AscendRange_IteratesOnValuesAscendingInRange(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)
	expected := []string{value1, value2}
	result := []string{}

	// test
	tree.AscendRange(key1, key3, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendRange_IteratesOnValuesAscendingInRangeAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.AscendRange(key1, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendRange_IteratesOnValuesDescendingInRange(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.DescendRange(key1, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendRange_IteratesOnValuesDescendingInRangeAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.DescendRange(key1, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendLessThan_IteratesOnValuesAscendingLessThanKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.AscendLessThan(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendLessThan_IteratesOnValuesAscendingLessThanKeyAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.AscendLessThan(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendLessOrEqual_IteratesOnValuesDescendingLessThanKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{value1}
	result := []string{}

	// test
	tree.DescendLessOrEqual(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendLessOrEqual_IteratesOnValuesDescendingLessThanKeyAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{value1}
	result := []string{}

	// test
	tree.DescendLessOrEqual(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendGreaterOrEqual_IteratesOnValuesAscendingGreaterThanKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{value1}
	result := []string{}

	// test
	tree.AscendGreaterOrEqual(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendGreaterOrEqual_IteratesOnValuesAscendingGreaterThanKeyAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{value1}
	result := []string{}

	// test
	tree.AscendGreaterOrEqual(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendGreaterThan_IteratesOnValuesDescendingGreaterThanKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.DescendGreaterThan(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendGreaterThan_IteratesOnValuesDescendingGreaterThanKeyAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.DescendGreaterThan(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendRange_IteratesOnValuesAscendingInRangeWithMultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.AscendRange(key1, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendRange_IteratesOnValuesAscendingInRangeWithMultipleKeysAndValuesAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.AscendRange(key1, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendRange_IteratesOnValuesDescendingInRangeWithMultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value3, value2}
	result := []string{}

	// test
	tree.DescendRange(key3, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendRange_IteratesOnValuesDescendingInRangeWithMultipleKeysAndValuesAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value3}
	result := []string{}

	// test
	tree.DescendRange(key3, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendLessThan_IteratesOnValuesAscendingLessThanKeyWithMultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value1, value2}
	result := []string{}

	// test
	tree.AscendLessThan(key3, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendLessThan_IteratesOnValuesAscendingLessThanKeyWithMultipleKeysAndValuesAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value1}
	result := []string{}

	// test
	tree.AscendLessThan(key3, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendLessOrEqual_IteratesOnValuesDescendingLessThanKeyWithMultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 3
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value3, value2, value1}
	result := []string{}

	// test
	tree.DescendLessOrEqual(key3, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendLessOrEqual_IteratesOnValuesDescendingLessThanKeyWithMultipleKeysAndValuesAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 1
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value3}
	result := []string{}

	// test
	tree.DescendLessOrEqual(key3, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendGreaterOrEqual_IteratesOnValuesAscendingGreaterThanKeyWithMultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 1
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value3, value2, value1}
	result := []string{}

	// test
	tree.AscendGreaterOrEqual(key3, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendGreaterOrEqual_IteratesOnValuesAscendingGreaterThanKeyWithMultipleKeysAndValuesAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 1
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value3}
	result := []string{}

	// test
	tree.AscendGreaterOrEqual(key3, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendGreaterThan_IteratesOnValuesDescendingGreaterThanKeyWithMultipleKeysAndValues(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 1
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value1, value2}
	result := []string{}

	// test
	tree.DescendGreaterThan(key3, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendGreaterThan_IteratesOnValuesDescendingGreaterThanKeyWithMultipleKeysAndValuesAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 1
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value1}
	result := []string{}

	// test
	tree.DescendGreaterThan(key3, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendRange_IteratesOnValuesAscendingInRangeWithMultipleKeysAndValuesAndNoValuesInRange(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 1
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value3, value2}
	result := []string{}

	// test
	tree.AscendRange(key3, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendRange_IteratesOnValuesAscendingInRangeWithMultipleKeysAndValuesAndNoValuesInRangeAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 1
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{value3}
	result := []string{}

	// test
	tree.AscendRange(key3, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendRange_IteratesOnValuesDescendingInRangeWithMultipleKeysAndValuesAndNoValuesInRange(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 1
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)

	expected := []string{}
	result := []string{}

	// test
	tree.DescendRange(key3, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendRange_IteratesOnValuesDescendingInRangeWithMultipleKeysAndValuesAndNoValuesInRangeAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"
	key2 := 2
	value2 := "value2"
	key3 := 1
	value3 := "value3"

	tree.Set(key1, &value1)
	tree.Set(key2, &value2)
	tree.Set(key3, &value3)
	expected := []string{}
	result := []string{}

	// test
	tree.DescendRange(key3, key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendLessThan_IteratesOnValuesAscendingLessThanKeyWithMultipleKeysAndValuesAndNoValuesLessThanKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.AscendLessThan(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendLessThan_IteratesOnValuesAscendingLessThanKeyWithMultipleKeysAndValuesAndNoValuesLessThanKeyAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{}
	result := []string{}

	// test
	tree.AscendLessThan(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendLessOrEqual_IteratesOnValuesDescendingLessThanKeyWithMultipleKeysAndValuesAndNoValuesLessThanKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{value1}
	result := []string{}

	// test
	tree.DescendLessOrEqual(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_DescendLessOrEqual_IteratesOnValuesDescendingLessThanKeyWithMultipleKeysAndValuesAndNoValuesLessThanKeyAndStopsWhenFuncReturnsFalse(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 3
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{value1}
	result := []string{}

	// test
	tree.DescendLessOrEqual(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return false
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_AscendGreaterOrEqual_IteratesOnValuesAscendingGreaterThanKeyWithMultipleKeysAndValuesAndNoValuesGreaterThanKey(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)
	expected := []string{value1}
	result := []string{}

	// test
	tree.AscendGreaterOrEqual(key1, func(key int, value *string) bool {
		result = append(result, *value)
		return true
	})

	// assert
	assert.Equal(t, expected, result)
}

func TestOrderedBTree_Clone(t *testing.T) {
	// setup
	tree := NewOrderedBTree[int, string]()
	key1 := 1
	value1 := "value1"

	tree.Set(key1, &value1)

	// test
	result := tree.Clone()

	// assert
	assert.True(t, result.Has(key1))
}
