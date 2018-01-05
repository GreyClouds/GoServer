package bean

import (
	"log"
	"sync"
	"time"

	"runtime/debug"

	"github.com/astaxie/beego/orm"
)

const (
	runing = iota + 1
	closed
	buffLen    = 1024
	backlog    = 4096
	keeperScan = time.Second
)

const (
	Once = iota + 1
)

type Cron struct {
	flag int
	cron func()
}

func newCron(f int, c func()) *Cron {
	return &Cron{
		flag: f,
		cron: c,
	}
}

type IdleCron struct {
	exec bool
	id   int
	cron map[int]*Cron
}

func newIdleCron() *IdleCron {
	return &IdleCron{
		exec: true,
		id:   0,
		cron: make(map[int]*Cron),
	}
}

func (i *IdleCron) AddCron(c *Cron) {
	i.id++
	i.cron[i.id] = c
}

func (i *IdleCron) Execute() *IdleCron {
	i.id = 0
	for k, v := range i.cron {
		v.cron()
		if v.flag == Once {
			delete(i.cron, k)
		}
	}
	i.exec = false
	return i
}

type ModelKeeper struct {
	status         int
	dbOrm          orm.Ormer
	closeCh        chan struct{}
	closeOK        chan struct{}
	rankCh         chan *Rank
	rankHistoryCh  chan *RankHistory
	rankRewardCh   chan *RankReward
	rankRoomCh     chan *RankRoom
	rankListCh     chan *RankList
	normalRoomCh   chan *BattleRoom
	cronCh         chan *Cron
	bufRank        *sync.Map
	bufRankHistory *sync.Map
	bufRankReward  *sync.Map
	bufRankRoom    *sync.Map
	idleCron       *IdleCron
}

func NewModelKeeper() *ModelKeeper {
	return &ModelKeeper{
		status:         closed,
		dbOrm:          nil,
		closeCh:        make(chan struct{}),
		closeOK:        make(chan struct{}),
		rankCh:         make(chan *Rank, buffLen),
		rankHistoryCh:  make(chan *RankHistory, buffLen),
		rankRewardCh:   make(chan *RankReward, buffLen),
		rankRoomCh:     make(chan *RankRoom, buffLen),
		rankListCh:     make(chan *RankList, buffLen),
		normalRoomCh:   make(chan *BattleRoom, backlog),
		cronCh:         make(chan *Cron, 10),
		bufRank:        &sync.Map{},
		bufRankHistory: &sync.Map{},
		bufRankReward:  &sync.Map{},
		bufRankRoom:    &sync.Map{},
		idleCron:       newIdleCron(),
	}
}

func (m *ModelKeeper) Close() *ModelKeeper {
	close(m.closeCh)
	t := time.Now()
	<-m.closeOK
	long := time.Since(t)
	log.Println("关闭数据管理员,耗时:", long)
	return m
}

func (m *ModelKeeper) SendRankInsert(r *Rank) *ModelKeeper {
	if m.status == runing {
		c := r.clone()
		c.op = ins
		m.rankCh <- c
	}
	return m
}

func (m *ModelKeeper) SendRankUpdate(r *Rank) *ModelKeeper {
	if m.status == runing {
		c := r.clone()
		c.op = up
		m.rankCh <- c
	}
	return m
}

func (m *ModelKeeper) SendHistory(r *RankHistory) *ModelKeeper {
	if m.status == runing {
		m.rankHistoryCh <- r
	}
	return m
}

func (m *ModelKeeper) SendReward(r *RankReward) *ModelKeeper {
	if m.status == runing {
		m.rankRewardCh <- r
	}
	return m
}

func (m *ModelKeeper) SendRoomInsert(r *RankRoom) *ModelKeeper {
	if m.status == runing {
		r.op = ins
		m.rankRoomCh <- r
	}
	return m
}

func (m *ModelKeeper) SendRoomUpdate(r *RankRoom) *ModelKeeper {
	if m.status == runing {
		r.op = up
		m.rankRoomCh <- r
	}
	return m
}

func (m *ModelKeeper) SendRankList(r *RankList) *ModelKeeper {
	m.rankListCh <- r
	return m
}

func (m *ModelKeeper) SendBattleRoom(bean *BattleRoom) *ModelKeeper {
	m.normalRoomCh <- bean
	return m
}

func (m *ModelKeeper) SendCron(f int, c func()) *ModelKeeper {
	m.cronCh <- newCron(f, c)
	return m
}

func (m *ModelKeeper) SetORM(o orm.Ormer) *ModelKeeper {
	m.dbOrm = o
	return m
}

func (m *ModelKeeper) Run() *ModelKeeper {
	m.SetORM(orm.NewOrm())
	m.status = runing
	go m.run()
	return m
}

func (m *ModelKeeper) run() *ModelKeeper {
	defer func() {
		err := recover()
		if nil != err {
			log.Println("obean/manager.go:func (m *ModelKeeper) run(), Error:", err)
			debug.PrintStack()
		}
	}()
	t := time.NewTicker(keeperScan)
	for m.status == runing {
		select {
		case <-t.C:
			m.scan()
		case insRank := <-m.rankCh:
			m.handleRank(insRank)
		case history := <-m.rankHistoryCh:
			m.handleRankHistory(history)
		case reward := <-m.rankRewardCh:
			m.handleRankReward(reward)
		case room := <-m.rankRoomCh:
			m.handleRankRoom(room)
		case rankList := <-m.rankListCh:
			m.handleRankList(rankList)
		case cron := <-m.cronCh:
			m.handleCron(cron)
		case <-m.closeCh:
			m.handleClose()
		}
	}
	return m
}

func (m *ModelKeeper) scan() {
	m.flushBuffer()
	if m.idleCron.exec && time.Now().Hour() == 4 {
		m.idleCron.Execute()
	}
	if !m.idleCron.exec && time.Now().Hour() != 4 {
		m.idleCron.exec = true
	}
}

func (m *ModelKeeper) flushBuffer() {
	m.flushRank().flushHistory().flushReward().flushRankRoom()
	m.flushNormalBattleRoom()
}

func (m *ModelKeeper) flushNormalBattleRoom() {
	var empty bool
	arr := make([]*BattleRoom, 0, backlog)
	for i := 0; i < backlog; i++ {
		select {
		case bean := <-m.normalRoomCh:
			arr = append(arr, bean)
		default:
			empty = true
		}

		if empty {
			break
		}
	}

	if !empty {
		log.Printf("普通战斗记录缓存队列已满: %d", backlog)
	}

	if l := len(arr); l > 0 {
		// t := time.Now()
		n, err := m.dbOrm.InsertMulti(1, arr)
		if err != nil || int(n) != l {
			log.Println("批量写入普通战斗记录时出错: %v", err)
		}
		// log.Printf("批量写入普通战斗记录%d耗时: %v", l, time.Since(t))
	}
}

func (m *ModelKeeper) flushRank() *ModelKeeper {
	var buf []*Rank
	fun := func(k interface{}, v interface{}) bool {
		r, ok := v.(*Rank)
		if ok {
			buf = append(buf, r)
		}
		m.bufRank.Delete(k)
		return true
	}
	m.bufRank.Range(fun)
	l := len(buf)
	if l > 0 {
		n, err := m.dbOrm.InsertMulti(1, buf)
		if err != nil || int(n) != l {
			log.Println("批量写入rank错误, 应写入:", l, "实写入:", n, err)
		}
	}
	return m
}

func (m *ModelKeeper) flushHistory() *ModelKeeper {
	var buf []*RankHistory
	fun := func(k interface{}, v interface{}) bool {
		r, ok := v.(*RankHistory)
		if ok {
			buf = append(buf, r)
		}
		m.bufRankHistory.Delete(k)
		return true
	}
	m.bufRankHistory.Range(fun)
	l := len(buf)
	if l > 0 {
		n, err := m.dbOrm.InsertMulti(1, buf)
		if err != nil || int(n) != l {
			log.Println("批量写入rank_history错误, 应写入:", l, "实写入:", n, err)
		}
	}
	return m
}

func (m *ModelKeeper) flushReward() *ModelKeeper {
	var reward []*RankReward
	var item []*RankRewardItem
	fun := func(k interface{}, v interface{}) bool {
		r, ok := v.(*RankReward)
		if ok {
			reward = append(reward, r)
			if len(r.Reward) > 0 {
				item = append(item, r.Reward...)
			}
		}
		m.bufRankReward.Delete(k)
		return true
	}
	m.bufRankReward.Range(fun)
	l := len(reward)
	if l > 0 {
		nReward, errReward := m.dbOrm.InsertMulti(1, reward)
		if errReward != nil || int(nReward) != l {
			log.Println("批量写入rank_reward错误, 应写入:", l, "实写入:", nReward, errReward)
		}
		lenItem := len(item)
		if lenItem > 0 {
			nItem, errItem := m.dbOrm.InsertMulti(1, item)
			if errItem != nil {
				log.Println("批量写入批量写入rank_reward_info错误, 应写入:", lenItem, "实写入:", nItem, errItem)
			}
		}
	}
	return m
}

func (m *ModelKeeper) flushRankRoom() *ModelKeeper {
	var buf []*RankRoom
	fun := func(k interface{}, v interface{}) bool {
		r, ok := v.(*RankRoom)
		if ok {
			buf = append(buf, r)
		}
		m.bufRankRoom.Delete(k)
		return true
	}
	m.bufRankRoom.Range(fun)
	l := len(buf)
	if l > 0 {
		n, err := m.dbOrm.InsertMulti(1, buf)
		if nil != err || int(n) != l {
			log.Println("批量写入rank_room错误, 应写入:", l, "实写入:", n, err)
		}
	}
	return m
}

func (m *ModelKeeper) handleRank(r *Rank) {
	if r != nil {
		switch r.op {
		case ins:
			m.rankInsert(r)
		case up:
			m.rankUpdate(r)
		default:
		}
	}
}

func (m *ModelKeeper) rankInsert(r *Rank) {
	m.bufRank.Store(r.UID, r)
}

func (m *ModelKeeper) rankUpdate(r *Rank) {
	_, ok := m.bufRank.Load(r.UID)
	if ok {
		m.bufRank.Store(r.UID, r)
	} else {
		_, err := m.dbOrm.Update(r)
		checkError("数据库更新rank数据错误:", err)
	}
}

func (m *ModelKeeper) handleRankHistory(r *RankHistory) {
	if r != nil {
		m.bufRankHistory.Store(r.UID, r)
	}
}

func (m *ModelKeeper) handleRankReward(r *RankReward) {
	if r != nil {
		m.bufRankReward.Store(r.UID, r)
	}
}

func (m *ModelKeeper) handleRankRoom(r *RankRoom) {
	if r != nil {
		switch r.op {
		case ins:
			m.rankRoomInsert(r)
		case up:
			m.rankRoomUpdate(r)
		default:
		}
	}
}

func (m *ModelKeeper) rankRoomInsert(r *RankRoom) {
	m.bufRankRoom.Store(r.ID, r)
}

func (m *ModelKeeper) rankRoomUpdate(r *RankRoom) {
	_, ok := m.bufRankRoom.Load(r.ID)
	if ok {
		m.bufRankRoom.Store(r.ID, r)
	} else {
		_, err := m.dbOrm.Update(r)
		checkError("更新房间结果", err)
	}
}

func (m *ModelKeeper) handleRankList(r *RankList) {
	if r != nil {
		switch r.op {
		case up:
			_, err := m.dbOrm.Update(r)
			checkError("Update rankList", err)
		case ins:
			_, err := m.dbOrm.Insert(r)
			checkError("Insert rankList", err)
		case del:
			_, err := m.dbOrm.Delete(r)
			checkError("Delete rankList", err)
		default:
		}
	}
}

func (m *ModelKeeper) handleCron(c *Cron) {
	if nil != c {
		m.idleCron.AddCron(c)
	}
}

func (m *ModelKeeper) handleClose() {
	if m.status == runing {
		log.Println("bean.ModelKeeper 准备关闭...")
		m.closeChanel()
		m.flushChanel()
		m.flushBuffer()
		close(m.closeOK)
		log.Println("bean.RankManager 关闭 ...OK")
	}
	m.status = closed
}

func (m *ModelKeeper) closeChanel() {
	close(m.rankCh)
	close(m.rankHistoryCh)
	close(m.rankRewardCh)
	close(m.rankRoomCh)
	close(m.rankListCh)
	close(m.cronCh)
}

func (m *ModelKeeper) flushChanel() {
	for v := range m.rankCh {
		m.handleRank(v)
	}
	for v := range m.rankHistoryCh {
		m.handleRankHistory(v)
	}
	for v := range m.rankRewardCh {
		m.handleRankReward(v)
	}
	for v := range m.rankRoomCh {
		m.handleRankRoom(v)
	}
	for v := range m.rankListCh {
		m.handleRankList(v)
	}
}

var defaultKeeper = NewModelKeeper()

func DefaultKeeper() *ModelKeeper {
	return defaultKeeper
}

func (m *ModelKeeper) LoadRankRoom(roomUID uint32) (*RankRoom, bool) {
	r := &RankRoom{
		ID: roomUID,
	}
	v, ok := m.bufRankRoom.Load(roomUID)
	if ok {
		r = v.(*RankRoom)
	} else {
		err := m.dbOrm.Read(r)
		if err != nil {
			return nil, false
		}
	}
	r.op = up
	return r, true
}

func (m *ModelKeeper) LoadRank(uid uint32) (*Rank, bool) {
	r := &Rank{
		UID: uid,
	}
	v, ok := m.bufRank.Load(uid)
	if ok {
		r = v.(*Rank)
	} else {
		err := m.dbOrm.Read(r)
		if err != nil {
			return nil, false
		}
	}
	r.op = up
	return r, true
}
