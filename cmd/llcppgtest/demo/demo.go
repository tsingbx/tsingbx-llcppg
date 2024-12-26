package demo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func runSingleDemo(demoPath string) {
	fmt.Printf("Testing demo: %s\n", demoPath)

	absPath, err := filepath.Abs(demoPath)
	if err != nil {
		panic(fmt.Sprintf("failed to get absolute path for %s: %v", demoPath, err))
	}

	configFile := filepath.Join(absPath, "llcppg.cfg")
	fmt.Printf("Looking for config file at: %s\n", configFile)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file not found: %s", configFile))
	}

	llcppgArgs := []string{"-v"}

	outDir := filepath.Join(absPath, "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output directory: %v", err))
	}

	// copy llcppg.cfg to out dir
	cfgFile := filepath.Join(outDir, "llcppg.cfg")
	cfgFileContent, err := os.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("failed to read config file: %v", err))
	}
	if err := os.WriteFile(cfgFile, cfgFileContent, 0644); err != nil {
		panic(fmt.Sprintf("failed to write config file: %v", err))
	}

	// run llcppg to gen pkg
	if err := runCommand(outDir, "llcppg", llcppgArgs...); err != nil {
		panic(fmt.Sprintf("llcppg execution failed: %v", err))
	}
	fmt.Printf("llcppg execution success\n")

	// check if the gen pkg is ok
	resultDir := filepath.Join(outDir, filepath.Base(absPath))
	if err := runCommand(resultDir, "llgo", "build", "."); err != nil {
		panic(fmt.Sprintf("llgo build failed in %s: %v", resultDir, err))
	}
	fmt.Printf("llgo build success\n")

	// go mod tidy,because the demo is dependent on the gen pkg
	modPath := filepath.Join(demoPath, "demo")
	if err := runCommand(modPath, "go", "mod", "tidy"); err != nil {
		panic(fmt.Sprintf("go mod tidy failed in %s: %v", demoPath, err))
	}
	fmt.Printf("go mod tidy success\n")

	// run the demo
	demos, err := os.ReadDir(modPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read demo directory: %v", err))
	}
	for _, demo := range demos {
		if demo.IsDir() {
			runCommand(filepath.Join(modPath, demo.Name()), "llgo", "run", ".")
		}
	}
}

// Get all first-level directories containing llcppg.cfg
func getFirstLevelDemos(baseDir string) []string {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		panic(fmt.Sprintf("failed to read directory: %v", err))
	}

	var demos []string
	for _, entry := range entries {
		if entry.IsDir() {
			demoPath := filepath.Join(baseDir, entry.Name())
			configPath := filepath.Join(demoPath, "llcppg.cfg")
			if _, err := os.Stat(configPath); err == nil {
				demos = append(demos, demoPath)
			}
		}
	}
	return demos
}

func TestDemos(path string) {
	fmt.Printf("Starting tests in directory: %s\n", path)

	stat, err := os.Stat(path)
	if err != nil || !stat.IsDir() {
		panic(fmt.Sprintf("specified path is not a directory or does not exist: %s", path))
	}

	demos := getFirstLevelDemos(path)
	if len(demos) == 0 {
		panic(fmt.Sprintf("no directories containing llcppg.cfg found in %s", path))
	}

	// Test each demo
	for _, demo := range demos {
		runSingleDemo(demo)
		fmt.Printf("Test completed for %s\n", demo)
	}
}

func runCommand(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
