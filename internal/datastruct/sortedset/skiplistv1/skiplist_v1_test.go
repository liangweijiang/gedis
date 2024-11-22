package skiplistv1

import (
	"fmt"
	"testing"
)

func TestSlipList(t *testing.T) {
	// 创建一个新的跳表实例
	skipList := NewSkipList()

	// 定义一些用于测试的值和分数
	testValues := []interface{}{"zywOo", "niko", "hunter", "HooXi", "simple"}
	testScores := []float64{10, 20, 30, 40, 50}

	// 测试插入非nil值
	for i, val := range testValues {
		result := skipList.Insert(val, testScores[i])
		if result != 0 {
			t.Errorf("Insert returned %d instead of Success for value %v and score %f", result, val, testScores[i])
		}
		fmt.Println(skipList.String())
		fmt.Println()
		fmt.Println()
	}
	result := skipList.Insert("donk", 10)
	if result != 0 {
		fmt.Printf("Insert returned %d instead of Success\n", result)
	}
	fmt.Println(skipList.String())

	fmt.Println(skipList.Find(10))
	fmt.Println("len: ", skipList.length)
	fmt.Println(skipList.Delete(10))
	fmt.Println(skipList.String())
	fmt.Println(skipList.Find(10))
	fmt.Println("len: ", skipList.length)

	/*// 尝试插入一个已经存在的值
	duplicateResult := skipList.Insert(3, 30)
	if duplicateResult != ValueExists {
		t.Errorf("Insert should have returned ValueExists for duplicate value, got %d instead", duplicateResult)
	}

	// 尝试插入nil值
	nilResult := skipList.Insert(nil, 60)
	if nilResult != NilValue {
		t.Errorf("Insert should have returned NilValue for nil value, got %d instead", nilResult)
	}

	fmt.Println(skipList.String())*/
}
