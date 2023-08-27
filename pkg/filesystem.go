package pkg

import (
	"errors"
	"fmt"
	"os"
	"log"
	"path/filepath"
	"syscall"
)

type FileSystem struct{}

func PivotRoot(root string) error {
	err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		return errors.New("挂载root文件出错")
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("创建临时root目录出错:%s", err.Error())
	}

	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("切换根文件系统失败:%s", err.Error())
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("切换工作目录失败:%s", err.Error())
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	//syscall.MNT_DETACH 是一个常量，通常在卸载文件系统时用于指定卸载操作的标志。它的作用是告诉操作系统在卸载文件系统时，如果有其他进程仍然在使用这个文件系统，不要立即卸载它，而是等待其他进程使用完后再卸载。
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf(":%s", err.Error())
	}
	return os.Remove(pivotDir)
}

func MountSystemKeyDir(){
	defaultMountFlags := syscall.MS_NOEXEC |syscall.MS_NOSUID | syscall.MS_STRICTATIME
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID | syscall.MS_STRICTATIME, "mode=755")
}

func Ls(dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历文件列表并打印文件名
	for _, file := range files {
		// 如果要获取文件的完整路径，可以使用 filepath.Join
		filePath := filepath.Join(dirPath, file.Name())
		fmt.Println(filePath)

		// 如果要分辨文件和目录，可以使用 file.IsDir() 方法
		if file.IsDir() {
			fmt.Println("这是一个目录")
		} else {
			fmt.Println("这是一个文件")
		}
	}
}
