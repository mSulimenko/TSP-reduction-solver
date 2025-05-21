package main

import (
    "fmt"
    "math"
    "math/rand"
    "sort"
    "time"
)

const INF = math.MaxInt32
const Q = 1000

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
            if delta == INF {
                break
            }
            for j := 0; j <= m; j++ {
                if used[j] {
                    u[p[j]] += delta
                } else {
                    minv[j] -= delta
                }
                v[j] -= delta
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

func SolveTSPWithAntAlgorithm(tspMatrix [][]int, iterations, numAnts int, alpha, beta, evaporationRate float64) (int, []int) {
    rand.Seed(time.Now().UnixNano())
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
            if tspMatrix[i][j] == INF {
                cMatrix[i][j] = 10 * maxTspValue
            } else {
                cMatrix[i][j] = tspMatrix[i][j]
            }
        }
    }

    pheromone := make([][]float64, n)
    for i := 0; i < n; i++ {
        pheromone[i] = make([]float64, n)
        for j := 0; j < n; j++ {
            pheromone[i][j] = 1.0
        }
    }

    bestCost := INF
    var bestAssignment []int

    for iter := 0; iter < iterations; iter++ {
        allSolutions := make([]struct {
            cost         int
            assignment   []int
            selectedWorks []int
        }, numAnts)

        for ant := 0; ant < numAnts; ant++ {
            visited := make(map[int]bool)
            currentCity := rand.Intn(n)
            visited[currentCity] = true
            selectedWorks := []int{currentCity}

            for len(selectedWorks) < n {
                probabilities := make([]float64, n)
                total := 0.0

                for nextCity := 0; nextCity < n; nextCity++ {
                    if visited[nextCity] || nextCity == currentCity {
                        probabilities[nextCity] = 0
                    } else {
                        pher := math.Pow(pheromone[currentCity][nextCity], alpha)
                        dist := math.Pow(1.0/float64(cMatrix[currentCity][nextCity]), beta)
                        prob := pher * dist
                        probabilities[nextCity] = prob
                        total += prob
                    }
                }

                if total > 0 {
                    for nextCity := 0; nextCity < n; nextCity++ {
                        probabilities[nextCity] /= total
                    }
                }

                r := rand.Float64()
                sumProb := 0.0
                nextCity := currentCity
                for i := 0; i < n; i++ {
                    sumProb += probabilities[i]
                    if sumProb >= r {
                        nextCity = i
                        break
                    }
                }

                visited[nextCity] = true
                selectedWorks = append(selectedWorks, nextCity)
                currentCity = nextCity
            }

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
                }
                filteredAssignment[j] = iAssigned
                assigned[iAssigned] = true
            }

            totalCost := 0
            for j := 0; j < n; j++ {
                iAssigned := filteredAssignment[j]
                totalCost += cSelected[iAssigned][j]
            }

            allSolutions[ant] = struct {
                cost         int
                assignment   []int
                selectedWorks []int
            }{totalCost, filteredAssignment, selectedWorks}

            if totalCost < bestCost {
                bestCost = totalCost
                bestAssignment = filteredAssignment
            }
        }

        for i := 0; i < n; i++ {
            for j := 0; j < n; j++ {
                pheromone[i][j] *= (1 - evaporationRate)
            }
        }

        for ant := 0; ant < numAnts; ant++ {
            sol := allSolutions[ant]
            cost := sol.cost
            assignment := sol.assignment
            works := sol.selectedWorks

            if cost == 0 {
                continue
            }
            deposit := float64(Q) / float64(cost)

            for j := 0; j < n; j++ {
                iAssigned := assignment[j]
                if iAssigned != j {
                    from := works[iAssigned]
                    to := works[j]
                    pheromone[from][to] += deposit
                }
            }
        }
    }

    return bestCost, bestAssignment
}
