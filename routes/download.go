package routes

import (
	"net/http"
	"path/filepath"
	"platform/middleware"
	"strings"
)

func download(ctx *middleware.Ctx) {
	cleanPath := filepath.Clean("./" + ctx.URL.Path)
	if !strings.HasPrefix(cleanPath, "files") {
		ctx.Error("Invalid file path", http.StatusForbidden)
		return
	}

	fileName := filepath.Base(cleanPath)

	ctx.SetHeader("Content-Disposition", "attachment; filename="+fileName)
	ctx.SetHeader("Content-Type", "application/octet-stream")
	ctx.SetHeader("Content-Transfer-Encoding", "binary")

	ctx.ServeFile(cleanPath)
}
