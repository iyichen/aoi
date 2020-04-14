package util

import "aoi/manager"

// 比较两个int类型切片内的元素是否相同
func IntEqual(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for _, a1 := range slice1 {
		find := false
		for _, a2 := range slice2 {
			if a1 == a2 {
				find = true
				break
			}
		}
		if !find {
			return false
		}
	}
	return true
}

// 获取两个int类型切片的交集
func IntIntersect(slice1, slice2 []int) (res []int) {
	for _, e1 := range slice1 {
		find := false
		for _, e2 := range slice2 {
			if e1 == e2 {
				find = true
				break
			}
		}
		if find {
			res = append(res, e1)
		}
	}
	return res
}

// 获取两个int类型切片的差集
func IntExcept(slice1, slice2 []int) (res []int) {
	for _, e1 := range slice1 {
		find := false
		for _, e2 := range slice2 {
			if e1 == e2 {
				find = true
				break
			}
		}
		if !find {
			res = append(res, e1)
		}
	}
	return res
}

func PlayerIntersect(ps1, ps2 []*manager.Player) (res []*manager.Player) {
	for _, e1 := range ps1 {
		find := false
		for _, e2 := range ps2 {
			if e1.Id == e2.Id {
				find = true
				break
			}
		}
		if find {
			res = append(res, e1)
		}
	}
	return res
}

func PlayerExcept(ps1, ps2 []*manager.Player) (res []*manager.Player) {
	for _, e1 := range ps1 {
		find := false
		for _, e2 := range ps2 {
			if e1.Id == e2.Id {
				find = true
				break
			}
		}
		if !find {
			res = append(res, e1)
		}
	}
	return res
}
