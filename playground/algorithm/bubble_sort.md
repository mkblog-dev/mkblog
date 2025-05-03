# Bubble Sort

Bubble Sort is one of the simplest sorting algorithms. It works by repeatedly swapping adjacent elements if they are in the wrong order. Although it's not suitable for large datasets due to its poor performance, it's a great starting point for understanding sorting.

## How It Works

1. Compare the first two elements.
2. If the first is greater than the second, swap them.
3. Move to the next pair and repeat step 2.
4. Continue this process for each element.
5. After each full pass through the list, the largest unsorted element will have "bubbled" up to its correct position.
6. Repeat the passes until no swaps are needed.

## Time Complexity

- **Best Case (Already Sorted):** O(n)
- **Average Case:** O(n²)
- **Worst Case:** O(n²)

## Example in Python

```python
def bubble_sort(arr):
    n = len(arr)
    for i in range(n):
        swapped = False
        for j in range(0, n - i - 1):
            if arr[j] > arr[j + 1]:
                arr[j], arr[j + 1] = arr[j + 1], arr[j]
                swapped = True
        if not swapped:
            break
    return arr

# Example usage:
numbers = [64, 34, 25, 12, 22, 11, 90]
sorted_numbers = bubble_sort(numbers)
print(f"Sorted array: {sorted_numbers}")
```

**Output:**
```
Sorted array: [11, 12, 22, 25, 34, 64, 90]
```

## Visual Illustration

Let's say we have the list `[5, 3, 8, 4, 2]`:

- First pass:
  - Compare 5 and 3 → Swap → [3, 5, 8, 4, 2]
  - Compare 5 and 8 → No swap
  - Compare 8 and 4 → Swap → [3, 5, 4, 8, 2]
  - Compare 8 and 2 → Swap → [3, 5, 4, 2, 8]
- Second pass:
  - Compare 3 and 5 → No swap
  - Compare 5 and 4 → Swap → [3, 4, 5, 2, 8]
  - Compare 5 and 2 → Swap → [3, 4, 2, 5, 8]
- Continue similarly until fully sorted.

## When to Use Bubble Sort

- For educational purposes to understand basic sorting concepts.
- When dealing with very small datasets.
- When the simplicity of code matters more than performance.

## Conclusion

Bubble Sort is easy to understand and implement, making it a great algorithm for beginners. However, for practical use with larger datasets, more efficient algorithms like quicksort or mergesort are preferred.