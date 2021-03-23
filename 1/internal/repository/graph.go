package repository

import (
	"gocompiler/internal/fsm"
	"gocompiler/internal/graph"
	"gocompiler/internal/visualizer"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/awalterschulze/gographviz"
)

type FileLoader struct {
	Path string
}

func NewFileLoader(path string) *FileLoader {
	return &FileLoader{Path: path}
}

// Load - загрузить граф из файла
func (fl FileLoader) Load() (*graph.Graph, error) {
	file, err := os.Open(fl.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	graphAst, err := gographviz.Parse(b)
	if err != nil {
		return nil, err
	}

	gr := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, gr); err != nil {
		return nil, err
	}

	var graphS = graph.NewGraph()
	for _, e := range gr.Edges.Edges {
		weight := strings.Replace(e.Attrs["label"], `<<font color="blue">`, "", -1)
		weight = strings.Replace(weight, `</font>>`, "", -1)
		graphS.AddEdge(&graph.Edge{
			From:   e.Src,
			To:     e.Dst,
			Weight: weight,
		})
	}
	var (
		starts, ends []string
	)
	for _, n := range gr.Nodes.Nodes {
		weight := strings.Replace(n.Attrs["label"], `<<font color="blue">`, "", -1)
		if strings.Contains(weight, "green") {
			starts = append(starts, n.Name)
		} else if strings.Contains(weight, "red") {
			ends = append(ends, n.Name)
		}
	}
	graphS.SetFirstLast(starts, ends)

	return graphS, nil
}

type Cache struct {
	c map[int]*fsm.FSM
}

func NewCache() *Cache {
	return &Cache{c: make(map[int]*fsm.FSM, 0)}
}

func (c *Cache) Put(code int, g *fsm.FSM) {
	c.c[code] = g
}

func (c *Cache) Take(code int) *fsm.FSM {
	return c.c[code]
}

// LoadGraf возвращает флаг надо ли пойти назад
func LoadGraf(
	prevCode, newCode int,
	cache Cache,
	action func(graf *fsm.FSM) error,
) {
	var graf = cache.Take(prevCode)
	log.Println("Выберите способ, чтобы задать граф")
	log.Println("1. Задать новый гриф с помощью файла")
	if graf != nil {
		log.Println("2. Использовать данные, полученные в прошлом действии")
	}
	log.Println("X. Назад")
	var repCode int
	GetInt(&repCode, `Ваш выбор?`)
	if repCode == 2 {
		if graf == nil {
			return
		}
	} else if repCode == 1 {
		var (
			text string
		)
		GetString(&text, `Введите путь до файла`)
		l := NewFileLoader(text)
		grafF, err := l.Load()
		if err != nil {
			log.Printf("Не удалось загрузить граф из файла: %s", err)
			return
		}
		graf = &fsm.FSM{grafF}
		err = visualizer.VisualizeFSM(graf, "assets", "loaded.dot")
		if err != nil {
			log.Printf("Не удалось визуализировтаь граф: %s", err)
			return
		}
	} else {
		return
	}

	err := action(graf)
	if err != nil {
		log.Printf("Не удалось выполнить действие: %s", err)
		return
	}
	err = visualizer.VisualizeFSM(graf, "assets", "main.dot")
	if err != nil {
		log.Printf("Не удалось визуализировать граф: %s", err)
		return
	}
	cache.Put(newCode, graf)
	return
}
