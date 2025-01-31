package routes

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"platform/db"
	"platform/log"
	"platform/middleware"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type Data struct {
	User    *db.User
	Flashes []middleware.Flash
}

type DataInterface interface {
	UpdateFlashes([]middleware.Flash)
}

const MAX_PASSWORD_LENGTH = 6
const CHALLENGE_FILES_DIR = "./files"

var USERNAME_REGEX = regexp.MustCompile(`[0-9a-zA-Z_!@#€\-&+]{4,32}`)

func (d *Data) UpdateFlashes(flashes []middleware.Flash) {
	d.Flashes = flashes
}

func is_visible_category(chals []db.Challenge) bool {
	for _, chal := range chals {
		if !chal.Hidden {
			return true
		}
	}
	return false
}

func getTemplate(ctx *middleware.Ctx, page string) *template.Template {
	// TODO: parse all templates at startup
	funcMap := template.FuncMap{
		"is_visible_category": is_visible_category,
		"split":               func(sep string, s string) []string { return strings.Split(s, sep) },
		"last":                func(s []string) string { return s[len(s)-1] },
		"inc":                 func(i int) int { return i + 1 },
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFiles("templates/base.html", fmt.Sprintf("templates/%s.html", page))
	if err != nil {
		ctx.InternalError(fmt.Errorf("error parsing template %v", err))
		return nil
	}
	return tmpl
}

func executeTemplate(ctx *middleware.Ctx, tmpl *template.Template, data DataInterface) {
	flashes, err := ctx.GetFlashes()
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting flashes: %v", err))
		return
	}
	data.UpdateFlashes(flashes)

	err = tmpl.ExecuteTemplate(ctx.Writer, "base", data)
	if err != nil {
		ctx.InternalError(fmt.Errorf("error executing template %v", err))
		return
	}
}

func unzipFile(tmpFileDir string, file *zip.File) error {
	tmpFile, err := file.Open()
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	f, err := os.OpenFile(tmpFileDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, tmpFile)
	if err != nil {
		return err
	}

	return nil
}

func unzip(chalName string, file multipart.File, header *multipart.FileHeader) (string, error) {
	if !strings.HasSuffix(header.Filename, ".zip") {
		return "", fmt.Errorf("file must be a zip archive")
	}

	zreader, err := zip.NewReader(file, header.Size)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	_, err = hash.Write([]byte(chalName))
	if err != nil {
		return "", err
	}
	fileDir := fmt.Sprintf("%s/%x/", CHALLENGE_FILES_DIR, hash.Sum(nil))

	err = os.MkdirAll(fileDir, 0755)
	if err != nil {
		return "", err
	}

	files := ""
	for i, file := range zreader.File {
		tmpFileDir := fmt.Sprintf("%s/%s", fileDir, file.Name)

		err := unzipFile(tmpFileDir, file)
		if err != nil {
			return "", err
		}

		if i > 0 {
			files += ","
		}
		files += tmpFileDir
	}

	return files, nil
}

func getChallFromForm(ctx *middleware.Ctx) *db.Challenge {
	id := ctx.FormValue("id")
	name := ctx.FormValue("name")
	flag := ctx.FormValue("flag")
	maxPoints := ctx.FormValue("points")
	category := ctx.FormValue("category")
	difficulty := ctx.FormValue("difficulty")
	description := ctx.FormValue("description")
	hint1 := ctx.FormValue("hint1")
	hint2 := ctx.FormValue("hint2")
	host := ctx.FormValue("host")
	port := ctx.FormValue("port")
	isHidden := ctx.FormValue("is_hidden")
	isExtra := ctx.FormValue("is_extra")

	ID, err := strconv.Atoi(id)
	if err != nil {
		ctx.InternalError(fmt.Errorf("error converting ID to int: %v", err))
		return nil
	}

	points, err := strconv.Atoi(maxPoints)
	if err != nil {
		ctx.InternalError(fmt.Errorf("error converting points to int: %v", err))
		return nil
	}

	var hidden, extra bool
	if isHidden == "on" {
		hidden = true
	}
	if isExtra == "on" {
		extra = true
	}

	chal := &db.Challenge{
		ID:          ID,
		Name:        name,
		Description: description,
		Difficulty:  difficulty,
		Points:      points,
		MaxPoints:   points,
		Solves:      0,
		Host:        host,
		Port:        port,
		Category:    category,
		Files:       "",
		Flag:        flag,
		Hint1:       hint1,
		Hint2:       hint2,
		Hidden:      hidden,
		IsExtra:     extra,
	}
	return chal
}

func isChallengeValid(chal *db.Challenge) error {
	if len(chal.Name) < 1 || len(chal.Name) > 128 {
		return fmt.Errorf("name must be between 1 and 128 characters")
	}
	if chal.Category == "" {
		return fmt.Errorf("category must be set")
	}
	if chal.Difficulty == "" {
		return fmt.Errorf("difficulty must be set")
	}
	if len(chal.Flag) < 1 || len(chal.Flag) > 128 {
		return fmt.Errorf("flag must be between 1 and 128 characters")
	}
	if chal.MaxPoints < 0 {
		return fmt.Errorf("points must be positive")
	}

	chal.Description = strings.TrimSpace(chal.Description)
	return nil
}

func isNewChallengeValid(chal *db.Challenge) error {
	if err := isChallengeValid(chal); err != nil {
		return err
	}
	exists, err := db.ChallengeExistsName(chal.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("challenge already exists")
	}
	if err := db.FlagExists(chal.Flag); err != nil {
		return err
	}
	return nil
}

func createChallenge(ctx *middleware.Ctx, chal *db.Challenge) {
	err := isNewChallengeValid(chal)
	if err != nil {
		ctx.AddFlash(err.Error())
		ctx.Redirect("/admin", http.StatusSeeOther)
		return
	}

	err = db.CreateChallenge(chal)
	if err != nil {
		log.Errorf("Error creating challenge: %v", err)
		ctx.AddFlash("Error creating challenge")
		ctx.Redirect("/admin", http.StatusSeeOther)
		return
	}

	ctx.AddFlash("Challenge created successfully", "success")
	ctx.Redirect("/admin", http.StatusSeeOther)
}

func renameChallenge(chal *db.Challenge) error {
	oldName, err := db.GetChallengeName(chal.ID)
	if err != nil {
		return err
	}

	if oldName == chal.Name {
		return nil
	}

	exists, err := db.ChallengeExistsName(chal.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("can't rename, challenge already exists")
	}

	hash := sha256.New()
	_, err = hash.Write([]byte(chal.Name))
	if err != nil {
		return err
	}
	newDir := fmt.Sprintf("%s/%x/", CHALLENGE_FILES_DIR, hash.Sum(nil))

	hash = sha256.New()
	_, err = hash.Write([]byte(oldName))
	if err != nil {
		return err
	}
	oldDir := fmt.Sprintf("%s/%x/", CHALLENGE_FILES_DIR, hash.Sum(nil))

	err = os.Rename(oldDir, newDir)
	if err != nil {
		return err
	}

	return nil
}

func extractChallengeFiles(ctx *middleware.Ctx, chal *db.Challenge) bool {
	file, header, err := ctx.FormFile("files")
	if err != nil {
		ctx.InternalError(fmt.Errorf("error getting file: %v", err))
		return false
	}
	defer file.Close()

	chal.Files, err = unzip(chal.Name, file, header)
	if err != nil {
		log.Errorf("Error unzipping file: %v", err)
		ctx.AddFlash("Error unzipping file")
		ctx.Redirect("/admin", http.StatusSeeOther)
		return false
	}
	return true
}

func deleteChallengeFiles(name string) error {
	hash := sha256.New()
	_, err := hash.Write([]byte(name))
	if err != nil {
		return err
	}
	fileDir := fmt.Sprintf("%s/%x/", CHALLENGE_FILES_DIR, hash.Sum(nil))

	err = os.RemoveAll(fileDir)
	if err != nil {
		return err
	}

	return nil
}

func deleteChallenge(ctx *middleware.Ctx, name string) bool {
	err := deleteChallengeFiles(name)
	if err != nil {
		log.Errorf("Error deleting challenge files: %s: %v", name, err)
		ctx.AddFlash("Error deleting challenge")
		ctx.Redirect("/admin", http.StatusSeeOther)
		return false
	}

	err = db.DeleteChallenge(name)
	if err != nil {
		log.Errorf("Error deleting challenge from DB: %s: %v", name, err)
		ctx.AddFlash("Error deleting challenge")
		ctx.Redirect("/admin", http.StatusSeeOther)
		return false
	}

	return true
}
