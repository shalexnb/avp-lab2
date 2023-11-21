package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Graph struct {
	Vertices map[int]map[int]int
}

type Path struct {
	nodes []int
	cost  int
}

func (g *Graph) addEdge(start, end, cost int) {
	if g.Vertices == nil {
		g.Vertices = make(map[int]map[int]int)
	}
	if g.Vertices[start] == nil {
		g.Vertices[start] = make(map[int]int)
	}
	g.Vertices[start][end] = cost
}

//для каждого соседа начальной вершины создаются новые горутины для параллельного выполнения рекурсивного поиска пути 
//выполняется поиск в глубину
func findShortestPath(graph Graph, start, end int, wg *sync.WaitGroup, resultChan chan Path) {
	defer wg.Done()

	visited := make(map[int]bool)
	currentPath := Path{nodes: []int{start}, cost: 0}

	findShortestPathRecursive(graph, start, end, visited, currentPath, wg, resultChan)
}

//здесь функция рекурсивно вызывает саму себя для каждого соседа данной вершины, не посещенного ранее
func findShortestPathRecursive(graph Graph, current, end int, visited map[int]bool, currentPath Path, wg *sync.WaitGroup, resultChan chan Path) {
	if current == end {
		resultChan <- currentPath
		return
	}

	visited[current] = true

	for neighbor, cost := range graph.Vertices[current] {
		if !visited[neighbor] {
			wg.Add(1)
			go findShortestPathRecursive(graph, neighbor, end, visited, Path{append(currentPath.nodes, neighbor), currentPath.cost + cost}, wg, resultChan)
		}
	}

	//после обхода соседей текущей вершины отменяется отметка о посещении этой вершины 
	//чтобы при возврате из рекурсии мы могли исключить текущую вершину из пути
	visited[current] = false
}

func readInput(prompt string) int {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Неверный ввод. Пожалуйста, введите целое число.")
		return readInput(prompt)
	}
	return input
}

func main() {
	graph := Graph{}

	fmt.Println("Введите вершины и ребра графа")
	fmt.Println("(Чтобы прекратить добавление ребер, введите 0 в качестве начальной вершины)")

	for {
		start := readInput("Введите начальную вершину (введите 0 чтобы оставновиться): ")
		if start == 0 {
			break
		}

		end := readInput("Введите конечную вершину: ")
		cost := readInput("Введите вес ребра: ")

		graph.addEdge(start, end, cost)
	}

	fmt.Println("\nВведите начальную и конечную вершины для нахождения кратчайшего пути:")
	start := readInput("Начальная вершина: ")
	end := readInput("Конечная вершина: ")

	var wg sync.WaitGroup
	resultChan := make(chan Path)

	wg.Add(1)
	go findShortestPath(graph, start, end, &wg, resultChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	shortestPath := <-resultChan
	fmt.Printf("Кратчайший путь от %d до %d: %v\n", start, end, shortestPath.nodes)
	fmt.Printf("Вес пути: %d\n", shortestPath.cost)
}
