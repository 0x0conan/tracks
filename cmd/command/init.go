package command

import (
	"github.com/airchains-network/decentralized-sequencer/config"
	"github.com/airchains-network/decentralized-sequencer/p2p"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var moniker string
var stationType string
var daType string
var daRPC string
var stationRPC string
var junctionRPC string

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the sequencer nodes",
	Run: func(cmd *cobra.Command, args []string) {
		moniker, _ = cmd.Flags().GetString("moniker")
		stationType, _ = cmd.Flags().GetString("stationType")
		daType, _ = cmd.Flags().GetString("daType")
		daRPC, _ = cmd.Flags().GetString("daRpc")
		stationRPC, _ = cmd.Flags().GetString("stationRpc")

		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err) // Handle error appropriately
		}
		tracksDir := filepath.Join(homeDir, config.DefaultTracksDir)

		conf := config.DefaultConfig()

		conf.RootDir = tracksDir
		conf.DA.DaType = daType // set daType
		conf.DA.DaRPC = daRPC   // set daRPC
		conf.Station.StationType = stationType
		conf.Station.StationRPC = stationRPC
		conf.SetRoot(conf.RootDir)
		config.EnsureRoot(conf.RootDir, conf)

		p2p.InititateIdentity(daType, moniker, stationType)
	},
}
