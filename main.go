package main

import (
	"GameBase2/app"
	"GameBase2/config"
	"GameBase2/room"
	"embed"
	R "tcg/res"
)

var (
	//go:embed res
	files         embed.FS
	StackRoom     *room.StackUIRoom
	PlayerManager *playerManager
	ActionManager *actionManager
	Info          *info
	Tip           *tip
	Round         *round
)

func main() {
	config.ViewSize = complex(1280, 720)
	config.Debug = true
	config.ShowFps = true
	config.Files = &files // 先使用内部资源 ，不存在  再寻找外部资源文件
	InitFont()
	app.Run(NewMainApp(), 1280, 720)
}

type MainApp struct {
	*app.App
}

// Init 必须先传入实例  初始化使用该方法
func NewMainApp() *MainApp {
	res := &MainApp{}
	res.App = app.NewApp()
	PlayerManager = NewPlayerManager()
	ActionManager = NewActionManager()
	temp := config.RoomFactory.LoadAndCreate(R.MAP.MAIN)
	temp.AddManager(PlayerManager)
	temp.AddManager(ActionManager)
	StackRoom = room.NewStackUIRoom(temp.(*room.Room))
	res.PushRoom(StackRoom)
	return res
}
