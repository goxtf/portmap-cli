package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var watchInterval int

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for port mapping changes in real time",
	Long:  "Continuously monitor port bindings and print added/removed entries as they change.",
	RunE:  runWatch,
}

func init() {
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 2, "Poll interval in seconds")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("failed to create scanner: %w", err)
	}

	formatter := ports.NewFormatter(false)

	onChange := func(added, removed []ports.PortEntry) {
		for _, e := range added {
			fmt.Fprintf(os.Stdout, "[+] %s\n", formatter.FormatEntry(e))
		}
		for _, e := range removed {
			fmt.Fprintf(os.Stdout, "[-] %s\n", formatter.FormatEntry(e))
		}
	}

	watcher, err := ports.NewWatcher(ports.WatcherConfig{
		Scanner:  scanner,
		Interval: time.Duration(watchInterval) * time.Second,
		OnChange: onChange,
	})
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Fprintf(os.Stdout, "Watching for port changes every %ds (Ctrl+C to stop)...\n", watchInterval)

	if err := watcher.Start(ctx); err != nil && err != context.Canceled {
		return err
	}
	return nil
}
