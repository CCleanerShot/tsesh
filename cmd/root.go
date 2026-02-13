package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/e-mar404/tsesh/config"
	"github.com/e-mar404/tsesh/picker"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tsesh",
		Short: "terminal sessionizer extending tmux",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				// TODO: this might want to be done through logging and only showup if there is a verbose flag
				fmt.Printf("could not load configuration: %v\n", err)
				fmt.Printf("using default paths: [\"~/\", \"~/projects\"]\n")
				cfg = &config.Config{
					SearchPaths: []string{
						"~/",
						"~/projects",
					},
				}
			}

			p := tea.NewProgram(picker.New(cfg.SearchPaths), tea.WithAltScreen())
			if pi, err := p.Run(); err != nil {
				fmt.Printf("%v\n", pi.(picker.Picker).Err)
				fmt.Printf("Encountered an error when trying to run the directory picker: %v\n", err)
				os.Exit(1)
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("could not execute command: %v\n", err)
		os.Exit(1)
	}
}
