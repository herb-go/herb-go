package middlewares

import (
	"github.com/herb-go/herb/middleware/middlewarefactory"
	"github.com/herb-go/providers/herb/hired/hiredmiddleware"
	"github.com/herb-go/providers/herb/requestpattern/requestpatterncondition"
)

func init() {
	//Register time condition factory.
	middlewarefactory.DefaultContext.RegisterConditionFactory("time", middlewarefactory.NewTimeConditionFactory())
	//Register request pattern condition factory
	middlewarefactory.DefaultContext.RegisterConditionFactory("pattren", requestpatterncondition.NewConditionFactory())
	//Register reponse condition factory.
	middlewarefactory.DefaultContext.RegisterFactory("response", middlewarefactory.NewResponseFactory())
	//Register hiredmiddleware factory.
	middlewarefactory.DefaultContext.RegisterFactory("hiredmiddleware", hiredmiddleware.NewFactory())
}
