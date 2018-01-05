package bean

import (
	"testing"

	"github.com/astaxie/beego/orm"
	pkgDB "yunjing.me/phoenix/go-phoenix/storage"
)

func init() {
	pDB := pkgDB.NewDBManager(false, "root:12345@tcp(localhost:3306)/test_db?charset=utf8")
	orm.RegisterModel(new(Account), new(ArenaLeaderboard), new(SimRole))
	//pDB.EnableLog()
	pDB.Run()
}

func BenchmarkPreloadAccount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PreloadAccount()
	}
}
