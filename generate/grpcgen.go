package generate

import (
	"adam/utils"
	"fmt"
	"github.com/myadamtest/logkit"
	"os"
	"os/exec"
)

const grpcexeUrl = "https://github.com/myadamtest/grpcexe.git"

func GrpcGenerate() {
	_, err := os.Stat("./grpcexe")
	if os.IsNotExist(err) { //fixme 把grpcexe文件放到path下
		cmd := exec.Command("git", "clone", grpcexeUrl)
		cmd.Stderr = os.Stdout
		err = cmd.Run()
		if err != nil {
			logkit.Infof("%s", err)
			return
		}
	}

	err = os.MkdirAll("./grpcservice", os.ModePerm)
	if err != nil {
		logkit.Errorf("%s", err)
		return
	}

	cmd := exec.Command("./grpcexe/protoc", "--go_out=plugins=grpc:./grpcservice/", "./protofile/helloworld2.proto")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = utils.GoModDownload("")

	//err = os.RemoveAll("./grpcexe")
	//if err != nil {
	//	logkit.Infof("%s", err)
	//	return
	//}
}
