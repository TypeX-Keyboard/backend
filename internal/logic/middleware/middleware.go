package middleware

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"io"
	"keyboard-api-go/internal/consts"
	"keyboard-api-go/internal/service"
	"keyboard-api-go/internal/util/secure"
	"sort"
	"strings"

	"keyboard-api-go/internal/util/cache"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type sMiddleware struct {
}

func init() {
	service.RegisterMiddleware(New())
}

func New() service.IMiddleware {
	return &sMiddleware{}
}

func (s *sMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func (s *sMiddleware) Auth(r *ghttp.Request) {
	// r.SetCtxVar(consts.AUTH_USERID, int64(1))
	// r.Middleware.Next()
	// return
	whiteList := []string{
		"v1/login",
		"v1/register",
	}
	adminList := []string{
		"/api/commission/displayAdminWithdrawalRecords",
	}
	if isAdmin := garray.NewStrArrayFrom(adminList).Contains(r.Request.URL.Path); isAdmin {
		g.Log().Info(r.Context(), "admin route", r.Request.URL.Path)
		tokenString := r.GetHeader("token")
		if len(tokenString) > 0 {
			if tokenString != "your_token" {
				g.Log().Info(r.Context(), "Token header is missing")
				r.Response.WriteStatusExit(401, "Unauthorized")
				return
			}
			r.Middleware.Next()
			return
		}
	}

	if isPass := garray.NewStrArrayFrom(whiteList).Contains(r.Request.URL.Path); isPass {
		g.Log().Info(r.Context(), "whitelist routes")
		r.Middleware.Next()
		return
	}
	// JWT authentication
	tokenString := r.GetHeader("Authorization")
	if tokenString == "" {
		g.Log().Info(r.Context(), "missing Authorization")
		r.Response.WriteStatusExit(401, "Unauthorized")
		return
	}
	// FirebaseAuthClient := firebase.GetFirebaseAuthClient()
	cacheRes, err := cache.GetCache().Get(r.Context(), tokenString)
	if err != nil || cacheRes.IsNil() {
		g.Log().Info(r.Context(), "Invalid JWT token")
		fmt.Printf("error message：%v\n", err)
		r.Response.WriteStatusExit(401, "Unauthorized")
		return
	}
	r.Middleware.Next()
}

func (s *sMiddleware) Context(r *ghttp.Request) {
	uuidStr := r.GetQuery(consts.DeviceIdParamKey, "")
	os := r.GetQuery(consts.OS, "")
	timezone := r.GetQuery(consts.TIMEZONE, "UTC")
	language := r.GetQuery(consts.LANGUAGE, "EN")
	aes, err := cache.GetCache().Get(r.GetCtx(), cache.GetAesKeyCacheKey(uuidStr.String()))
	if err != nil {
		r.Response.WriteStatusExit(400, err.Error())
		return
	}
	fallback_aes_key := secure.GetMD5(uuidStr.String())
	r.SetCtxVar(consts.OS, strings.ToUpper(os.String()))
	r.SetCtxVar("uuid", uuidStr.String())
	r.SetCtxVar("uid", uuidStr.String()) // 测试用
	r.SetCtxVar("enable_encrypt", uuidStr.String() != "")
	r.SetCtxVar("request_id", uuid.New().String())
	r.SetCtxVar("aes_key", aes.String())
	r.SetCtxVar("fallback_aes_key", fallback_aes_key)
	r.SetCtxVar(consts.TIMEZONE, timezone.String())
	r.SetCtxVar(consts.LANGUAGE, language.String())
	r.SetCtxVar(consts.DeviceIdParamKey, uuidStr.String())
	r.Middleware.Next()
}

func (s *sMiddleware) Encrypted(r *ghttp.Request) {
	r.Middleware.Next()
	if r.Response.BufferLength() > 0 {
		key := r.GetCtx().Value("aes_key")
		if key == nil || key.(string) == "" {
			r.Response.WriteStatusExit(419, "Signature error")
			return
		}
		content := r.Response.BufferString()
		g.Log().Infof(r.GetCtx(), "key: %v response body: %v", key, content)

		encryptedContent, err := secure.ECBEncrypt([]byte(content), []byte(key.(string)))
		if err != nil {
			r.Response.WriteStatusExit(419, err.Error())
			return
		}
		base64Content := base64.StdEncoding.EncodeToString(encryptedContent)
		//g.Log().Infof(r.GetCtx(), "body: %v", base64Content)
		r.Response.ClearBuffer()
		r.Response.Write([]byte(base64Content))
	}

}

func (s *sMiddleware) Decrypted(r *ghttp.Request) {
	if enableEncrypt := r.GetCtx().Value("enable_encrypt"); enableEncrypt != nil && enableEncrypt.(bool) {
		if err := s.checkRequest(r); err != nil {
			r.Response.WriteStatusExit(419, err.Error())
			return
		}
		body := r.GetBody()
		if len(body) != 0 {
			key := r.GetCtx().Value("aes_key")
			if key == nil || key.(string) == "" {
				r.Response.WriteStatusExit(419, "Signature error")
				return
			}
			g.Log().Infof(r.GetCtx(), "key: %v body: %v", key, string(body))
			origData, err := base64.StdEncoding.DecodeString(string(body))
			if err != nil {
				r.Response.WriteStatusExit(419, err.Error())
				return
			}
			res, err := secure.ECBDecrypt(origData, []byte(key.(string)))
			if err != nil {
				r.Response.WriteStatusExit(419, err.Error())
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(res))
			r.ReloadParam()
			g.Log().Infof(r.GetCtx(), "body: %v", string(r.GetBody()))
		}
		r.Middleware.Next()
	}
}

var requredQueryParams = []string{"timestamp", "os", "version", "device", "device_id", "network", "timezone", "language"}

func (s *sMiddleware) checkRequest(r *ghttp.Request) error {
	for _, param := range requredQueryParams {
		if r.GetQuery(param) == nil {
			return fmt.Errorf("Parameters are missing: %s", param)
		}
	}
	sign := r.GetQuery(consts.SignParamKey)
	needFallback := s.needAesKeyFallback(r)
	g.Log().Infof(r.Context(), "url:%s needFallback: %v", r.Request.URL.Path, needFallback)
	if needFallback {
		newSign := s.makeSign(r, r.GetCtx().Value("fallback_aes_key").(string))
		g.Log().Infof(r.Context(), "needFallback sign: %v newSign: %v", sign.String(), newSign)
		if sign.String() == newSign {
			r.SetCtxVar("aes_key", r.GetCtx().Value("fallback_aes_key").(string))
			return nil
		}
	}
	newSign := s.makeSign(r, r.GetCtx().Value("aes_key").(string))
	g.Log().Infof(r.Context(), "sign: %v newSign: %v", sign.String(), newSign)
	if sign.String() != newSign {
		return fmt.Errorf("Signature error")
	}
	return nil
}

func (s *sMiddleware) makeSign(r *ghttp.Request, aesKey string) string {
	g.Log().Infof(r.Context(), "checkSign url:%s aesKey: %v", r.Request.URL.Path, aesKey)
	queryParams := make([]string, 0)
	for k := range r.GetQueryMapStrStr() {
		if k == "sign" {
			continue
		}
		queryParams = append(queryParams, k)
	}
	sort.Strings(queryParams)
	temp := ""
	for _, param := range queryParams {
		temp += fmt.Sprintf("&%s=%v", param, r.GetQuery(param))
	}
	str2Sign := r.Request.URL.Path + temp[1:] + r.GetBodyString() + aesKey
	g.Log().Infof(r.Context(), "str2Sign: %v", str2Sign)
	return secure.GetMD5(str2Sign)
}

func (s *sMiddleware) needAesKeyFallback(r *ghttp.Request) bool {
	list := []string{
		"/api/v1/secure/acquire",
		"/api/v1/secure/submit",
		"/api/v1/system/settings",
	}
	return garray.NewStrArrayFrom(list).Contains(r.Request.URL.Path)
}
