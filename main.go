package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Depado/cuitesite/cmd"
	"github.com/Depado/cuitesite/fetch"
	"github.com/Depado/cuitesite/infra"
	"github.com/Depado/cuitesite/router"
)

// Build number and versions injected at compile time
var (
	Version = "unknown"
	Build   = "unknown"
)

// Main command that will be run when no other command is provided on the
// command-line
var rootc = &cobra.Command{
	Use:   "cuitesite <options>",
	Short: "Cuitesite backend",
	Long:  "Backend app that will aggregate playlists",
	Run:   func(cmd *cobra.Command, args []string) { run() },
}

// Version command that will display the build number and version (if any)
var versionc = &cobra.Command{
	Use:   "version",
	Short: "Show build and version",
	Run:   func(c *cobra.Command, args []string) { fmt.Printf("Build: %s\nVersion: %s\n", Build, Version) },
}

func run() {
	strids := viper.GetStringSlice("user_ids")
	ids := make([]int, len(strids))
	for i := 0; i < len(strids); i++ {
		var err error
		if ids[i], err = strconv.Atoi(strids[i]); err != nil {
			logrus.WithError(err).Fatal("not an integer ID")
		}
	}
	s := infra.NewServer(
		viper.GetString("server.host"),
		viper.GetInt("server.port"),
		viper.GetString("server.mode"),
		infra.NewCorsConfig(
			viper.GetBool("server.cors.disabled"),
			viper.GetBool("server.cors.all"),
			viper.GetStringSlice("server.cors.origins"),
			viper.GetStringSlice("server.cors.methods"),
			viper.GetStringSlice("server.cors.headers"),
			viper.GetStringSlice("server.cors.expose"),
		),
	)
	fc := fetch.NewClient(
		viper.GetString("client_id"),
		ids,
	)
	if err := fc.Fetch(); err != nil {
		logrus.WithError(err).Fatal("Unable to fetch content from Soundcloud")
	}
	gr := &router.GinRouter{Playlists: fc.Playlists, Tracks: fc.Tracks, PlaylistsMap: &fc.PlaylistsMap}
	gr.AddRoutes(s.Router)
	s.Start()
}

func main() {
	// Initialize Cobra and Viper
	cobra.OnInitialize(cmd.Initialize)
	cmd.AddAllFlags(rootc)

	// Run the command
	if err := rootc.Execute(); err != nil {
		log.Fatal("Couldn't start program:", err)
	}
}
