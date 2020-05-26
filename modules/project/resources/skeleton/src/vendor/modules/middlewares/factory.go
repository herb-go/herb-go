package middlewares

import (
	"github.com/herb-go/herb/middleware/middlewarefactory"
	"github.com/herb-go/providers/hired/hiredmiddleware"
)

func init() {
	//Register time condition factory.
	middlewarefactory.DefaultContext.RegisterConditionFactory("time", middlewarefactory.NewTimeConditionFactory())
	//Register reponse condition factory.
	middlewarefactory.DefaultContext.RegisterFactory("response", middlewarefactory.NewResponseFactory())
	//Register hiredmiddleware factory.
	middlewarefactory.DefaultContext.RegisterFactory("hiredmiddleware", hiredmiddleware.NewFactory())
}
