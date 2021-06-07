package account

import (
	"context"
	"fmt"
	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"time"
)

var (
	ACCOUNTMGR AccountMgr
)

type (
	AccountDB struct {
		AccountId   int64  `sql:"primary;name:account_id"` //主键
		AccountName string `sql:"name:account_name"`
		Status      int    `sql:"name:status"`
		LoginTime   int64  `sql:"datetime;name:login_time"`  //日期
		LogoutTime  int64  `sql:"datetime;name:logout_time"` //日期
		LoginIp     string `sql:"name:login_ip"`
	}
	Account struct {
		AccountDB
	}
	IAccount interface {
		CheckLoginTime() bool
		UpdateAccountLogoutTime()
	}
	AccountMgr struct {
		actor.Actor

		m_AccountMap     map[int64]*Account
		m_AccountNameMap map[string]*Account
	}
)

func (this *Account) CheckLoginTime() bool {
	return false
}

func (this *Account) UpdateAccountLogoutTime() {
	this.LogoutTime = time.Now().Unix()
	//db
}
func (this *AccountMgr) Init(num int) {
	this.Actor.Init(1000)
	this.m_AccountMap = make(map[int64]*Account)
	this.m_AccountNameMap = make(map[string]*Account)
	//this.RegisterTimer(1000 * 1000 * 1000, this.Update)//定时器
	//账号登录处理
	this.RegisterCall("Account_Login", func(ctx context.Context, accountName string, accountId int, socketId int, id int) string {
		rpcHead := ctx.Value("rpcHead").(rpc3.RpcHead)
		client := network.SocketServer.GetClientById(rpcHead.SocketId)
		fmt.Println("AddAccount222:", rpcHead.SocketId, socketId, id)
		head := rpc3.RpcHead{Code: 200, Msg: "ok", ActorName: "Account"}
		byteD := rpc.Marshal(head, "Account_Login", "test", 2, 4, 9)
		client.DoSend(byteD)
		return "11109999"
	})
	this.Actor.Start()
}

func (this *AccountMgr) AddAccount(accountId int) *Account {
	fmt.Sprintf("AddAccount111:", accountId)

	return nil
}