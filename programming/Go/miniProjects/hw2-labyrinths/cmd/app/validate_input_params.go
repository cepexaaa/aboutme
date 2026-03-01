package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw2-labyrinths/internal/domain"
)

var (
	patternCord = `^\s*(0|[1-9]\d*)\s*,\s*(0|[1-9]\d*)\s*$`
	coordinates = regexp.MustCompile(patternCord)
)

type config struct {
	Algorithm string
	Width     int
	Height    int
	Output    string
	File      string
	Start     string
	End       string
	Unicode   bool
}

// parse input parameters
func parseSetUps(args []string) (*config, error) {
	if len(os.Args) < 2 {
		info()
	}
	var config *config
	var err error
	switch strings.ToLower(os.Args[1]) {
	case "generate":
		config, err = validateInput(generateConfig, validateGenerateParams, args)
	case "solve":
		config, err = validateInput(solveConfig, validateSolveParams, args)
	default:
		err = errors.New("Справка по использованию программы")
		fmt.Println(about)
	}

	if err != nil {
		return nil, err
	}
	return config, nil
}

// function for both generate and solve
func validateInput(cfg func(args []string) (*config, error), valid func(cfg *config) error, args []string) (*config, error) {
	c, err := cfg(args)
	if err != nil {
		fmt.Println("not valid input parametrs")
		return nil, err
	}
	err = valid(c)
	if err != nil {
		fmt.Println("you are using unknown values for this labyrinth")
		return nil, err
	}
	return c, nil
}

func validateSolveParams(config *config) error {
	if config.Algorithm != "astar" && config.Algorithm != "dijkstra" && config.Algorithm != "digger" && config.Algorithm != "bfs" && config.Algorithm != "optimal" {
		return errors.New("unsupported algorithm for solving labyrinth")
	}
	if !coordinates.MatchString(config.End) {
		return errors.New("uncorrect finish coordinates, expected format: x,y")
	}
	if !coordinates.MatchString(config.Start) {
		return errors.New("uncorrect start coordinates, expected format: x,y")
	}
	if len(config.File) == 0 {
		return errors.New("file with maze is unfound")
	}
	return nil
}

func validateGenerateParams(config *config) error {
	if config.Algorithm != "dfs" && config.Algorithm != "prim" && config.Algorithm != "kruskal" {
		return errors.New("unsupported algorithm for generating labyrinth")
	}
	if config.Width <= 2 || config.Height <= 2 {
		return errors.New("so few maze. Input size bigger than it")
	}
	return nil
}

func info() {
	fmt.Println("Usage:")
	fmt.Println("  generate --algorithm=dfs --width=10 --height=10 --output=maze.txt")
	fmt.Println("  solve --algorithm=dijkstra --file=maze.txt --start=0,0 --end=9,9 --output=solution.txt")
}

func generateConfig(args []string) (*config, error) {
	var config config
	flags := flag.NewFlagSet("generate", flag.ContinueOnError)
	flags.StringVar(&config.Algorithm, "algorithm", "dfs", "Algorithm for labyrinth generation")
	flags.IntVar(&config.Width, "width", 10, "Width of the maze")
	flags.IntVar(&config.Height, "height", 10, "Height of the maze")
	flags.StringVar(&config.Output, "output", "", "Output file name (stdout if empty)")
	flags.BoolVar(&config.Unicode, "unicode", false, "Use unicode characters for maze display")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	return &config, nil
}

func solveConfig(args []string) (*config, error) {
	var config config
	flags := flag.NewFlagSet("solve", flag.ContinueOnError)
	flags.StringVar(&config.Algorithm, "algorithm", "dijkstra", "Algorithm for solving maze (dijkstra, astar, etc.)")
	flags.StringVar(&config.File, "file", "", "Input maze file")
	flags.StringVar(&config.Start, "start", "0,0", "Start coordinates (x,y)")
	flags.StringVar(&config.End, "end", "0,0", "End coordinates (x,y)")
	flags.StringVar(&config.Output, "output", "", "Output file name (stdout if empty)")
	flags.BoolVar(&config.Unicode, "unicode", false, "Use unicode characters for maze display")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	return &config, nil
}

func parseCoords(str string) (domain.Pair[uint], error) {
	s1 := strings.Split(str, ",")
	if len(s1) != 2 {
		return domain.Pair[uint]{}, errors.New("uncorrect coordinates")
	}
	x, err1 := strconv.Atoi(s1[0])
	y, err2 := strconv.Atoi(s1[1])
	if err1 != nil || err2 != nil {
		return domain.Pair[uint]{}, errors.New("uncorrect coordinates")
	}
	return domain.Pair[uint]{First: uint(x), Second: uint(y)}, nil
}
