package generate

import (
	"fmt"
	"github.com/bookrun-go/fileutils/fileinfo"
	"github.com/bookrun-go/parseutil/gofile"
	"github.com/deckarep/golang-set"
	"github.com/myadamtest/adam/utils"
	"github.com/myadamtest/logkit"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"text/template"
)

func GrpcGenerate() {
	err := generateBefore()
	if err != nil {
		logkit.Errorf("%s", err)
		return
	}

	fd, err := ioutil.ReadDir(protoBaseDir)
	if err != nil {
		logkit.Errorf("%s", err)
		return
	}

	serviceList := make([]*RpcServiceInfo, 0)
	// 生成grpc文件
	for i := 0; i < len(fd); i++ {
		if !fd[i].IsDir() && fileinfo.GetFileSuffix(fd[i].Name()) == ".proto" {
			stru, err := grpcGenerateByFilename(protoBaseDir + fd[i].Name())
			if err != nil {
				logkit.Errorf("generate[%s] code fail", fd[i].Name())
				return
			}
			if stru != nil && len(stru.Interfaces) > 0 {
				serviceList = append(serviceList, stru)
			}
		}
	}

	//生成grpc启动文件
	startTmpl, err := template.New("stp").Parse(grpcStartTemplate)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = os.Remove("./grpcservice/grpc_service.go")
	startFd, err := os.OpenFile("./grpcservice/grpc_service.go", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer startFd.Close()

	thisProjectName, _ := utils.GetProjectName()
	//packageList := make([]string,0)
	packageList := mapset.NewSet()
	for i := 0; i < len(serviceList); i++ {
		//"{{$v.ProjectName}}/grpcservice/pb/{{$v.PackageName}}"
		packageList.Add(fmt.Sprintf("\"%s/grpcservice/pb/%s\"", serviceList[i].ProjectName, serviceList[i].PackageName))
	}

	err = startTmpl.Execute(startFd, map[string]interface{}{"ProjectName": thisProjectName, "ServiceList": serviceList, "PackageList": packageList.ToSlice()})
	if err != nil {
		fmt.Println(err)
	}
	//生成grpc启动文件 结束

	err = utils.GoModTidy("")
	if err != nil {
		fmt.Println("gen grpc fail!", err)
		return
	}
	fmt.Println(fmt.Sprintf("success gen grpc file"))
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

func grpcGenerateByFilename(fileName string) (*RpcServiceInfo, error) {
	stru, err := GetProtoStruct(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//simplyName := fileinfo.GetFileSimpleName(fileName)
	//projectName, _ := utils.GetProjectName()

	err = os.MkdirAll(fmt.Sprintf("./grpcservice/pb/%s/", stru.PackageName), os.ModePerm)
	if err != nil {
		logkit.Errorf("%s", err)
		return nil, err
	}

	cmd := exec.Command("protoc", fmt.Sprintf("--go_out=plugins=grpc:./grpcservice/pb/%s", stru.PackageName), "--proto_path=./protofile/", fileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	generateImplement(stru)
	return stru, nil
}

// 生成实现模板
func generateImplement(stru *RpcServiceInfo) {
	if stru == nil || len(stru.Interfaces) == 0 {
		return
	}

	targetFileName := fmt.Sprintf("./grpcservice/%s.go", stru.FileName)

	_, err := os.Stat(targetFileName)
	if !os.IsNotExist(err) && err != nil {
		fmt.Println(err)
		return
	}

	tp := methodTemplate
	if os.IsNotExist(err) {
		tp = commonTemplate
	} else {
		err = gofile.IParseFactory.DoByFile(targetFileName)
		if err == nil { // 去重已经有的函数。
			fs := gofile.IParseFactory.GetRootFunctions()
			newInterface := make([]RpcServiceInterfaceInfo, 0)
			for _, i := range stru.Interfaces {
				tempI := &i
				for _, f := range fs {
					if i.Name == f.Name {
						tempI = nil
						break
					}

				}

				if tempI != nil {
					newInterface = append(newInterface, *tempI)
				}
			}
			stru.Interfaces = newInterface
		}
	}

	tmpl, err := template.New("tp").Parse(tp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fd, err := os.OpenFile(targetFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fd.Close()

	err = tmpl.Execute(fd, stru)
	if err != nil {
		fmt.Println(err)
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
