package main

import (
	"facebook_login/fb"
	"log"
	"os/exec"
)

func main() {

	//启动登录密文生成程序。无需使用时，你可在任务管理器结束
	log.Println(exec.Command("./encpass-server/server-win-x64.exe").Start())

	//初始化，使用代理方式，填空则不需要代理
	fb.Init("http://127.0.0.1:1080")

	//页面ID
	var documentId string

	//配置账户密码，建议使用邮箱类型用户名
	fbUsername := "example@gmail.com"
	fbPassword := "example"

	//登录网页版Facebook
	success, cookie := fb.Login(fbUsername, fbPassword)
	log.Println("登录返回：", success, cookie)

	//获取账户名称
	success, accountName := fb.GetAccountName(cookie)
	log.Println("获取账户名称返回：", success, accountName)

	//通过cookie取得后续操作所需的一些参数
	paramOk, graphqlParam := fb.GetGraphqlParam(cookie)
	log.Println("取参数返回：", paramOk, graphqlParam)

	//通过账户ID添加好友
	friendId := "100069534465329"
	documentId = "5554947034589513"
	success, response := fb.FriendRequest(graphqlParam, documentId, friendId)
	log.Println("添加好友返回：", success, response)

	//在个人主页发帖，最后一个参数代表帖子可见权限，0 所有人可见；1 仅自己可见
	documentId = "6403121713035238"
	success, postsUrl := fb.Post(graphqlParam, documentId, "我就发个贴！！！", 0)
	log.Println("发帖到个人主页返回：", success, postsUrl)

	//在指定的小组中发帖，前提是已加入
	documentId = "6403121713035238"
	groupId := "2239661079665553"
	success, groupPostsUrl := fb.PostGroup(graphqlParam, documentId, groupId, "小组的朋友们，大家晚上好啊！")
	log.Println("小组发帖返回：", success, groupPostsUrl)

	//获取当前账户下已加入的小组列表、
	number := 10 //获取数量
	documentId = "4001103780013203"
	success, groupsList := fb.GetMyGroups(graphqlParam, documentId, number)
	log.Println("获取小组列表返回：", success, groupsList)

}
