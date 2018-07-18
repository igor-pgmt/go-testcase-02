package main

import (
	"container/heap"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

// graph это структура основного объекта - графа
type graph struct {
	rawMap   map[string]map[string]int
	pq       PriorityQueue
	fileName string
	fullMap  map[string]map[string]string
}

// NewGraph создаёт граф и импортирует в него файл по-умолчанию
func NewGraph() *graph {
	g := &graph{}
	g.importFile("cities")
	return g
}

// InsertVertex Добавляет в граф новые соседствующие вершины  и расстяние между ними
func (g graph) InsertVertex(name, adjacent string, distance int) {
	if _, ok := g.rawMap[name]; ok {
		g.rawMap[name][adjacent] = distance
	} else {
		g.rawMap[name] = map[string]int{adjacent: distance}
	}

	if _, ok := g.rawMap[adjacent]; ok {
		g.rawMap[adjacent][name] = distance
	} else {
		g.rawMap[adjacent] = map[string]int{name: distance}
	}
}

// generateMap на основании содержимого файла
// генерирует карту вершин графа
func (g *graph) generateMap() {
	g.rawMap = make(map[string]map[string]int)
	file, err := os.Open(g.fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ' '
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		distance, err := strconv.Atoi(record[2])
		if err != nil {
			panic(err)
		}
		g.InsertVertex(record[0], record[1], distance)
	}
}

// generateQueue создаёт структуру очереди с приоритетом
// используется в алгоритме Дейкстры
func (g *graph) generateQueue(start string) {
	g.pq = make(PriorityQueue, len(g.rawMap))
	i := 0
	for name, adjacents := range g.rawMap {
		g.pq[i] = &Item{
			name:      name,
			adjacents: adjacents,
			priority:  maxPriority,
			index:     i,
		}
		if name == start {
			g.pq[i].priority = uint32(0)
		}
		i++
	}
	heap.Init(&g.pq)
}

// findPaths создаёт карту кратчайших путей для каждой из вершин графа
func (g *graph) findPaths() {
	g.fullMap = make(map[string]map[string]string, len(g.rawMap))
	for start := range g.rawMap {
		prev := make(map[string]string, len(g.rawMap))
		dist := make(map[string]uint32, len(g.rawMap))
		for name := range g.rawMap {
			if name == start {
				dist[name] = 0
			} else {
				dist[name] = maxPriority
			}
		}

		g.generateQueue(start)

		// Алгоритм Дейкстры
		for g.pq.Len() > 0 {
			city := heap.Pop(&g.pq).(*Item)
			for name, distance := range city.adjacents {

				if dist[name] > dist[city.name]+uint32(distance) {
					dist[name] = dist[city.name] + uint32(distance)
					prev[name] = city.name

					for i := range g.pq {
						if g.pq[i].name == name {
							item := g.pq[i]
							g.pq.update(item, dist[name])
						}
					}

				}
			}
		}
		g.fullMap[start] = prev
	}
}

// extractPath возвращает кратчайшее расстяние между двумя вершинами
// из сгенерированной ранее карты путей
func (g *graph) extractPath(start, end string) ([]string, error) {
	startCheck, endCheck := g.checkVertices(start, end)
	if !startCheck {
		return nil, errors.New("Вершина не существует: " + start)
	}
	if !endCheck {
		return nil, errors.New("Вершина не существует: " + end)
	}
	pathMap := g.fullMap[end]
	path := []string{start}
	for currentVertice := pathMap[start]; ; currentVertice = pathMap[currentVertice] {
		path = append(path, currentVertice)
		if currentVertice == end {
			break
		}
	}

	return path, nil
}

// checkVertices проверяет наличие вершин в графе
func (g *graph) checkVertices(start, end string) (bool, bool) {
	_, startCheck := g.fullMap[start]
	_, endCheck := g.fullMap[end]
	return startCheck, endCheck
}
