package test

import "fmt"

func TestSlice() {
	arr01 := make([]int, 2, 2) // 长度为2，容量为2
	arr01[0] = 1
	arr01[1] = 2
	arr01 = append(arr01, 1)
	arr01 = append(arr01, 2)
	arr01 = append(arr01, 3)
	arr01 = append(arr01, 4)
	arr01 = append(arr01, 5)

	var numbers4 = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	myslice := numbers4[4:6:8] // 从4到6的元素，容量到第八个元素
	fmt.Printf("myslice为 %d, 其长度为: %d\n", myslice, len(myslice))

	myslice = myslice[:cap(myslice)]
	fmt.Printf("myslice的第四个元素为: %d", myslice[3])

	fmt.Printf("arr01:%d %d %d %T\n", arr01[0], len(arr01), cap(arr01), arr01)
}
