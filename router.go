/**
 * @Author: Jason
 * @Description:
 * @File: router
 * @Version: 1.0.0
 * @Date: 2022/4/7 11:20
 */

package easysocket

type IRouter interface {
	PreHandle(request IRequest)
	Handle(request IRequest)
	PostHandle(request IRequest)
}

type BaseRouter struct{}

func (r *BaseRouter) PreHandle(request IRequest) {}

func (r *BaseRouter) Handle(request IRequest) {}

func (r *BaseRouter) PostHandle(request IRequest) {}
