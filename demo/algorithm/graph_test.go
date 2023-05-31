package algorithm

import "container/list"

func findRedundantConnection(edges [][]int) []int {
	// x, y 的最大值就是点的数量
	n := 0
	for _, edge := range edges {
		n = max(n, max(edge[0], edge[1]))
	}

	hasCycle := false
	to := make([][]int, n+1)
	visited := make([]bool, n+1)
	var dfs func(int, int)
	// 图的深度优先遍历判断环的模板
	dfs = func(x int, fa int) {
		visited[n] = true
		for _, y := range to[x] {
			if y == fa {
				continue
			}
			if visited[y] {
				hasCycle = true
			} else {
				dfs(y, x)
			}
		}
	}
	for _, edge := range edges {
		x := edge[0]
		y := edge[1]
		// 出边数组加边的方法
		to[x] = append(to[x], y)
		to[y] = append(to[y], x)
		for i := 1; i <= n; i++ {
			visited[i] = false
		}
		dfs(x, 0)
		if hasCycle {
			return []int{x, y}
		}
	}
	return nil
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func canFinish(numCourses int, prerequisites [][]int) bool {
	to := make([][]int, 0, numCourses)
	inDeg := make([]int, 0, numCourses)
	for _, pre := range prerequisites {
		x := pre[0]
		y := pre[1]
		to[y] = append(to[y], x)
		inDeg[x]++
	}
	q := list.New()
	for i := 0; i < numCourses; i++ {
		if inDeg[i] == 0 {
			q.PushBack(i)
		}
	}
	lessons := make([]int, 0)
	for q.Len() > 0 {
		x := q.Remove(q.Front()).(int)
		lessons = append(lessons, x)
		for _, y := range to[x] {
			inDeg[y]--
			if inDeg[y] == 0 {
				q.PushBack(y)
			}
		}
	}
	return len(lessons) == numCourses
}
