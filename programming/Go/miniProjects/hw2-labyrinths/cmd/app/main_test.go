package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestGenerateMaze(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "Valid DFS generation",
			args:      []string{"generate", "--algorithm=dfs", "--width=10", "--height=10"},
			wantError: false,
		},
		{
			name:      "Valid Prim generation",
			args:      []string{"generate", "--algorithm=prim", "--width=15", "--height=15"},
			wantError: false,
		},
		{
			name:      "Generation with unicode",
			args:      []string{"generate", "--algorithm=dfs", "--width=8", "--height=8", "--unicode"},
			wantError: false,
		},
		{
			name:      "Invalid algorithm",
			args:      []string{"generate", "--algorithm=invalid", "--width=10", "--height=10"},
			wantError: true,
		},
		{
			name:      "Missing width",
			args:      []string{"generate", "--algorithm=dfs", "--height=10"},
			wantError: false,
		},
		{
			name:      "Missing height",
			args:      []string{"generate", "--algorithm=dfs", "--width=10"},
			wantError: false,
		},
		{
			name:      "Zero dimensions",
			args:      []string{"generate", "--algorithm=dfs", "--width=0", "--height=0"},
			wantError: true,
		},
		{
			name:      "Negative dimensions",
			args:      []string{"generate", "--algorithm=dfs", "--width=-5", "--height=-5"},
			wantError: true,
		},
		{
			name:      "Too small dimensions",
			args:      []string{"generate", "--algorithm=dfs", "--width=1", "--height=1"},
			wantError: true,
		},
		{
			name:      "Even dimensions correction",
			args:      []string{"generate", "--algorithm=dfs", "--width=12", "--height=12"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			oldStdout := os.Stdout
			defer func() {
				os.Args = oldArgs
				os.Stdout = oldStdout
			}()

			os.Args = append([]string{"cmd"}, tt.args...)

			r, w, _ := os.Pipe()
			os.Stdout = w

			main()

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			if tt.wantError {
				if !strings.Contains(output, "Uncorrect input parametrs") &&
					!strings.Contains(output, "Error") {
					t.Errorf("%s: ожидалась ошибка, но получен вывод: %s", tt.name, output)
				}
			} else {
				if strings.Contains(output, "Error") || strings.Contains(output, "Uncorrect") {
					t.Errorf("%s: не ожидалась ошибка, но получено: %s", tt.name, output)
				}
				if !strings.Contains(output, "#") && !strings.Contains(output, "┌") && !strings.Contains(output, " ") {
					t.Errorf("%s: вывод не содержит ожидаемого лабиринта: %s", tt.name, output)
				}
			}
		})
	}
}

func TestSolveMaze(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "Valid A* solution",
			args:      []string{"solve", "--algorithm=astar", "--file=test_maze.txt", "--start=0,0", "--end=4,4"},
			wantError: false,
		},
		{
			name:      "Valid Dijkstra solution",
			args:      []string{"solve", "--algorithm=dijkstra", "--file=test_maze.txt", "--start=1,1", "--end=3,3"},
			wantError: false,
		},
		{
			name:      "Valid BFS solution",
			args:      []string{"solve", "--algorithm=bfs", "--file=test_maze.txt", "--start=0,1", "--end=4,3"},
			wantError: false,
		},
		{
			name:      "Solution with output file",
			args:      []string{"solve", "--algorithm=astar", "--file=test_maze.txt", "--start=0,0", "--end=4,4", "--output=solution.txt"},
			wantError: false,
		},
		{
			name:      "Solution with unicode",
			args:      []string{"solve", "--algorithm=astar", "--file=test_maze.txt", "--start=0,0", "--end=4,4", "--unicode"},
			wantError: false,
		},
		{
			name:      "Non-existent file",
			args:      []string{"solve", "--algorithm=astar", "--file=nonexistent.txt", "--start=0,0", "--end=4,4"},
			wantError: true,
		},
		{
			name:      "Invalid coordinates format",
			args:      []string{"solve", "--algorithm=astar", "--file=test_maze.txt", "--start=0-0", "--end=4,4"},
			wantError: true,
		},
		{
			name:      "Coordinates out of bounds",
			args:      []string{"solve", "--algorithm=astar", "--file=test_maze.txt", "--start=10,10", "--end=20,20"},
			wantError: true,
		},
		{
			name:      "Start in wall",
			args:      []string{"solve", "--algorithm=astar", "--file=test_maze.txt", "--start=0,0", "--end=4,4"},
			wantError: true,
		},
		{
			name:      "End in wall",
			args:      []string{"solve", "--algorithm=astar", "--file=test_maze.txt", "--start=1,1", "--end=0,0"},
			wantError: true,
		},
		{
			name:      "Missing start",
			args:      []string{"solve", "--algorithm=astar", "--file=test_maze.txt", "--end=4,4"},
			wantError: true,
		},
		{
			name:      "Missing end",
			args:      []string{"solve", "--algorithm=astar", "--file=test_maze.txt", "--start=0,0"},
			wantError: true,
		},
		{
			name:      "Missing file",
			args:      []string{"solve", "--algorithm=astar", "--start=0,0", "--end=4,4"},
			wantError: true,
		},
		{
			name:      "Invalid solving algorithm",
			args:      []string{"solve", "--algorithm=invalid", "--file=test_maze.txt", "--start=0,0", "--end=4,4"},
			wantError: true,
		},
	}

	createTestMazeFile()
	defer os.Remove("test_maze.txt")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = append([]string{"cmd"}, tt.args...)

			r, w, _ := os.Pipe()
			oldStdout := os.Stdout
			os.Stdout = w

			main()

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			if tt.wantError {
				if !strings.Contains(output, "Error") &&
					!strings.Contains(output, "Uncorrect") &&
					!strings.Contains(output, "out of") &&
					!strings.Contains(output, "on wall") &&
					!strings.Contains(output, "can not read") {
					t.Errorf("%s: ожидалась ошибка, но получен вывод: %s", tt.name, output)
				}
			} else {
				if strings.Contains(output, "Error") || strings.Contains(output, "Uncorrect") {
					t.Errorf("%s: не ожидалась ошибка, но получено: %s", tt.name, output)
				}
			}
		})
	}
}

func createTestMazeFile() string {
	content := `#####
#   #
# # #
#   #
#####`

	tmpfile, err := os.Create("test_maze.txt")
	if err != nil {
		panic(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		panic(err)
	}

	if err := tmpfile.Close(); err != nil {
		panic(err)
	}

	return tmpfile.Name()
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "No arguments",
			args:      []string{},
			wantError: true,
		},
		{
			name:      "Unknown command",
			args:      []string{"unknown"},
			wantError: true,
		},
		{
			name:      "Help command",
			args:      []string{"help"},
			wantError: true,
		},
		{
			name:      "Very large maze",
			args:      []string{"generate", "--algorithm=dfs", "--width=100", "--height=100"},
			wantError: false,
		},
		{
			name:      "Minimum valid maze",
			args:      []string{"generate", "--algorithm=dfs", "--width=3", "--height=3"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = append([]string{"cmd"}, tt.args...)

			r, w, _ := os.Pipe()
			oldStdout := os.Stdout
			os.Stdout = w

			main()

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			if tt.wantError {
				if !strings.Contains(output, "Error") &&
					!strings.Contains(output, "Uncorrect") &&
					!strings.Contains(output, "Usage:") &&
					!strings.Contains(output, "can not read") {
					t.Errorf("%s: ожидалась ошибка, но получен вывод: %s", tt.name, output)
				}
			} else {
				if strings.Contains(output, "Error") || strings.Contains(output, "Uncorrect") {
					t.Errorf("%s: не ожидалась ошибка, но получено: %s", tt.name, output)
				}
			}
		})
	}
}
