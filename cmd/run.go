/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"mydocker/pkg"

	"github.com/spf13/cobra"
)

const selfExe = "/proc/self/exe"

var RunArgsD RunArgs

type RunArgs struct{
	KeepStedin *bool
	TTY *bool
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "在容器中执行一个命令",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		initArgs := []string{"init"}
		initArgs = append(initArgs, args...)
		err := pkg.NewProcess().Run(selfExe, pkg.PorcessOption{
			KeepStedin: *RunArgsD.KeepStedin,
			TTY: *RunArgsD.TTY,
			Args: initArgs,
		})
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	RunArgsD.KeepStedin = runCmd.Flags().BoolP("interactive","i",false,"Keep STDIN open even if not attached")
	RunArgsD.TTY = runCmd.Flags().BoolP("tty","t",false,"Allocate a pseudo-TTY")
}
