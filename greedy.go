package main

import (
    "fmt"
    "math"
    "sort"
)

const INF = math.MaxInt32

func SolveTSPViaTwoLevelAssignment(tspMatrix [][]int) (int, []int) {
    n := len(tspMatrix)
    if n == 0 {
        return 0, nil
    }

    maxTspValue := 0
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if tspMatrix[i][j] != INF && tspMatrix[i][j] > maxTspValue {
                maxTspValue = tspMatrix[i][j]
            }
        }
    }

    cMatrix := make([][]int, n)
    for i := 0; i < n; i++ {
        cMatrix[i] = make([]int, n)
        for j := 0; j < n; j++ {
            tspVal := tspMatrix[i][j]
            if tspVal == INF {
                tspVal = 10 * maxTspValue
            }
            cMatrix[i][j] = tspVal
        }
    }

    workScores := make([]int, n)
    for i := 0; i < n; i++ {
        score := 0
        for j := 0; j < n; j++ {
            if i != j {
                score += cMatrix[i][j]
            }
        }
        workScores[i] = score
    }

    indexes := make([]int, n)
    for i := range indexes {
        indexes[i] = i
    }

    sort.Slice(indexes, func(i, j int) bool {
        return workScores[indexes[i]] < workScores[indexes[j]]
    })

    selectedWorks := indexes[:n]

    cSelected := make([][]int, n)
    for i := 0; i < n; i++ {
        cSelected[i] = make([]int, n)
        for j := 0; j < n; j++ {
            cSelected[i][j] = cMatrix[selectedWorks[i]][selectedWorks[j]]
        }
    }

    rawAssignment := HungarianMethod(cSelected)

    filteredAssignment := make([]int, n)
    copy(filteredAssignment, rawAssignment)

    assigned := make([]bool, n)
    for j := 0; j < n; j++ {
        iAssigned := filteredAssignment[j]
        if iAssigned == j {
            for k := 0; k < n; k++ {
                if k != j && !assigned[k] {
                    iAssigned = k
                    assigned[k] = true
                    break
                }
            }
            filteredAssignment[j] = iAssigned
        } else {
            assigned[iAssigned] = true
        }
    }

    totalCost := 0
    for j := 0; j < n; j++ {
        iAssigned := filteredAssignment[j]
        totalCost += cSelected[iAssigned][j]
    }

    return totalCost, filteredAssignment
}


func HungarianMethod(matrix [][]int) []int {
    n := len(matrix)
    m := len(matrix[0])

    u := make([]int, n+1)
    v := make([]int, m+1)
    p := make([]int, m+1)
 way := make([]int, m+1)

    for i := 1; i <= n; i++ {
        p[0] = i
        minv := make([]int, m+1)
        used := make([]bool, m+1)
        for j := 0; j <= m; j++ {
            minv[j] = INF
            used[j] = false
        }

        f := 0
        var j0 int
        for {
            used[j0] = true
            i0 := p[j0]
            delta := INF
            j1 := 0
            for j := 1; j <= m; j++ {
                if !used[j] {
                    cur := matrix[i0-1][j-1] - u[i0] - v[j]
                    if cur < minv[j] {
                        minv[j] = cur
                        way[j] = j0
                    }
                    if minv[j] < delta {
                        delta = minv[j]
                        j1 = j
                    }
                }
            }
            for j := 0; j <= m; j++ {
                if used[j] {
                    u[p[j]] += delta
                    v[j] -= delta
                } else {
                    minv[j] -= delta
                }
            }
            j0 = j1
            if p[j0] == 0 {
                break
            }
        }
        for {
            j1 := way[j0]
            p[j0] = p[j1]
            j0 = j1
            if j0 == 0 {
                break
            }
        }
    }

    result := make([]int, n)
    for j := 1; j <= m; j++ {
        if p[j] != 0 {
            result[p[j]-1] = j - 1
        }
    }
    return result
}
