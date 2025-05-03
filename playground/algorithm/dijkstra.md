# Dijkstra's Algorithm

Dijkstra's algorithm is a famous algorithm for finding the shortest paths between nodes in a graph, which may represent, for example, road networks. It is widely used in routing and navigation systems.

## How It Works

1. Assign to every node a tentative distance value: set it to zero for the initial node and to infinity for all other nodes.
2. Set the initial node as current. Mark all other nodes unvisited.
3. For the current node, consider all of its unvisited neighbors and calculate their tentative distances through the current node. Compare the newly calculated tentative distance to the current assigned value and assign the smaller one.
4. After considering all neighbors of the current node, mark the current node as visited. A visited node will not be checked again.
5. If the destination node has been marked visited or if the smallest tentative distance among the nodes is infinity, the algorithm stops.
6. Otherwise, select the unvisited node with the smallest tentative distance, set it as the new current node, and repeat from step 3.

## Time Complexity

- **Using a simple array:** O(V²)
- **Using a priority queue (min-heap):** O((V + E) log V)

Where:
- V is the number of vertices
- E is the number of edges

## Example in Python

```python
import heapq

def dijkstra(graph, start):
    queue = []
    heapq.heappush(queue, (0, start))
    distances = {node: float('inf') for node in graph}
    distances[start] = 0

    while queue:
        current_distance, current_node = heapq.heappop(queue)

        if current_distance > distances[current_node]:
            continue

        for neighbor, weight in graph[current_node].items():
            distance = current_distance + weight
            if distance < distances[neighbor]:
                distances[neighbor] = distance
                heapq.heappush(queue, (distance, neighbor))

    return distances

# Example usage:
graph = {
    'A': {'B': 1, 'C': 4},
    'B': {'A': 1, 'C': 2, 'D': 5},
    'C': {'A': 4, 'B': 2, 'D': 1},
    'D': {'B': 5, 'C': 1}
}

distances = dijkstra(graph, 'A')
print(distances)
```

**Output:**
```
{'A': 0, 'B': 1, 'C': 3, 'D': 4}
```

## Visual Illustration

Given this graph:

```
A --1-- B
|      / |
4    2  5
|  /    |
C --1-- D
```

Finding shortest paths from A:

- A → B = 1
- A → C = 3 (A → B → C)
- A → D = 4 (A → B → C → D)

## When to Use Dijkstra's Algorithm

- When you need the shortest path in a weighted graph with non-negative weights.
- In routing, network analysis, and pathfinding for games.

## Conclusion

Dijkstra's algorithm is a cornerstone technique in graph theory and computer science. It is an essential algorithm to know, especially when dealing with shortest-path problems in weighted graphs.