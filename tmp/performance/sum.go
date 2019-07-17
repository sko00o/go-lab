package sum

//测试数据
var testData = CreateTestData()

func GetValue(index int) int {
	return testData[index]
}

//获取测试数据长度
func GetDataLen() int {
	return len(testData)
}

//创建一个6000000000大小的整型切片
func CreateTestData() []int {
	data := make([]int, 600000000)
	for index := range data {
		data[index] = index % 128
	}
	return data
}

func sum(result *int) {
	for i := 0; i < GetDataLen(); i++ {
		*result += GetValue(i)
	}
}

func GetData() []int {
	return testData
}

// 指令级优化： 消除连续的函数调用和不必要的存储器引用
func sum1(result *int) {
	k := *result
	data := GetData()
	dataLength := GetDataLen()
	for i := 0; i < dataLength; i++ {
		k += data[i]
	}
	*result = k
}

// 低级优化： 循环展开，提高并行性
func sum2(result *int) {
	k := *result
	k1 := 0
	data := GetData()
	dataLength := GetDataLen()

	for i := 1; i < dataLength; i += 2 {
		k += data[i]
		k1 += data[i-1]
	}
	if dataLength%2 == 1 {
		k += data[dataLength-1]
	}
	*result = k + k1
}
