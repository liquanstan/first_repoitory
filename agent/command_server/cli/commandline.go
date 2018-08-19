package cli

import (
	"os/exec"
	"io/ioutil"
	"context"
	"time"
)

func ExecuteCli(args ...string) (string , error) {


	ctx,_:=context.WithTimeout(context.Background(),time.Second*5)         //设置5秒执行时间限定
	cmd := exec.CommandContext(ctx,args[0],args[1:]...)                   //将第一个参数作为shell，后面的作为具体参数

	stdout,err := cmd.StdoutPipe()
	if err != nil {
		return " ",err
	}

	err =cmd.Start()
	if err != nil {
		return " ",err
	}

	bytes,err := ioutil.ReadAll(stdout)
	if err != nil {
		return " ",err
	}

	err = cmd.Wait()
	if err != nil {

		return " ",err
	}

	return string(bytes),nil

}
