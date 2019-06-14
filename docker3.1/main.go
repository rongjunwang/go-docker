package main

import (
	"fmt"
	"os"

	"./container"
	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `mydocker is a simple container runtime  implementation.
				The purpose of this project is to learn how docker works and how to write a docker by ourselvers 
				Enjoy it, just for fun.`

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		fmt.Println("test")
		//Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})

		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var runCommand = cli.Command{
	Name:  "run",
	Usage: `Create a contailner with namespace and cgroup limit mydocker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},

	/*
		这里是 run 命令执行的真正函数
		1.判断参数是否包含 command
		2.获取用户指定的command
		3.调用 Run function 去准备启动容器
	*/
	Action: func(context *cli.Context) error {
		fmt.Println("test runCommand")
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		container.Run(tty, cmd)
		return nil
	},
}

//这里定义了initCommand 的具体操作，此操作是内部方法，禁止外部调用
var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's pricess in container,Do not call it outside",
	/*
		1.执行传递过去的command参数
		2.执行容器初始化的操作
	*/
	Action: func(context *cli.Context) error {
		fmt.Println("test initCommand")
		log.Infof("init come on")
		cmd := context.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
