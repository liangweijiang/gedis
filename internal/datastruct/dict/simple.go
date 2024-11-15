package dict

type SimpleDict struct {
	m map[string]interface{}
}

func NewSimpleDict() *SimpleDict {
	return &SimpleDict{
		m: make(map[string]interface{}),
	}
}

func (d *SimpleDict) Put(key string, val interface{}) (result int) {
	_, exists := d.m[key]
	d.m[key] = val
	if !exists {
		result = 1
	}
	return
}

func (d *SimpleDict) Get(key string) (val interface{}, exists bool) {
	val, exists = d.m[key]
	return
}

func (d *SimpleDict) Remove(key string) (val interface{}, result int) {
	val, exists := d.m[key]
	if exists {
		delete(d.m, key)
		result = 1
	}
	return
}

func (d *SimpleDict) Foreach(f func(key string, val interface{}) bool) {
	for k, v := range d.m {
		if !f(k, v) {
			break
		}
	}
}

func (d *SimpleDict) Keys() []string {
	keys := make([]string, len(d.m))
	i := 0
	for k := range d.m {
		keys[i] = k
		i++
	}
	return keys
}

func (d *SimpleDict) Len() int {
	return len(d.m)
}

func (d *SimpleDict) Clear() {
	d.m = make(map[string]interface{})
}
