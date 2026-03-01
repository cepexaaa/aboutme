package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/application/service/generator"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/application/service/solver"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/infrastructure"
)

//go:embed about.txt
var about string

func main() {
	if len(os.Args) < 2 {
		info()
		return
	}
	config, err := parseSetUps(os.Args[2:])
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Uncorrect input parametrs. Example to run")
		info()
		return
	}

	writer, closeFn, err := getOutputWriter(config.Output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer closeFn()

	var labyrinth *domain.Labyrinth
	isSolve := false
	if os.Args[1] == "generate" {
		labyrinth = generateMaze(config)
	} else {
		labyrinth, err = solveMaze(config)
		isSolve = true
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	infrastructure.ShowResult(labyrinth, writer, config.Unicode, isSolve)
}

// read maze file, use input parameters to solve maze
func solveMaze(config *config) (*domain.Labyrinth, error) {
	field := infrastructure.ReadMaze(config.File)
	if field == nil {
		return nil, errors.New("can not read labyrint in the file " + config.File)
	}
	startPoint, err1 := parseCoords(config.Start)
	endPoint, err2 := parseCoords(config.End)
	if err1 != nil || err2 != nil {
		return nil, errors.Join(err1, err2)
	}
	solve := solver.NewSolver(startPoint, endPoint)
	solve.NewMaze(config.Algorithm, uint(len(field)), uint(len(field[0])), field)
	err := solve.IsLogicValid()
	if err != nil {
		return solve.Maze, err
	}
	solve.Solve()
	return solve.Maze, nil
}

// generate maze by configuration
func generateMaze(config *config) *domain.Labyrinth {
	gen := generator.NewGenerator(config.Algorithm, uint(config.Width), uint(config.Height))
	gen.Generate()
	return gen.Maze
}

// return: writer and close function
func getOutputWriter(output string) (io.Writer, func() error, error) {
	if output == "" {
		return os.Stdout, func() error { return nil }, nil
	}

	file, err := os.Create(output)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create file %s: %v", output, err)
	}

	return file, file.Close, nil
}
