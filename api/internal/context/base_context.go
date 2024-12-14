package context

import (
	c "context"
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"go.uber.org/zap"
)

type BaseContext struct {
	Ctx            c.Context
	baseStageList  []Stage
	deferStageList []Stage
	StageErr       error
	Canceller      c.CancelFunc
	ExpireDuration time.Duration
}

type Stage struct {
	handler      StageHandler
	stageMsgChan chan *StageMsg
}

type StageMsg struct {
	needBreak bool
	errMsg    string
	err       error
}

type StageHandler func(ctx c.Context) (needBreak bool, errMsg string, err error)

func (stage *Stage) String() string {
	if reflect.TypeOf(stage.handler).Kind() != reflect.Func {
		panic("input is not a function")
	}
	fullName := runtime.FuncForPC(reflect.ValueOf(stage.handler).Pointer()).Name()
	parts := strings.Split(fullName, ".")
	name := parts[len(parts)-1]
	return strings.TrimRight(name, "-fm")
}

func (ctx *BaseContext) Init(expireDuration time.Duration) {
	ctx.baseStageList = make([]Stage, 0)
	ctx.deferStageList = make([]Stage, 0)
	ctx.ExpireDuration = expireDuration

	if ctx.ExpireDuration > 0 {
		ctx.Ctx, ctx.Canceller = c.WithTimeout(c.Background(), expireDuration)
	} else {
		ctx.Ctx, ctx.Canceller = c.WithCancel(c.Background())
	}
}

func (context *BaseContext) AddBaseHandler(handler StageHandler) (ctx *BaseContext) {
	ctx = context
	context.baseStageList = append(context.baseStageList, Stage{
		handler:      handler,
		stageMsgChan: make(chan *StageMsg, 1),
	})
	return
}

func (context *BaseContext) AddDeferHandler(handler StageHandler) (ctx *BaseContext) {
	ctx = context
	context.deferStageList = append(context.deferStageList, Stage{
		handler:      handler,
		stageMsgChan: make(chan *StageMsg, 1),
	})
	return
}

func (context *BaseContext) Run() {
	defer context.Canceller()

	stageHandlerWrapper := func(stage Stage) (newHandler func(ctx c.Context)) {
		return func(ctx c.Context) {
			defer printPanic()
			defer close(stage.stageMsgChan)
			stageMsg := &StageMsg{needBreak: true}
			defer func() {
				stage.stageMsgChan <- stageMsg
			}()

			stageMsg.needBreak, stageMsg.errMsg, stageMsg.err = stage.handler(ctx)
			if stageMsg.err != nil {
				log.Logger().Error("[BaseContext] error in stage handler", zap.String("stage", stage.String()), zap.Error(stageMsg.err))
				stageMsg.needBreak = true
			}
			return
		}
	}

	for _, stage := range context.deferStageList {
		defer stageHandlerWrapper(stage)(context.Ctx)
	}
	defer printPanic()

	for _, stage := range context.baseStageList {
		go stageHandlerWrapper(stage)(context.Ctx)
		select {
		case <-context.Ctx.Done():
			if context.Ctx.Err() != nil {
				log.Logger().Error("[BaseContext] Ctx.Done closed", zap.Error(context.Ctx.Err()))
				context.StageErr = context.Ctx.Err()
			} else {
				log.Logger().Error("[BaseContext] Ctx.Done closed, no error, get sys message to quit", zap.String("stage", stage.String()))
			}
		case stageMsg, ok := <-stage.stageMsgChan:
			if !ok {
				panic(fmt.Sprintf("[BaseContext] stage.stageMsgChan closed without sending any message, stage: %s", stage.String()))
			}
			context.StageErr = stageMsg.err
			if stageMsg.needBreak {
				log.Logger().Info("[BaseContext] stageMsg.needBreak is true, will cancel task", zap.String("stage", stage.String()))
				return
			}
		}
	}

	return
}

func printPanic() {
	if e := recover(); e != nil {
		log.Logger().Error("======== Panic ========")
		errMsg := fmt.Sprintf("Panic: %v\nTraceBack:\n%s\n", e, string(debug.Stack()))
		log.Logger().Error(errMsg)
		log.Logger().Error("=======================")
	}
}
