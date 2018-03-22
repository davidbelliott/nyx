package board

import (
	"fmt"
	"github.com/pressly/chi"
	"github.com/tidwall/buntdb"
	"github.com/davidbelliott/nyx/http/errw"
	"github.com/davidbelliott/nyx/http/middle"
	"github.com/davidbelliott/nyx/resources"
	_ "image/gif"
	_ "image/jpeg"
	"net/http"
)

func handleNewReply(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	err = r.ParseMultipartForm(4 * 1024 * 1024)
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	if !resources.VerifyCaptcha(r) {
		http.Redirect(w, r,
			fmt.Sprintf("/%s/%s/thread.html?err=wrong_captcha",
				chi.URLParam(r, "board"), chi.URLParam(r, "thread")),
			http.StatusSeeOther)
		return
	}

	var reply = &resources.Reply{}

	err = parseReply(r, reply)
	if err == trollThrottle {
		http.Redirect(w, r,
			fmt.Sprintf("/%s/%s/thread.html?err=trollthrottle",
				chi.URLParam(r, "board"), chi.URLParam(r, "thread")),
			http.StatusSeeOther)
		return
	} else if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	db := middle.GetDB(r)
	if err = db.Update(func(tx *buntdb.Tx) error {
		thread, err := resources.GetThread(tx, r.Host, reply.Board, reply.Thread)
		if err != nil {
			return err
		}
		return resources.NewReply(tx, r.Host, reply.Board, thread, reply, false)
	}); err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%d/thread.html", chi.URLParam(r, "board"), reply.Thread), http.StatusSeeOther)
}
