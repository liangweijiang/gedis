package datastruct

// Dict 接口定义了一个简单的键值对存储结构的操作方法。
// 它允许用户存储和检索与字符串键关联的值。
type Dict interface {
	// Get 方法根据键获取对应的值。
	// 它接受一个字符串类型的键，返回值和一个布尔值指示键是否存在。
	// 如果键存在，返回对应的值和true；如果键不存在，返回nil和false。
	Get(key string) (val interface{}, exists bool)

	// Len 方法返回当前存储结构中键值对的数量。
	Len() int

	// Put 方法向存储结构中插入一个键值对。
	// 如果键已存在，它会用新值覆盖旧值，并返回旧值。
	// 参数key是字符串类型的键，val是任意类型的值。
	// 返回值result是被覆盖的旧值的类型。
	Put(key string, val interface{}) (result int)

	// Keys 方法返回当前存储结构中所有键的列表。
	// 返回值是一个字符串切片，包含所有的键。
	Keys() []string

	// Remove 方法从存储结构中移除指定键的键值对。
	// 参数key是字符串类型的键。
	// 返回值val是被移除的值，result表示移除操作的结果。
	Remove(key string) (val interface{}, result int)

	// Clear 方法清空存储结构中的所有键值对。
	Clear()

	// Foreach 方法遍历存储结构中的所有键值对，并对每个键值对执行给定的函数f。
	// 函数f接受键和值作为参数，并返回一个布尔值。
	// 如果f返回false，遍历会提前终止。
	Foreach(f func(key string, val interface{}) bool)
}
