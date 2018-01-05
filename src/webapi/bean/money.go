package bean

import (
	"log"
	"sync"
	"time"

	"github.com/astaxie/beego/orm"
)

type Money struct {
	mtx     sync.Mutex `orm:"-"`
	UID     uint32     `orm:"column(uid);pk"`
	Gold    int32      `orm:"column(gold)"`
	Diamond int32      `orm:"column(diamond)"`
	Updated time.Time  `orm:"auto_now;type(datetime);column(updated)"`
}

func NewMoney(uid uint32) *Money {
	return &Money{
		UID: uid,
	}
}

func (g *Money) Read() *Money {
	err := defaultOrm.Read(g)
	checkError("从数据库读取玩家金钱数据错误:", err)
	return g
}

func (g *Money) ChangeGold(gold int32) *Money {
	g.mtx.Lock()
	g.Gold = g.Gold + gold
	_, err := defaultOrm.Update(g)
	if nil != err {
		g.Gold = g.Gold - gold
		log.Println("数据库更新玩家金钱数据,错误:", err)
	}
	g.mtx.Unlock()
	return g
}

func (g *Money) GetGold() int32 {
	g.mtx.Lock()
	gold := g.Gold
	g.mtx.Unlock()
	return gold
}

func (g *Money) GetDiamond() int32 {
	g.mtx.Lock()
	d := g.Diamond
	g.mtx.Unlock()
	return d
}

func (g *Money) ChangeDiamond(d int32, src int) error {
	g.mtx.Lock()
	g.Diamond = g.Diamond + d
	_, err := defaultOrm.Update(g)
	if err == nil {
		e := NewBalance(g.UID, g.Diamond, d, src).Insert()
		if nil != e {
			log.Println("数据库更新玩家钻石变动出错", e)
		}
	} else {
		log.Println("数据库更新玩家金钱数据,错误", err)
		g.Diamond = g.Diamond - d
	}
	g.mtx.Unlock()
	return err
}

func CreateMoney(uid uint32) *Money {
	g := NewMoney(uid)
	_, err := defaultOrm.Insert(g)
	checkError("数据库插入玩家金钱数据,错误:", err)
	return g
}

func LoadMoney(uid uint32) *Money {
	g := NewMoney(uid)
	err := defaultOrm.Read(g)
	if err == orm.ErrNoRows {
		_, err1 := defaultOrm.Insert(g)
		checkError("数据库插入玩家金钱数据,错误:", err1)
	} else {
		checkError("数据库读取玩家金钱数据,错误:", err)
	}
	return g
}

const (
	DiamondChangeError      = iota //错误的方式
	DiamondChangePay               //渠道充值购买
	DiamondChangeTaskReward        //任务奖励
	DiamondChangeCDKey             //CDKey兑换
	DiamondChangeGM                //GM修改
	DiamondChangeBuySkin           //购买皮肤
	DiamondChangeBuyGold           //购买金币
)

type Balance struct {
	ID      uint64    `orm:"column(id);auto;pk"`
	Date    time.Time `orm:"auto_now;type(datetime);column(date)"`
	UID     uint32    `orm:"column(uid);index"`
	Change  int32     `orm:"column(change)"`
	Diamond int32     `orm:"column(diamond)"`
	Source  int       `orm:"column(source)"`
}

func NewBalance(uid uint32, diamond int32, change int32, source int) *Balance {
	return &Balance{
		UID:     uid,
		Change:  change,
		Diamond: diamond,
		Source:  source,
	}
}

func (b *Balance) Insert() error {
	_, err := defaultOrm.Insert(b)
	return err
}
