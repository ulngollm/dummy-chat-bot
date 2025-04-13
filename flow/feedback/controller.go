package feedback

import (
	"fmt"

	"github.com/ulngollm/teleflow"
	tele "gopkg.in/telebot.v4"
)

const FlowName string = "default"

type FlowController struct {
	manager *teleflow.FlowManager
}

type eventData map[string]string

func New(manager *teleflow.FlowManager) *FlowController {
	return &FlowController{manager: manager}
}

func (r *FlowController) Init(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		id := c.Sender().ID
		flow := teleflow.NewSimpleFlow(id, StateUndefined, FlowName)
		if err := r.manager.InitFlow(flow); err != nil {
			return err
		}
		teleflow.SaveToCtx(c, flow)
		return next(c)
	}
}

func (r *FlowController) AskCategory(nextState string) tele.HandlerFunc {
	return func(c tele.Context) error {
		inner := func(context tele.Context) error {
			return c.Send("назовите категорию")
		}
		return r.buildHandler(c, inner, nextState)(c)
	}
}

func (r *FlowController) AskProduct(c tele.Context) error {
	inner := func(context tele.Context) error {
		return c.Send("назовите продукт")
	}
	return r.buildHandler(c, inner, StateAskedProduct)(c)
}

func (r *FlowController) AskDetails(c tele.Context) error {
	inner := func(context tele.Context) error {
		return c.Send("опишите, как воспроизвести ошибку")
	}
	return r.buildHandler(c, inner, StateAskedDetails)(c)
}

func (r *FlowController) AskScreenshot(c tele.Context) error {
	inner := func(context tele.Context) error {
		return c.Send("приложите скриншот")
	}
	return r.buildHandler(c, inner, StateAskedScreenshot)(c)
}

func (r *FlowController) Thank(c tele.Context) error {
	inner := func(context tele.Context) error {
		return c.Send("спасибо за обратную связь! мы передали ваше сообщение в поддержку")
	}
	return r.buildHandler(c, inner, StateComplete)(c)
}

func (r *FlowController) buildHandler(
	c tele.Context,
	inner tele.HandlerFunc,
	stateTo string,
) tele.HandlerFunc {
	return func(context tele.Context) error {
		flow := teleflow.GetCurrentFlow(c)
		var d eventData
		if flow.Data() == nil {
			d = make(eventData)
		} else {
			d = flow.Data().(eventData)
		}
		state := flow.State()
		d[state] = c.Message().Text
		flow.SetData(d)

		if err := inner(c); err != nil {
			return fmt.Errorf("inner: %v", err)
		}
		flow.SetState(stateTo)
		return nil
	}
}
