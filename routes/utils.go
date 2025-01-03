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
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/sessions"
)

type Flash struct {
	Message string
	Type    string
}

type Data struct {
	User    *db.User
	Flashes []Flash
}

type DataInterface interface {
	UpdateFlashes([]Flash)
}

const MAX_PASSWORD_LENGTH = 6
const CHALLENGE_FILES_DIR = "./files"

// ! TODO: change
var store = sessions.NewCookieStore([]byte("GrazieDarioGrazieDarioGrazieDP_1"))
var USERNAME_REGEX = regexp.MustCompile(`[0-9a-zA-Z_!@#€\-&+]{4,32}`)

func (d *Data) UpdateFlashes(flashes []Flash) {
	d.Flashes = flashes
}

func getSession(w http.ResponseWriter, r *http.Request) (*sessions.Session, bool) {
	session, err := store.Get(r, "session")
	if err != nil {
		log.Errorf("Error getting session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, false
	}
	return session, true
}

func saveSession(w http.ResponseWriter, r *http.Request, s *sessions.Session) bool {
	err := s.Save(r, w)
	if err != nil {
		log.Errorf("Error saving session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}
	return true
}

func getSessionUser(s *sessions.Session) *db.User {
	val := s.Values["apikey"]
	apiKey, ok := val.(string)
	if !ok {
		return nil
	}
	user, err := db.GetUserByAPIKey(apiKey)
	if err != nil {
		log.Errorf("Error getting user by apikey: %v", err)
		return nil
	}
	return user
}

func addFlash(s *sessions.Session, args ...string) {
	if len(args) < 1 || len(args) > 2 {
		return
	}
	var flashType string
	if len(args) == 1 {
		flashType = "danger"
	} else {
		flashType = args[1]
	}
	s.AddFlash(&Flash{args[0], flashType})
}

func getFlashes(w http.ResponseWriter, r *http.Request, s *sessions.Session) []Flash {
	tmp := s.Flashes()
	flashes := make([]Flash, len(tmp))
	for i, flash := range tmp {
		flashes[i] = *flash.(*Flash)
	}
	err := s.Save(r, w)
	if err != nil {
		log.Errorf("Error saving session: %v", err)
	}
	return flashes
}

func is_visible_category(chals []*db.Challenge) bool {
	for _, chal := range chals {
		if !chal.Hidden {
			return true
		}
	}
	return false
}

func getTemplate(w http.ResponseWriter, page string) (*template.Template, error) {
	// TODO: parse all templates at startup
	funcMap := template.FuncMap{
		"is_visible_category": is_visible_category,
		"split":               func(sep string, s string) []string { return strings.Split(s, sep) },
		"last":                func(s []string) string { return s[len(s)-1] },
		"inc":                 func(i int) int { return i + 1 },
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFiles("templates/base.html", fmt.Sprintf("templates/%s.html", page))
	if err != nil {
		log.Errorf("Error parsing template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}
	return tmpl, nil
}

func executeTemplate(w http.ResponseWriter, r *http.Request, s *sessions.Session, tmpl *template.Template, data DataInterface) {
	data.UpdateFlashes(getFlashes(w, r, s))

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Errorf("Error executing template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

func getChallFromForm(w http.ResponseWriter, r *http.Request) (*db.Challenge, error) {
	id := r.FormValue("id")
	name := r.FormValue("name")
	flag := r.FormValue("flag")
	maxPoints := r.FormValue("points")
	category := r.FormValue("category")
	difficulty := r.FormValue("difficulty")
	description := r.FormValue("description")
	hint1 := r.FormValue("hint1")
	hint2 := r.FormValue("hint2")
	host := r.FormValue("host")
	port := r.FormValue("port")
	isHidden := r.FormValue("is_hidden")
	isExtra := r.FormValue("is_extra")

	ID, err := strconv.Atoi(id)
	if err != nil {
		log.Errorf("Error converting ID to int: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}

	points, err := strconv.Atoi(maxPoints)
	if err != nil {
		log.Errorf("Error converting points to int: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
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
	return chal, nil
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
	if err := db.CheckChallengeExists(chal.Name); err != nil {
		return err
	}
	if err := db.CheckFlagExists(chal.Flag); err != nil {
		return err
	}
	return nil
}

func createChallenge(w http.ResponseWriter, r *http.Request, s *sessions.Session, chal *db.Challenge) {
	err := isNewChallengeValid(chal)
	if err != nil {
		addFlash(s, err.Error())
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		return
	}

	err = db.CreateChallenge(chal)
	if err != nil {
		log.Errorf("Error creating challenge: %v", err)
		addFlash(s, "Error creating challenge")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		return
	}

	addFlash(s, "Challenge created successfully", "success")
	if saveSession(w, r, s) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func renameChallenge(chal *db.Challenge) error {
	oldName, err := db.GetChallengeName(chal.ID)
	if err != nil {
		return err
	}

	if oldName == chal.Name {
		return nil
	}

	if db.CheckChallengeExists(chal.Name) != nil {
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

func extractChallengeFiles(w http.ResponseWriter, r *http.Request, s *sessions.Session, chal *db.Challenge) error {
	file, header, err := r.FormFile("files")
	if err != nil {
		log.Errorf("Error getting file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	defer file.Close()

	chal.Files, err = unzip(chal.Name, file, header)
	if err != nil {
		log.Errorf("Error unzipping file: %v", err)
		addFlash(s, "Error unzipping file")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		return err
	}

	return nil
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

func deleteChallenge(w http.ResponseWriter, r *http.Request, s *sessions.Session, name string) error {
	err := deleteChallengeFiles(name)
	if err != nil {
		log.Errorf("Error deleting challenge files: %s: %v", name, err)
		addFlash(s, "Error deleting challenge")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		return err
	}

	err = db.DeleteChallenge(name)
	if err != nil {
		log.Errorf("Error deleting challenge from DB: %s: %v", name, err)
		addFlash(s, "Error deleting challenge")
		if saveSession(w, r, s) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		return err
	}

	return nil
}
