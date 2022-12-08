/**
 * @Author: Jason
 * @Description:
 * @File: router
 * @Version: 1.0.0
 * @Date: 2022/4/7 11:20
 */

package easysocket

import "google.golang.org/protobuf/proto"

type IRouter interface {
	Handle(request IRequest, message proto.Message)
}
