package main

func Map() {
	m1 := map[string]string{
		"key1": "value1",
	}

	val, ok := m1["key1"]
	if ok {
		println("第一步:", val)
	}

	val = m1["key2"]
	println("第二步:", val)

	m2 := make(map[string]string, 4)
	m2["key2"] = "value2"

	println(len(m2))
	for k, v := range m1 {
		println(k, v)
	}

	for k := range m1 {
		println(k)
	}

	delete(m1, "key1")
}
