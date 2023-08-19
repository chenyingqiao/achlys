/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("init start %v", args)
		if len(args) < 1 {
			log.Fatalf("init failed!")
		}
		if err := InitContainerProcess(args[0],args[1:]); err != nil {
			log.Fatal(err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


func InitContainerProcess(command string, args []string) error {
	log.Printf("start init container first process: %s\n", command)
	syscall.Mount("proc","/proc","proc",uintptr(syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV), "")
	if err := syscall.Exec(command,args,os.Environ()); err != nil {
		log.Printf("exec command error:%s\n",err.Error())
		return  err
	}
	return nil
}
