package etcdv3

import (
	"context"
	"fmt"
	"github.com/erDong01/micro-kit/tools"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

const (
	uuid_dir = "uuid/"
)

type Snowflake struct {
	id      int64
	client  *clientv3.Client
	lease   clientv3.Lease
	leaseId clientv3.LeaseID
}

func (this *Snowflake) Key() string {
	return uuid_dir + fmt.Sprintf("%d", this.id)
}

func (this *Snowflake) Run() {
	for {
	TrySet:
		key := this.Key()
		tx := this.client.Txn(context.Background())
		leaseResp, err := this.lease.Grant(context.Background(), 60)
		if err != nil {
			goto TrySet
		}
		this.leaseId = leaseResp.ID
		tx.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Then(clientv3.OpPut(key, "", clientv3.WithLease(this.leaseId))).Else()
		txnRes, err := tx.Commit()
		if err != nil || !txnRes.Succeeded {
			resp, err := this.client.Get(context.Background(), uuid_dir)
			if err == nil && (resp != nil && resp.Kvs != nil) {
				Ids := [tools.WorkeridMax + 1]bool{}
				for _, v := range resp.Kvs {
					Id := tools.Int(string(v.Value[len(uuid_dir)+1]))
					Ids[Id] = true
				}
				for i, v := range Ids {
					if v == false {
						this.id = int64(i) & tools.WorkeridMax
						goto TrySet
					}
				}
			}
			this.id++
			this.id = this.id & tools.WorkeridMax
			goto TrySet
		}
		tools.UUID.Init(this.id)
		//保持ttl
	TryTTL:
		_, err = this.lease.KeepAliveOnce(context.Background(), this.leaseId)
		if err != nil {
			goto TrySet
		} else {
			time.Sleep(time.Second * 10)
			goto TryTTL
		}
	}
}

func (this *Snowflake) Init(endpoints []string) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	lease := clientv3.NewLease(etcdClient)
	this.id = 0
	this.client = etcdClient
	this.lease = lease
	this.Start()
}

func (this *Snowflake) Start() {
	go this.Run()
}
