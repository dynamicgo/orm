package orm

import (
	"sync"

	"github.com/dynamicgo/slf4go"

	"github.com/go-xorm/xorm"
)

// PageOrder type
type PageOrder string

// Order types
var (
	ASC  = PageOrder("asc")
	DESC = PageOrder("desc")
)

// Page .
type Page struct {
	Offset  uint64    `binding:"required" json:"offset"`
	Size    uint64    `binding:"required" json:"size"`
	OrderBy string    `json:"orderby"`
	Order   PageOrder `json:"order"`
}

// SyncHandler sync handler prototype
type SyncHandler func() []interface{}

type registerImpl struct {
	slf4go.Logger
	sync.RWMutex
	handlers []SyncHandler
}

func (register *registerImpl) Register(handler SyncHandler) {
	register.Lock()
	defer register.Unlock()

	register.handlers = append(register.handlers, handler)
}

func (register *registerImpl) Sync(engine *xorm.Engine) error {
	register.RLock()
	defer register.RUnlock()

	var tables []interface{}

	for _, handler := range register.handlers {
		tables = append(tables, handler()...)
	}

	return engine.Sync2(tables...)
}

var register = &registerImpl{
	Logger: slf4go.Get("orm"),
}

// Register .
func Register(handler SyncHandler) {
	register.Register(handler)
}

// RegisterWithName .
func RegisterWithName(name string, handler SyncHandler) {
	register.DebugF("register orm module %s", name)
	register.Register(handler)
}

// Sync .
func Sync(engine *xorm.Engine) error {
	return register.Sync(engine)
}
