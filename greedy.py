def SolveTSPViaTwoLevelAssignment(tspMatrix):
    n = len(tspMatrix)
    maxTspValue = 0
    INF = float('inf')
  
    for i in range(n):
        for j in range(n):
            if tspMatrix[i][j] != INF and tspMatrix[i][j] > maxTspValue:
                maxTspValue = tspMatrix[i][j]
    
    cMatrix = [[0] * n for _ in range(n)]
    for i in range(n):
        for j in range(n):
            tspValue = tspMatrix[i][j]
            if tspValue == INF:
                tspValue = 10 * maxTspValue
            cMatrix[i][j] = tspValue

      
    workScores = [sum(cMatrix[i][j] for j in range(n) if j != i) for i in range(n)]
    
    indices = sorted(range(n), key=lambda x: workScores[x])
    selectedWorks = indices  

    c_selected = [[0] * n for _ in range(n)]
    for i in range(n):
        for j in range(n):
            c_selected[i][j] = cMatrix[selectedWorks[i]][selectedWorks[j]]
    
    rawAssignment = greedy_solve_assignment(c_selected)
    
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
    

    totalCost = sum(c_selected[i][j] for j, i in enumerate(filteredAssignment))
    
    return totalCost, filteredAssignment

def greedy_solve_assignment(cost_matrix):
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
