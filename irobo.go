package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type Stage struct {
	Commands []string `yaml:"commands"`
}
type Stages struct {
	Stages map[string]Stage `yaml:"stages"`
}

func runCommand(cmd string) error {
	fmt.Println(color.GreenString("Running command:"), cmd)
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	s.Stop()
	if err != nil {
		fmt.Println(color.RedString(string(out)))
		return fmt.Errorf("error running command: %w", err)
	}
	fmt.Println(string(out))
	return nil
}
func runStage(stage string, commands []string) error {
	fmt.Println(color.CyanString("\nRunning stage:"), stage)
	for _, cmd := range commands {
		if err := runCommand(cmd); err != nil {
			return fmt.Errorf("error running stage %q: %w", stage, err)
		}
	}
	return nil
}
func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid env format: %q", line)
		}
		os.Setenv(parts[0], parts[1])
	}
	return scanner.Err()
}
func main() {
	var stageFile, envFile, stageName string
	flag.StringVar(&stageFile, "tasks", "stages.yml", "Path to the stages file")
	flag.StringVar(&envFile, "env", ".env", "Path to the environment file")
	flag.StringVar(&stageName, "stage", "", "Name of the stage to run")
	flag.Parse()

	if err := loadEnv(envFile); err != nil {
		fmt.Println(color.RedString("Error loading .env file:"), err)
		return
	}
	yamlFile, err := ioutil.ReadFile(stageFile)
	if err != nil {
		fmt.Println(color.RedString("Error loading stages.yml file:"), err)
		return
	}
	var stages Stages
	if err := yaml.Unmarshal(yamlFile, &stages); err != nil {
		fmt.Println(color.RedString("Error parsing stages.yml file:"), err)
		return
	}
	if stageName != "" {
		stage, ok := stages.Stages[stageName]
		if !ok {
			fmt.Println(color.RedString("Error: stage", stageName, "not found"))
			return
		}
		if err := runStage(stageName, stage.Commands); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		for stage, s := range stages.Stages {
			if err := runStage(stage, s.Commands); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
