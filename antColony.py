import numpy as np
import random
from copy import deepcopy

def SolveTSPWithAntAlgorithm(tspMatrix, iterations=100, numAnts=10, alpha=1, beta=2, evaporationRate=0.1, Q=100):
    n = len(tspMatrix)
    INF = float('inf')
    
    pheromone = [[1.0 for _ in range(n)] for _ in range(n)]
    maxTspValue = max(max(val for val in row if val != INF) for row in tspMatrix)
    
    cMatrix = [[0] * n for _ in range(n)]
    for i in range(n):
        for j in range(n):
            tspValue = tspMatrix[i][j]
            if tspValue == INF:
                tspValue = 10 * maxTspValue
            cMatrix[i][j] = tspValue
    
    bestCost = float('inf')
    bestAssignment = None
    
    for iteration in range(iterations):
        all_solutions = []
        
        for ant in range(numAnts):
            visited = set()
            currentCity = random.randint(0, n-1)
            visited.add(currentCity)
            selectedWorks = [currentCity]
            
            while len(selectedWorks) < n:
                probabilities = [0.0] * n
                total = 0.0
                
                for nextCity in range(n):
                    if nextCity in visited or nextCity == currentCity:
                        probabilities[nextCity] = 0.0
                    else:
                        prob = (pheromone[currentCity][nextCity] ** alpha) * \
                              ((1.0 / cMatrix[currentCity][nextCity]) ** beta)
                        probabilities[nextCity] = prob
                        total += prob
                
                if total > 0:
                    probabilities = [p/total for p in probabilities]
                    nextCity = random.choices(range(n), weights=probabilities, k=1)[0]
                else:
                    for city in range(n):
                        if city not in visited and city != currentCity:
                            nextCity = city
                            break
                
                selectedWorks.append(nextCity)
                visited.add(nextCity)
                currentCity = nextCity
            
            C_selected = [[0] * n for _ in range(n)]
            for i in range(n):
                for j in range(n):
                    C_selected[i][j] = cMatrix[selectedWorks[i]][selectedWorks[j]]
            
            rawAssignment = solve_assignment(C_selected)     

            filteredAssignment = rawAssignment.copy()
            assigned = set(filteredAssignment)
            
            for j in range(n):
                i_assigned = filteredAssignment[j]
                if i_assigned == j:
                    for k in range(n):
                        if k != j and k not in assigned:
                            i_assigned = k
                            break
                    filteredAssignment[j] = i_assigned
                    assigned.add(i_assigned)

            totalCost = sum(C_selected[i][j] for j, i in enumerate(filteredAssignment))
            
            all_solutions.append((totalCost, filteredAssignment, selectedWorks))
            
            if totalCost < bestCost:
                bestCost = totalCost
                bestAssignment = filteredAssignment

        for i in range(n):
            for j in range(n):
                pheromone[i][j] *= (1 - evaporationRate)
        
        for solution in all_solutions:
            cost, assignment, works = solution
            if cost == 0: 
                continue
            for j in range(n):
                i_assigned = assignment[j]
                pheromone[works[i_assigned]][works[j]] += Q / cost
    
    return bestCost, bestAssignment


def solve_assignment(cost_matrix):
    n = len(cost_matrix)
    assignment = []
    used_rows = set()
    
    for j in range(n):
        min_cost = float('inf')
        best_i = -1
        for i in range(n):
            if i not in used_rows and cost_matrix[i][j] < min_cost:
                min_cost = cost_matrix[i][j]
                best_i = i
        if best_i == -1:
            for i in range(n):
                if i not in used_rows:
                    best_i = i
                    break
        assignment.append(best_i)
        used_rows.add(best_i)
    return assignment
