package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/raz0red/radiscordpresence/internal/buildinfo"
	"github.com/raz0red/radiscordpresence/internal/config"
	"github.com/raz0red/radiscordpresence/internal/presence"
	"github.com/raz0red/radiscordpresence/internal/svc"
)

func main() {
	// If --config-dir is present in os.Args, apply it before anything else and
	// strip it so cobra never sees an unrecognised flag. This is injected at
	// 'install' time so that when Windows SCM starts the service as LocalSystem
	// (which has no user APPDATA), it still finds the correct config file.
	for i, arg := range os.Args {
		if arg == "--config-dir" && i+1 < len(os.Args) {
			config.OverrideDir = os.Args[i+1]
			os.Args = append(os.Args[:i], os.Args[i+2:]...)
			break
		}
	}

	// When invoked by a service manager (Windows SCM, macOS launchd, Linux systemd),
	// hand off to the service runner immediately — no CLI parsing needed.
	if !svc.Interactive() {
		if err := svc.Run(); err != nil {
			log.Fatalf("service error: %v", err)
		}
		return
	}

	root := &cobra.Command{
		Use:     "radiscordpresence",
		Short:   "RetroAchievements Discord Rich Presence",
		Version: buildinfo.String(),
		Long: `RADiscordPresence mirrors your RetroAchievements session to Discord Rich Presence.

Quick start:
  radiscordpresence set --username YOUR_NAME --apikey YOUR_KEY
  radiscordpresence run`,
		SilenceUsage: true,
	}

	root.AddCommand(
		cmdSet(),
		cmdRun(),
		cmdInstall(),
		cmdUninstall(),
		cmdStart(),
		cmdStop(),
		cmdStatus(),
		cmdVersion(),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// cmdSet — save credentials and options to the config file.
func cmdSet() *cobra.Command {
	var username, apikey string
	var interval int

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Save credentials and settings to config",
		Example: `  radiscordpresence set --username YOUR_NAME --apikey YOUR_KEY
  radiscordpresence set --interval 30`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			changed := false
			if cmd.Flags().Changed("username") {
				cfg.Username = username
				changed = true
			}
			if cmd.Flags().Changed("apikey") {
				cfg.APIKey = apikey
				changed = true
			}
			if cmd.Flags().Changed("interval") {
				cfg.Interval = interval
				changed = true
			}

			if !changed {
				// No flags given — show current config and usage hint.
				dir, _ := config.Dir()
				fmt.Printf("Config file: %s\\config.json\n\n", dir)
				if cfg.Username != "" {
					fmt.Printf("  username: %s\n", cfg.Username)
				} else {
					fmt.Printf("  username: (not set)\n")
				}
				if cfg.APIKey != "" {
					fmt.Printf("  api_key:  (set)\n")
				} else {
					fmt.Printf("  api_key:  (not set)\n")
				}
				fmt.Printf("  interval: %d seconds\n", cfg.Interval)
				fmt.Println()
				fmt.Println("To update, use:")
				fmt.Println("  radiscordpresence set --username YOUR_RA_USERNAME --apikey YOUR_WEB_API_KEY")
				fmt.Println()
				fmt.Println("Your Web API key is at: https://retroachievements.org/controlpanel.php")
				return nil
			}

			if err := config.Save(cfg); err != nil {
				return err
			}
			dir, _ := config.Dir()
			fmt.Printf("Config saved to %s\n", dir)
			return nil
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "RetroAchievements username")
	cmd.Flags().StringVar(&apikey, "apikey", "", "RetroAchievements Web API key")
	cmd.Flags().IntVar(&interval, "interval", 10, "Poll interval in seconds")
	return cmd
}

// cmdRun — run in the foreground (Ctrl+C to stop).
func cmdRun() *cobra.Command {
	var username, apikey string
	var interval int

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run in the foreground (Ctrl+C to stop)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			// Inline flags override saved config.
			if cmd.Flags().Changed("username") {
				cfg.Username = username
			}
			if cmd.Flags().Changed("apikey") {
				cfg.APIKey = apikey
			}
			if cmd.Flags().Changed("interval") {
				cfg.Interval = interval
			}
			if cfg.Username == "" || cfg.APIKey == "" {
				return fmt.Errorf(
					"username and apikey are required\n" +
						"  save permanently: radiscordpresence set --username X --apikey Y\n" +
						"  or pass inline:   radiscordpresence run --username X --apikey Y",
				)
			}

			w := presence.New(cfg.Username, cfg.APIKey, cfg.Interval)

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt)
			go func() {
				<-sig
				fmt.Println()
				w.Stop()
			}()

			w.Run()
			return nil
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "Override username from config")
	cmd.Flags().StringVar(&apikey, "apikey", "", "Override API key from config")
	cmd.Flags().IntVar(&interval, "interval", 0, "Override poll interval in seconds")
	return cmd
}

// cmdInstall — register as a system service.
func cmdInstall() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install as a system service (requires elevated privileges)",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := svc.Install(); err != nil {
				return fmt.Errorf("install failed: %w\nTip: on Windows run as Administrator; on Linux/macOS use sudo", err)
			}
			fmt.Println("Service installed. Run: radiscordpresence start")
			return nil
		},
	}
}

// cmdUninstall — remove the system service.
func cmdUninstall() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the system service",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := svc.Uninstall(); err != nil {
				return fmt.Errorf("uninstall failed: %w", err)
			}
			fmt.Println("Service uninstalled.")
			return nil
		},
	}
}

// cmdStart — start the installed service.
func cmdStart() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the installed service",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := svc.Start(); err != nil {
				return fmt.Errorf("start failed: %w", err)
			}
			fmt.Println("Service started.")
			return nil
		},
	}
}

// cmdStop — stop the running service.
func cmdStop() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the running service",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := svc.Stop(); err != nil {
				return fmt.Errorf("stop failed: %w", err)
			}
			fmt.Println("Service stopped.")
			return nil
		},
	}
}

// cmdStatus — show current service status.
func cmdStatus() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show service status",
		RunE: func(_ *cobra.Command, _ []string) error {
			st, err := svc.Status()
			if err != nil {
				return fmt.Errorf("status failed: %w", err)
			}
			fmt.Printf("Service: %s\n", st)
			return nil
		},
	}
}

// cmdVersion — print build version information.
func cmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("RADiscordPresence %s\n", buildinfo.String())
		},
	}
}
