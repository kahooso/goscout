package twosum

func TwoSum(lst []int, target int) (int, int, bool) {
	m := make(map[int]int)
	for i, v := range lst {
		complement := target - v
		if idx, ok := m[complement]; ok {
			return idx, i, true
		}
		m[v] = i
	}
	return 0, 0, false
}
