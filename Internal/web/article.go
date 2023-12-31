package web

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/wureny/webook/webook/Internal/domain"
	"github.com/wureny/webook/webook/Internal/service"
	ijwt "github.com/wureny/webook/webook/Internal/web/jwt"
	"github.com/wureny/webook/webook/pkg/ginx"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	svc service.ArticleService
	inte service.InteractiveService
	biz string
}

func NewArticleHandler(svc service.ArticleService,inte service.InteractiveService) handler {
	return &ArticleHandler{
		svc: svc
		inte: inte
	}
}

func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)
	// 创作者的查询接口
	// 这个是获取数据的接口，理论上来说（遵循 RESTful 规范），应该是用 GET 方法
	// GET localhost/articles => List 接口
	g.POST("/list",
		ginx.WrapBodyAndToken[ListReq, ijwt.UserClaims](a.List))
	g.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](a.Detail))
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//	a.l.Error("未发现用户的 session 信息")
		return
	}
	// 检测输入，跳过这一步
	// 调用 svc 的代码
	id, err := a.svc.Save(ctx, req.toDomain(int64(claims.Id)))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？
		//		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		//TODO
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		//	a.l.Error("未发现用户的 session 信息")
		return
	}
	id, err := a.svc.Publish(ctx, req.toDomain(int64(claims.Id)))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？
		//	h.l.Error("发表帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 检测输入，跳过这一步
	// 调用 svc 的代码
	err := h.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Id,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})

}

func (h *ArticleHandler) List(ctx *gin.Context, req ListReq, uc ijwt.UserClaims) (ginx.Result, error) {
	res, err := h.svc.List(ctx, uc.Id, req.Offset, req.Limit)
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, nil
	}
	// 在列表页，不显示全文，只显示一个"摘要"
	// 比如说，简单的摘要就是前几句话
	// 强大的摘要是 AI 帮你生成的
	return ginx.Result{
		Data: slice.Map[domain.Article, ArticleVO](res,
			func(idx int, src domain.Article) ArticleVO {
				return ArticleVO{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstract(),
					Status:   src.Status.ToUint8(),
					// 这个列表请求，不需要返回内容
					//Content: src.Content,
					// 这个是创作者看自己的文章列表，也不需要这个字段
					//Author: src.Author
					Ctime: src.Ctime.Format(time.DateTime),
					Utime: src.Utime.Format(time.DateTime),
				}
			}),
	}, nil
}

func (a *ArticleHandler) Detail(ctx *gin.Context, usr ijwt.UserClaims) (ginx.Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		//ctx.JSON(http.StatusOK, )
		//a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return ginx.Result{
			Code: 4,
			Msg:  "参数错误",
		}, err
	}
	art, err := a.svc.GetById(ctx, id)
	if err != nil {
		//ctx.JSON(http.StatusOK, )
		//a.l.Error("获得文章信息失败", logger.Error(err))
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	// 这是不借助数据库查询来判定的方法
	if art.Author.Id != usr.Id {
		//ctx.JSON(http.StatusOK)
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		//a.l.Error("非法访问文章，创作者 ID 不匹配",
		//	logger.Int64("uid", usr.Id))
		return ginx.Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		}, fmt.Errorf("非法访问文章，创作者 ID 不匹配 %d", usr.Id)
	}
	return ginx.Result{
		Data: ArticleVO{
			Id:    art.Id,
			Title: art.Title,
			// 不需要这个摘要信息
			//Abstract: art.Abstract(),
			Status:  art.Status.ToUint8(),
			Content: art.Content,
			// 这个是创作者看自己的文章列表，也不需要这个字段
			//Author: art.Author
			Ctime: art.Ctime.Format(time.DateTime),
			Utime: art.Utime.Format(time.DateTime),
		},
	}, nil
}
