package middleware

import (
	"net/http"
	"platform/log"
)

func HandleFunc(pattern string, handler func(ctx *Ctx)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx, err := InitCtx(w, r)
		if err != nil {
			log.Errorf("Error initializing context: %v", err)
			return
		}

		handler(ctx)
	})
}

func AuthHandleFunc(pattern string, handler func(ctx *Ctx)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx, err := InitCtx(w, r)
		if err != nil {
			log.Errorf("Error initializing context: %v", err)
			return
		}

		if ctx.User == nil {
			ctx.AddFlash("You must be logged in to access that page")
			ctx.Redirect("/", http.StatusSeeOther)
			return
		}

		if ctx.IsValid() != nil {
			ctx.Redirect("/newpw", http.StatusSeeOther)
			return
		}

		handler(ctx)
	})
}

func AdminHandleFunc(pattern string, handler func(ctx *Ctx)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx, err := InitCtx(w, r)
		if err != nil {
			log.Errorf("Error initializing context: %v", err)
			return
		}

		if ctx.User == nil || !ctx.User.IsAdmin {
			ctx.AddFlash("You must be logged in to access that page")
			ctx.Redirect("/", http.StatusSeeOther)
			return
		}

		handler(ctx)
	})
}
