package generate

import (
	"adam/utils"
	"fmt"
	"github.com/bookrun-go/fileutils/fileinfo"
	"github.com/myadamtest/logkit"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

const grpcexeUrl = "https://github.com/myadamtest/grpcexe.git"
const grpcexeFilename = "grpcexe"

func GrpcGenerate() {
	err := generateBefore()
	if err != nil {
		logkit.Errorf("%s", err)
		return
	}

	fd, err := ioutil.ReadDir("./protofile/")
	if err != nil {
		logkit.Errorf("%s", err)
		return
	}

	// 生成grpc文件
	for i := 0; i < len(fd); i++ {
		if !fd[i].IsDir() && fileinfo.GetFileSuffix(fd[i].Name()) == ".proto" {
			grpcGenerateByFilename("./protofile/" + fd[i].Name())
		}
	}

	err = utils.GoModDownload("")
}

func generateBefore() error {
	// 检查生成grpc需要的可执行文件是否存在
	exist, err := checkMustExeExist()
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	// 不存在清除可执行文件。重置环境
	_ = cleanMustExe()

	// 克隆可执行文件
	return cloneMustExeFile()
}

func grpcGenerateByFilename(fileName string) {
	simplyName := fileinfo.GetFileSimpleName(fileName)

	err := os.MkdirAll(fmt.Sprintf("./grpcservice/pb/%s/", simplyName), os.ModePerm)
	if err != nil {
		logkit.Errorf("%s", err)
		return
	}

	cmd := exec.Command("protoc", fmt.Sprintf("--go_out=plugins=grpc:./grpcservice/pb/%s/", simplyName), "--proto_path=./protofile/", fileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func checkMustExeExist() (bool, error) {
	rootBin := fmt.Sprintf("%s/bin", runtime.GOROOT())
	_, err := os.Stat(fmt.Sprintf("%s/protoc.exe", rootBin))
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	_, err = os.Stat(fmt.Sprintf("%s/protoc-gen-go.exe", rootBin))
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func cleanMustExe() error {
	rootBin := fmt.Sprintf("%s/bin", runtime.GOROOT())
	err := os.Remove(fmt.Sprintf("%s/protoc.exe", rootBin))
	if err != nil {
		return err
	}
	err = os.Remove(fmt.Sprintf("%s/protoc-gen-go.exe", rootBin))
	if err != nil {
		return err
	}

	return nil
}

func cloneMustExeFile() error {
	cmd := exec.Command("git", "clone", grpcexeUrl)
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	//var err error
	rootBin := fmt.Sprintf("%s/bin", runtime.GOROOT())
	_, err = copyFile(fmt.Sprintf("%s/protoc.exe", rootBin), fmt.Sprintf("./%s/protoc.exe", grpcexeFilename))
	if err != nil {
		return err
	}

	_, err = copyFile(fmt.Sprintf("%s/protoc-gen-go.exe", rootBin), fmt.Sprintf("./%s/protoc-gen-go.exe", grpcexeFilename))
	if err != nil {
		return err
	}

	err = os.RemoveAll(fmt.Sprintf("./%s/", grpcexeFilename))
	if err != nil {
		return err
	}

	return nil
}

func copyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
