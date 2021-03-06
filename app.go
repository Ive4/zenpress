package main

import (
	"log"

	"github.com/insionng/vodka"
	"github.com/insionng/vodka/engine/fasthttp"
	m "github.com/insionng/vodka/middleware"
	"github.com/insionng/zenpress/handler"
	"github.com/vodka-contrib/pongor"
	"github.com/vodka-contrib/session"
	//_ "github.com/vodka-contrib/session/redis"
	"github.com/vodka-contrib/vodkapprof"
)

/*
func TokenHandler(self *vodka.Context) error {

	var claims = map[string]interface{}{}
	claims["username"] = "Insion"
	token, err := jwt.NewToken(key, claims)
	if err != nil {
		return err
	}
	// show the token use for test
	return self.String(http.StatusOK, "%s", token)
}
*/
// init 初始化Sesssion
func init() {
	opt := session.Options{"file", `{"cookieName":"vodkaSid","gclifetime":3600,"ProviderConfig":"./data/session"}`}
	if err := session.Setup(opt); err != nil {
		log.Fatalln("session errors:", err)
	}
}

func main() {

	v := vodka.New()

	v.Use(m.Logger())
	v.Use(m.Recover())
	v.Use(m.Gzip())
	v.Use(m.Secure())
	v.Use(m.BodyLimit("2M"))
	v.Use(session.Sessioner())
	v.Pre(m.AddTrailingSlash())
	v.Use(m.CSRFWithConfig(m.CSRFConfig{
		TokenLookup: "header:X-XSRF-TOKEN",
	}))
	v.Use(m.CORSWithConfig(m.CORSConfig{
		AllowOrigins: []string{"https://github.com", "http://yougam.com"},
		AllowHeaders: []string{vodka.HeaderOrigin, vodka.HeaderContentType, vodka.HeaderAcceptEncoding},
	}))
	v.Use(m.Static("static"))
	v.SetRenderer(pongor.Renderor())

	v.File("/favicon.ico", "static/ico/favicon.ico")

	g := v.Group("")
	g.Get("/", handler.MainHandler)

	g.Get("/signup/", handler.SignupGetHandler)
	g.Post("/signup/", handler.SignupPostHandler)

	g.Get("/signin/", handler.SigninGetHandler)
	g.Post("/signin/", handler.SigninPostHandler)

	g.Get("/signout/", handler.SignoutHandler)

	g.Any("/search/", handler.SearchHandler)
	g.Get("/node/:nid/", handler.NodeHandler)
	g.Get("/view/:tid/", handler.ViewHandler)
	g.Get("/category/:cid/", handler.MainHandler)

	// Restricted group
	r := v.Group("")
	r.Use(m.JWTWithConfig(m.JWTConfig{
		SigningKey:  []byte("ZeNpReSe"),
		TokenLookup: "query:token",
	}))
	r.Get("/new/category/", handler.NewCategoryGetHandler)
	r.Post("/new/category/", handler.NewCategoryPostHandler)

	r.Get("/new/node/", handler.NewNodeGetHandler)
	r.Post("/new/node/", handler.NewNodePostHandler)

	r.Get("/new/topic/", handler.NewTopicGetHandler)
	r.Post("/new/topic/", handler.NewTopicPostHandler)

	r.Post("/new/reply/:tid/", handler.NewReplyPostHandler)

	r.Get("/modify/category/", handler.ModifyCatGetHandler)
	r.Post("/modify/category/", handler.ModifyCatPostHandler)

	r.Get("/modify/node/", handler.ModifyNodeGetHandler)
	r.Post("/modify/node/", handler.ModifyNodePostHandler)

	r.Any("/topic/delete/:tid/", handler.TopicDeleteHandler)

	r.Get("/topic/edit/:tid/", handler.TopicEditGetHandler)
	r.Post("/topic/edit/:tid/", handler.TopicEditPostHandler)

	r.Any("/node/delete/:nid/", handler.NodeDeleteHandler)

	r.Get("/node/edit/:nid/", handler.NodeEditGetHandler)
	r.Post("/node/edit/:nid/", handler.NodeEditPostHandler)

	r.Any("/delete/reply/:rid/", handler.DeleteReplyHandler)

	//hotness
	r.Any("/like/:name/:id/", handler.LikeHandler)
	r.Any("/hate/:name/:id/", handler.HateHandler)

	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	vodkapprof.Wrapper(v)
	v.Run(fasthttp.New(":9000"))
}
