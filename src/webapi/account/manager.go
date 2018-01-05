package account

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/go-playground/pool.v3"

	obean "webapi/bean"
	pkgConfig "webapi/config"
)

const (
	kBGSaveIntervalSecond = 60
)

type Manager struct {
	loginCounter int32 // 上次登陆记数
	lastSaveTS   int64
	mux          sync.RWMutex
	roles        map[uint32]*Role
	taskPool     pool.Pool
}

func New() *Manager {
	return &Manager{
		roles:    make(map[uint32]*Role),
		taskPool: pool.NewLimited(10),
	}
}

func (self *Manager) GetOnlineNum() int {
	self.mux.RLock()
	defer self.mux.RUnlock()

	return len(self.roles)
}

func (self *Manager) GetLoginNum() int32 {
	return atomic.SwapInt32(&self.loginCounter, 0)
}

func (self *Manager) doRoleLoad(uid uint32) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		err, character := obean.LoadCharacter(uid)
		if err != nil {
			log.Printf("加载角色%d基本数据时出错：%v", uid, err)
			return nil, err
		}

		err, skins := obean.LoadSkins(uid)
		if err != nil {
			log.Printf("加载角色%d皮肤列表时出错: %v", uid, err)
			return nil, err
		}

		err, skinTasks := obean.LoadSkinTasks(uid)
		if err != nil {
			log.Printf("加载七夕活动%d皮肤进度列表时出错: %v", uid, err)
			return nil, err
		}

		err, notifies := obean.LoadNotifies(uid)
		if err != nil {
			log.Printf("加载角色%d通知列表时出错: %v", uid, err)
			return nil, err
		}

		err, payments := obean.LoadAndroidPayments(uid)
		if err != nil {
			log.Printf("加载角色%d支付订单时出错: %v", uid, err)
			return nil, err
		}

		roleob := NewRole(character, skins, skinTasks, notifies)
		roleob.money = obean.LoadMoney(uid)
		roleob.winReward = obean.LoadWinReward(uid)
		roleob.adReward = obean.LoadAdReward(uid)
		rank := obean.NewRank(roleob.Uid()).Read()
		if rank.Tiers > 0 {
			roleob.SetRank(rank)
		}

		// 加载离线期间收到的支付订单
		if payments != nil {
			n := len(payments)
			for i := 0; i < n; i++ {
				v := payments[i]
				self.OnPaid(roleob, v)
			}
		}

		return roleob, nil
	}
}

// 新建角色
func (self *Manager) CreateNewRole(uid uint32, arr []int32, testBattleId int32) {
	self.mux.Lock()
	defer self.mux.Unlock()

	if _, exists := self.roles[uid]; exists {
		log.Printf("[Error]新建角色%d时对象已存在", uid)
		return
	}

	err, character := obean.NewCharacter(uid)
	if err != nil {
		log.Printf("新建%d角色基本数据时出错: %v", uid, err)
		return
	}

	character.TestBattleID = testBattleId

	err, skins := obean.NewSkins(uid, arr, obean.SkinSrcInit)
	if err != nil {
		log.Printf("新建角色%d皮肤列表%v时出错: %v", uid, arr, err)
		return
	}

	roleob := NewRole(character, skins, nil, nil)
	roleob.money = obean.CreateMoney(uid)
	roleob.winReward = obean.CreateWinReward(uid)
	roleob.adReward = obean.CreateAdReward(uid)
	self.roles[uid] = roleob

	atomic.AddInt32(&self.loginCounter, 1)
}

func (self *Manager) GetRole(uid uint32) *Role {
	self.mux.RLock()
	defer self.mux.RUnlock()

	roleob, exists := self.roles[uid]
	if exists {
		return roleob
	}

	return nil
}

func (self *Manager) setRole(role *Role) *Role {
	self.mux.Lock()
	defer self.mux.Unlock()

	uid := role.Uid()

	oldRole, exists := self.roles[uid]
	if exists {
		return oldRole
	}

	self.roles[uid] = role

	return role
}

// 加载角色
func (self *Manager) LoadRole(uid uint32) *Role {
	roleob := self.GetRole(uid)
	if roleob == nil {
		roleob = &Role{ID: uid}
		self.setRole(roleob)
	}
	return roleob
}
//func (self *Manager) LoadRole(uid uint32) *Role {
//	roleob := self.GetRole(uid)
//	if roleob == nil {
//		task := self.taskPool.Queue(self.doRoleLoad(uid))
//		task.Wait()
//
//		if task.Error() != nil {
//			log.Println("account/manager:func (self *Manager) LoadRole(uid uint32) *Role ", task.Error())
//			return nil
//		}
//
//		atomic.AddInt32(&self.loginCounter, 1)
//
//		roleob = task.Value().(*Role)
//		self.setRole(roleob)
//	}
//	return roleob
//}

func (self *Manager) Kick(uid uint32) bool {
	self.mux.Lock()
	defer self.mux.Unlock()

	roleob, exists := self.roles[uid]
	if !exists {
		return false
	}

	roleob.Save()
	delete(self.roles, uid)

	return true
}

func (self *Manager) BGSave() {
	self.mux.Lock()
	defer self.mux.Unlock()

	//for _, roleob := range self.roles {
	//	roleob.Save()
	//}
}

func (self *Manager) OnHeartBeat() {
	self.lastSaveTS++
	if self.lastSaveTS > kBGSaveIntervalSecond {
		self.lastSaveTS = 0
		self.BGSave()
		self.cleanRole()
	}
}

func (self *Manager) cleanRole() *Manager {
	self.mux.Lock()
	defer self.mux.Unlock()
	t := time.Now().Add(-10 * time.Minute).UnixNano()
	for k, v := range self.roles {
		if v.GetLastActionTime() <= t {
			delete(self.roles, k)
		}
	}
	return self
}

// 资源结算
// 支持皮肤、金币等
func (self *Manager) Reward(roleob *Role, resources map[int32]int32, src int) map[int32]int32 {
	return self.MultipleReward(roleob, resources, 1, src)
}

func (self *Manager) MultipleReward(roleob *Role, resources map[int32]int32, multiple int32, src int) map[int32]int32 {
	results := make(map[int32]int32)

	for id, count := range resources {
		if count <= 0 {
			continue
		}

		switch {
		case id < 1000:
			s := skinSrcFromReward(src)
			if ok := roleob.AddSkin(id, s); !ok {
				continue
			}
			results[id] = 1
		case 5001 == id: // 金币
			roleob.GoldAdd(count * multiple)
			results[id] = count * multiple
		case 5002 == id: // 钻石
			roleob.DiamondAdd(count*multiple, src)
			results[id] = count * multiple
		default:
			log.Printf("未知的资源编号: %d", id)
			continue
		}
	}

	return results
}

func skinSrcFromReward(src int) int {
	switch src {
	case obean.DiamondChangeError:
		return obean.SkinSrcError
	case obean.DiamondChangePay:
		return obean.SkinSrcPay
	case obean.DiamondChangeTaskReward:
		return obean.SkinSrcTaskReward
	case obean.DiamondChangeCDKey:
		return obean.SkinSrcCDKey
	case obean.DiamondChangeGM:
		return obean.SkinSrcGM
	case obean.DiamondChangeBuySkin:
		return obean.SkinSrcDiamond
	case obean.DiamondChangeBuyGold:
		return obean.SkinSrcDiamond
	default:
		return obean.SkinSrcError
	}
}

func (self *Manager) OnPaid(role *Role, bean *obean.AndroidPayment) {
	iap := pkgConfig.Singleton().GetFYSDKIAP(bean.SKU)
	if iap == nil {
		log.Printf("FYSDK SKU配置不存在: %s", bean.SKU)
		return
	}

	if (iap.ProductPrice * 100) != bean.Amount {
		log.Printf("FYSDK支付金额与配置不符[%d][%s]: %d != %d", role.Uid(), bean.SKU, bean.Amount, (iap.ProductPrice * 100))
		return
	}

	isFirstRecharge := iap.IsFirstRedeem && role.IsFirstRecharge()
	multiple := int32(1)
	if isFirstRecharge {
		multiple = 2
	}

	gains := self.MultipleReward(role, iap.ItemInfo, multiple, obean.DiamondChangePay)
	if iap.AdsMove {
		role.RemoveAD()
	}

	if isFirstRecharge {
		id := pkgConfig.Singleton().GetInt(pkgConfig.KeyFirstRechargeSkin)

		if ok := role.AddSkin(id, obean.SkinSrcPay); ok {
			gains[id] = 1
		}

		role.CleanFirstRechargeFlag()
	}

	notify, err := obean.AddAndroidPaymentNotify(bean.OrderID, bean.UID, role.IsRemoveAD(), gains, role.SerializeMoney())
	if err != nil {
		return
	}

	role.AddNotify(notify)
}
