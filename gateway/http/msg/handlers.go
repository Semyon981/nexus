package msg

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Semyon981/nexus/gateway/http/auth"
	"github.com/Semyon981/nexus/proto/msgpb"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	msgclient msgpb.ServiceClient
}

func NewHandler(msgclient msgpb.ServiceClient) *Handler {
	return &Handler{
		msgclient: msgclient,
	}
}

type SendMessageInput struct {
	Id_to int64  `json:"id_to" binding:"required"`
	Msg   string `json:"msg" binding:"required"`
}

func (h *Handler) SendMessage(c *gin.Context) {
	inp := new(SendMessageInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	Id_users := c.MustGet(auth.CtxUserKey).(int64)

	_, err := h.msgclient.SendMessage(c.Request.Context(), &msgpb.SendMessageRequest{IdFrom: Id_users, IdTo: inp.Id_to, Msg: inp.Msg, Time: time.Now().Unix()})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetMessages(c *gin.Context) {
	Limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	Offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}

	Id_to, err := strconv.ParseInt(c.Query("id_to"), 10, 64)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

	}
	Id_users := c.MustGet(auth.CtxUserKey).(int64)

	res, err := h.msgclient.GetMessages(c.Request.Context(), &msgpb.GetMessagesRequest{
		Limit:  Limit,
		Offset: Offset,
		IdFrom: Id_users,
		IdTo:   Id_to})
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, res.Messages)
}
