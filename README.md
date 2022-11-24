# Facebook 网页爬虫工具

> 爬虫模拟实现Facebook的个人主页发帖、小组发帖、添加好友、获取已加入小组列表、获取账户名称。
> [![License](https://img.shields.io/badge/license-MIT-db5149.svg)](https://github.com/trry071/facebook-crawler-tools/blob/master/LICENSE)

## 使用

#### 引入 fb 爬虫工具包

```ssh
import ("facebook_login/fb")
```

## 调用例子

运行前请先配置好国外代理，否则将无法运行。

```go
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

	//页面ID，固定的
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

	//获取当前账户下已加入的小组列表
	number := 10 //获取数量
	documentId = "4001103780013203"
	success, groupsList := fb.GetMyGroups(graphqlParam, documentId, number)
	log.Println("获取小组列表返回：", success, groupsList)

}

```

## 注意

encpass-server 文件夹下的程序是用于生成登录密文的，它是开源的： [facebook-login-encpass](https://github.com/trry071/facebook-login-encpass)

## 打赏作者

开源不易，如果此项目对你有帮助，不妨可以打赏一下，金额不限，非常感谢！！！

<img src="https://s2.loli.net/2022/11/24/E4ernNydpBY3tCk.jpg" width = "20%" height = "20%"  /> <img src="https://s2.loli.net/2022/11/24/majpKl1g2q5O3GL.jpg" width = "20%" height = "20%"  />

## License

[MIT](LICENSE)  
这个项目是无人维护的。您可以使用它，但问题和拉取请求可能会被忽略。
