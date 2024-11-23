package c

import (
	"fmt"
	"log"
	"time"

	"github.com/kr/pretty"
	"github.com/vault-thirteen/JSON-RPC-M1"
	"github.com/vault-thirteen/TR1/src/models/rpc"
)

type MetaData struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Forum struct {
	MetaData
	Id      int `gorm:"primarykey"`
	Name    string
	Threads []Thread `gorm:"foreignKey:ForumId"`
}

type Thread struct {
	MetaData
	Id      int `gorm:"primarykey"`
	ForumId int
	Name    string
}

func (c *Controller) test(p *rpc.RegisterUserParams) (result *rpc.RegisterUserResult, re *jrm1.RpcError) {
	db := c.GetDb()

	var err error
	err = db.AutoMigrate(&Forum{})
	if err != nil {
		log.Println(err)
	}
	err = db.AutoMigrate(&Thread{})
	if err != nil {
		log.Println(err)
	}

	/**/
	var thread1 = Thread{Name: "Thread 1"}
	var thread2 = Thread{Name: "Thread 2"}
	var forum1 = Forum{Name: "Forum 1", Threads: []Thread{thread1, thread2}}
	db.Create(&forum1)

	var forums []Forum
	res := db.Limit(10).Offset(0).Find(&forums)
	if res.Error != nil {
		log.Println(res.Error)
	}
	fmt.Println(res.RowsAffected)
	fmt.Println(forums)
	/**/

	var forum = Forum{Id: 1}
	n := db.Model(&forum).Association("Threads").Count()
	fmt.Println(n)

	var thread3 = Thread{Name: "Thread 3"}
	err = db.Model(&forum).Association("Threads").Append(&thread3)
	if err != nil {
		log.Println(err)
	}

	n = db.Model(&forum).Association("Threads").Count()
	fmt.Println(n)

	db.Preload("Threads").First(&forum, 1)
	pretty.Println(forum)

	return nil, jrm1.NewRpcErrorByUser(rpc.ErrorCode_FeatureIsNotImplemented, rpc.ErrorMsg_FeatureIsNotImplemented, nil)
}
