package middleware

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"platform/db"
	"platform/log"
	"strings"

	"github.com/gorilla/sessions"
)

type Ctx struct {
	Writer     http.ResponseWriter
	request    *http.Request
	session    *sessions.Session
	User       *db.User
	RequestURI string
	URL        *url.URL
}

type Flash struct {
	Message string
	Type    string
}

func InitCtx(w http.ResponseWriter, r *http.Request) (*Ctx, error) {
	var ctx Ctx

	err := ctx.Init(w, r)
	if err != nil {
		return nil, err
	}
	return &ctx, nil
}

func (c *Ctx) Init(w http.ResponseWriter, r *http.Request) error {
	c.Writer = w
	c.request = r
	c.RequestURI = r.RequestURI
	c.URL = r.URL

	session, err := store.Get(r, "session")
	c.session = session
	if err != nil {
		c.AddFlash("Error getting session")
		c.Redirect("/", http.StatusSeeOther)
		return fmt.Errorf("error getting session: %v", err)
	}

	c.loadUser()
	return nil
}

func (c *Ctx) Save() error {
	err := c.session.Save(c.request, c.Writer)
	if err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}
	return nil
}

func (c *Ctx) Redirect(url string, code int) {
	err := c.Save()
	if err != nil {
		c.InternalError(fmt.Errorf("error saving session when redirecting: %v", err))
	}
	http.Redirect(c.Writer, c.request, url, code)
}

func (c *Ctx) InternalError(err error) {
	log.Errorf("%v", err)
	http.Error(c.Writer, "Internal Server Error", http.StatusInternalServerError)
}

func (c *Ctx) Error(msg string, code int) {
	err := c.Save()
	if err != nil {
		c.InternalError(fmt.Errorf("error saving session when throwing error: %v", err))
	}
	http.Error(c.Writer, msg, code)
}

func (c *Ctx) WriteHeader(code int) {
	err := c.Save()
	if err != nil {
		c.InternalError(fmt.Errorf("error saving session when writing header: %v", err))
	}
	c.Writer.WriteHeader(code)
}

func (c *Ctx) Write(data []byte) {
	n, err := c.Writer.Write(data)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
	}
	if n != len(data) {
		log.Errorf("Error writing response: %v", err)
	}
}

func (c *Ctx) AddFlash(args ...string) {
	if len(args) < 1 || len(args) > 2 {
		return
	}

	var flashType string
	if len(args) == 1 {
		flashType = "danger"
	} else {
		flashType = args[1]
	}

	c.session.AddFlash(&Flash{args[0], flashType})
}

func (c *Ctx) GetFlashes() ([]Flash, error) {
	flashesObjs := c.session.Flashes()

	flashes := make([]Flash, len(flashesObjs))
	for i, flash := range flashesObjs {
		flashes[i] = *flash.(*Flash)
	}

	err := c.Save()
	if err != nil {
		return nil, fmt.Errorf("error saving session when getting flashes: %v", err)
	}

	return flashes, nil
}

func (c *Ctx) loadUser() {
	apiKeyObj, ok := c.session.Values["apikey"]
	if !ok {
		return
	}

	apiKey, ok := apiKeyObj.(string)
	if !ok {
		log.Errorf("Error casting apikey")
		return
	}

	user, err := db.GetUserByAPIKey(apiKey)
	if err != nil {
		log.Warningf("Error getting user by apikey: %v", err)
		return
	}

	c.User = user
}

func (c *Ctx) IsValid() error {
	val := c.session.Values["apikey"]
	apiKey, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid session")
	}

	if strings.HasPrefix(apiKey, db.INVALID_PREFIX) {
		return fmt.Errorf("invalid session")
	}

	return nil
}

func (c *Ctx) ExpireCookie() {
	c.session.Options.MaxAge = -1
}

func (c *Ctx) FormValue(key string) string {
	return c.request.FormValue(key)
}

func (c *Ctx) PathValue(key string) string {
	return c.request.PathValue(key)
}

func (c *Ctx) SetSessionValue(key interface{}, value interface{}) {
	c.session.Values[key] = value
}

func (c *Ctx) ParseMultipartForm() {
	err := c.request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.InternalError(fmt.Errorf("error parsing multipart form: %v", err))
	}
}

func (c *Ctx) MultipartForm() *multipart.Form {
	return c.request.MultipartForm
}

func (c *Ctx) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.request.FormFile(key)
}

func (c *Ctx) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Ctx) ServeFile(path string) {
	http.ServeFile(c.Writer, c.request, path)
}
