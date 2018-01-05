package other

import (
	"net/http"
	"strings"

	pkgConf "webapi/config"
)

// 评价地址跳转
func HandleHomePage(skeleton interface{}, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	channel := strings.TrimSpace(r.FormValue("channel"))
	if channel == "" {
		http.Redirect(w, r, "http://eat.crazyant.com", http.StatusPermanentRedirect)
		return
	}

	pConfManager := pkgConf.Singleton()
	cc := pConfManager.GetChannel(channel)
	if cc != nil && cc.Comment != "" {
		http.Redirect(w, r, cc.Comment, http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "http://eat.crazyant.com", http.StatusPermanentRedirect)
}
