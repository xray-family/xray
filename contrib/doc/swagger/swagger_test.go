package swagger

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/swag"
	"github.com/xray-family/xray"
	httpAdapter "github.com/xray-family/xray/contrib/adapter/http"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func onANY(r *xray.Router, path string, handler xray.HandlerFunc) {
	r.GET(path, handler)
	r.POST(path, handler)
	r.PUT(path, handler)
	r.DELETE(path, handler)
	r.OnEvent(http.MethodOptions, path, handler)
}

type mockedSwag struct{}

func (s *mockedSwag) ReadDoc() string {
	return `{
}`
}

func TestWrapHandler(t *testing.T) {
	router := xray.New()

	router.GET("/:any", WrapHandler(swaggerFiles.Handler, URL("https://github.com/swaggo/gin-swagger")))

	assert.Equal(t, http.StatusOK, performRequest("GET", "/index.html", router).Code)
}

func TestWrapCustomHandler(t *testing.T) {
	router := xray.New()
	onANY(router, "/:any", CustomWrapHandler(&Config{}, swaggerFiles.Handler))

	w1 := performRequest(http.MethodGet, "/index.html", router)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, w1.Header()["Content-Type"][0], "text/html; charset=utf-8")

	assert.Equal(t, http.StatusInternalServerError, performRequest(http.MethodGet, "/doc.json", router).Code)

	doc := &mockedSwag{}
	swag.Register(swag.Name, doc)

	w2 := performRequest(http.MethodGet, "/doc.json", router)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, w2.Header()["Content-Type"][0], "application/json; charset=utf-8")

	// Perform body rendering validation
	w2Body, err := ioutil.ReadAll(w2.Body)
	assert.NoError(t, err)
	assert.Equal(t, doc.ReadDoc(), string(w2Body))

	w3 := performRequest(http.MethodGet, "/favicon-16x16.png", router)
	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Equal(t, w3.Header()["Content-Type"][0], "image/png")

	w4 := performRequest(http.MethodGet, "/swagger-ui.css", router)
	assert.Equal(t, http.StatusOK, w4.Code)
	assert.Equal(t, w4.Header()["Content-Type"][0], "text/css; charset=utf-8")

	w5 := performRequest(http.MethodGet, "/swagger-ui-bundle.js", router)
	assert.Equal(t, http.StatusOK, w5.Code)
	assert.Equal(t, w5.Header()["Content-Type"][0], "application/javascript")

	assert.Equal(t, http.StatusNotFound, performRequest(http.MethodGet, "/notfound", router).Code)

	assert.Equal(t, http.StatusMethodNotAllowed, performRequest(http.MethodPost, "/index.html", router).Code)

	assert.Equal(t, http.StatusMethodNotAllowed, performRequest(http.MethodPut, "/index.html", router).Code)
}

func TestDisablingWrapHandler(t *testing.T) {

	t.Run("", func(t *testing.T) {
		router := xray.New()
		disablingKey := "SWAGGER_DISABLE"
		router.GET("/simple/:any", DisablingWrapHandler(swaggerFiles.Handler, disablingKey))

		assert.Equal(t, http.StatusOK, performRequest(http.MethodGet, "/simple/index.html", router).Code)
		assert.Equal(t, http.StatusOK, performRequest(http.MethodGet, "/simple/doc.json", router).Code)

		assert.Equal(t, http.StatusOK, performRequest(http.MethodGet, "/simple/favicon-16x16.png", router).Code)
		assert.Equal(t, http.StatusNotFound, performRequest(http.MethodGet, "/simple/notfound", router).Code)
	})

	t.Run("", func(t *testing.T) {
		router := xray.New()
		disablingKey := "SWAGGER_DISABLE"
		_ = os.Setenv(disablingKey, "true")
		router.GET("/disabling/:any", DisablingWrapHandler(swaggerFiles.Handler, disablingKey))

		assert.Equal(t, http.StatusNotFound, performRequest(http.MethodGet, "/disabling/index.html", router).Code)
		assert.Equal(t, http.StatusNotFound, performRequest(http.MethodGet, "/disabling/doc.json", router).Code)
		assert.Equal(t, http.StatusNotFound, performRequest(http.MethodGet, "/disabling/oauth2-redirect.html", router).Code)
		assert.Equal(t, http.StatusNotFound, performRequest(http.MethodGet, "/disabling/notfound", router).Code)
	})
}

func TestDisablingCustomWrapHandler(t *testing.T) {
	disablingKey := "SWAGGER_DISABLE2"

	t.Run("", func(t *testing.T) {
		router := xray.New()
		router.GET("/simple/:any", DisablingCustomWrapHandler(&Config{}, swaggerFiles.Handler, disablingKey))
		assert.Equal(t, http.StatusOK, performRequest(http.MethodGet, "/simple/index.html", router).Code)
	})

	t.Run("", func(t *testing.T) {
		router := xray.New()
		_ = os.Setenv(disablingKey, "true")
		router.GET("/disabling/:any", DisablingCustomWrapHandler(&Config{}, swaggerFiles.Handler, disablingKey))
		assert.Equal(t, http.StatusNotFound, performRequest(http.MethodGet, "/disabling/index.html", router).Code)
	})

}

func performRequest(method, target string, router *xray.Router) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, bytes.NewBufferString(""))
	w := httptest.NewRecorder()
	httpAdapter.NewAdapter(router).ServeHTTP(w, r)
	return w
}

func TestURL(t *testing.T) {
	cfg := Config{}

	expected := "https://github.com/swaggo/http-swagger"
	configFunc := URL(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.URL)
}

func TestDocExpansion(t *testing.T) {
	var cfg Config

	expected := "list"
	configFunc := DocExpansion(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.DocExpansion)

	expected = "full"
	configFunc = DocExpansion(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.DocExpansion)

	expected = "none"
	configFunc = DocExpansion(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.DocExpansion)
}

func TestDeepLinking(t *testing.T) {
	var cfg Config
	assert.Equal(t, false, cfg.DeepLinking)

	configFunc := DeepLinking(true)
	configFunc(&cfg)
	assert.Equal(t, true, cfg.DeepLinking)

	configFunc = DeepLinking(false)
	configFunc(&cfg)
	assert.Equal(t, false, cfg.DeepLinking)

}

func TestDefaultModelsExpandDepth(t *testing.T) {
	var cfg Config

	assert.Equal(t, 0, cfg.DefaultModelsExpandDepth)

	expected := -1
	configFunc := DefaultModelsExpandDepth(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.DefaultModelsExpandDepth)

	expected = 1
	configFunc = DefaultModelsExpandDepth(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.DefaultModelsExpandDepth)
}

func TestInstanceName(t *testing.T) {
	var cfg Config

	assert.Equal(t, "", cfg.InstanceName)

	expected := swag.Name
	configFunc := InstanceName(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.InstanceName)

	expected = "custom_name"
	configFunc = InstanceName(expected)
	configFunc(&cfg)
	assert.Equal(t, expected, cfg.InstanceName)
}

func TestPersistAuthorization(t *testing.T) {
	var cfg Config
	assert.Equal(t, false, cfg.PersistAuthorization)

	configFunc := PersistAuthorization(true)
	configFunc(&cfg)
	assert.Equal(t, true, cfg.PersistAuthorization)

	configFunc = PersistAuthorization(false)
	configFunc(&cfg)
	assert.Equal(t, false, cfg.PersistAuthorization)
}

func TestOauth2DefaultClientID(t *testing.T) {
	var cfg Config
	assert.Equal(t, "", cfg.Oauth2DefaultClientID)

	configFunc := Oauth2DefaultClientID("default_client_id")
	configFunc(&cfg)
	assert.Equal(t, "default_client_id", cfg.Oauth2DefaultClientID)

	configFunc = Oauth2DefaultClientID("")
	configFunc(&cfg)
	assert.Equal(t, "", cfg.Oauth2DefaultClientID)
}
