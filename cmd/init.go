/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"mydocker/pkg"
	"os"
	"os/exec"
	"path/filepath"
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
		log.Printf("开始init初始化: %v",args)
		argsCommand,err := pkg.ReadExtraFd()
		if err != nil {
			log.Printf("读取管道错误:%s", err.Error())
		}
		log.Printf("管道参数：%s",argsCommand)
		if len(argsCommand) == 0 {
			log.Fatalf("init参数错误\n")
		}
		if err := InitContainerProcess(argsCommand[0],argsCommand[0:]); err != nil {
			log.Fatal(err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func InitContainerProcess(command string, args []string) error {
	log.Printf("start init container first process: %s\n", command)
	//获取当前进程运行的目录作为根目录
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Printf("当前目录：%s",pwd)
	err  = pkg.PivotRoot(pwd)
	if err != nil {
		return err
	}
	pkg.MountSystemKeyDir()
	pwd, err = os.Getwd()
	if err != nil {
		return err
	}
	log.Printf("当前目录：%s",pwd)
	files,_ := filepath.Glob("*")
	log.Printf("文件列表:%v",files)
	log.Printf("执行commamd： %v %v",command, args)
	path,err := exec.LookPath(command)
	if err != nil {
		log.Printf("找不到对应的可执行程序：%s",path)
		return err
	}
	log.Printf("执行命令位置：%s",path)
	if err := syscall.Exec(path, args, os.Environ()); err != nil {
		log.Printf("exec command error:%s\n", err.Error())
		return err
	}
	return nil
}
