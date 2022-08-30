package util

import (
	"casbin_demo"
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func testAuthzRequest(t *testing.T, router *gin.Engine, user string, path string, method string, code int) {
	r, _ := http.NewRequestWithContext(context.Background(), method, path, nil)
	r.SetBasicAuth(user, "123")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != code {
		t.Errorf("%s, %s, %s: %d, supposed to be %d", user, path, method, w.Code, code)
	}
}

func TestPathWildcard(t *testing.T) {
	router := gin.New()
	e, _ := casbin.NewEnforcer("rbac_model.conf", "test.csv")
	router.Use(casbin_demo.NewAuthorizer(e))
	router.Any("/*anypath", func(c *gin.Context) {
		c.Status(200)
	})

	testAuthzRequest(t, router, "alice", "/v1/h1", "GET", 200)
	testAuthzRequest(t, router, "alice", "/v1/h1", "POST", 200)
	testAuthzRequest(t, router, "alice", "/v2/h1", "GET", 403)
	testAuthzRequest(t, router, "alice", "/v3/h1", "GET", 403)
	testAuthzRequest(t, router, "bob", "/v2/h1", "GET", 200)
	testAuthzRequest(t, router, "bob", "/v3/code", "GET", 200)
	testAuthzRequest(t, router, "bob", "/v3/code", "DELETE", 403)
	testAuthzRequest(t, router, "ciel", "/v2/h1", "POST", 200)
	testAuthzRequest(t, router, "ciel", "/v3/h1", "DELETE", 200)
	testAuthzRequest(t, router, "ciel", "/v3/code", "GET", 200)
	testAuthzRequest(t, router, "ciel", "/v3/code", "DELETE", 200)
	testAuthzRequest(t, router, "ciel", "/v1/h1", "DELETE", 200)
	testAuthzRequest(t, router, "rookie", "/v3/code", "DELETE", 403)
	testAuthzRequest(t, router, "admin", "/v3/code", "DELETE", 200)
	testAuthzRequest(t, router, "rookie", "/v5/code", "DELETE", 200)

	_, err := e.DeleteRolesForUser("bob")
	if err != nil {
		t.Errorf("got error %v", err)
	}

	testAuthzRequest(t, router, "bob", "/v2/h1", "GET", 403)
}

func Test2(t *testing.T) {
	//Start()
}

func TestRoleApi(t *testing.T) {
	T()
}

func T() {
	Enf, _ = casbin.NewEnforcer("./rbac_model.conf", "test_copy.csv") //Adapter
	Enf.EnableAutoSave(true)
	Logger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{}), zapcore.NewMultiWriteSyncer(rotateLog(), zapcore.AddSync(os.Stdout)), zap.DebugLevel)).Sugar()

	//casbin_demo.Router = gin.New()
	//casbin_demo.router(casbin_demo.Router)
	//casbin_demo.Router.Run(":8080")
}
