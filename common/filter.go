package common

import "net/http"

//声明一个拦截器的函数类型
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

//拦截器结构体
type Filter struct {
	//用来存储需要拦截的URI
	filterMap map[string]FilterHandle
}

//初始化Filter
func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

//注册路由和对应的过滤函数到filterMap
func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

//根据uri获取对应的handler
func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

//声明一个新的函数类型，用于处理正常通过过滤器的请求
type WebHandle func(rw http.ResponseWriter, req *http.Request)

//执行操作
func (f *Filter) Handle(webHandle WebHandle) WebHandle {
	return func(rw http.ResponseWriter, req *http.Request) {
		//执行拦截器的拦截操作
		for path, handler := range f.filterMap {
			//如果当前访问的uri和拦截器中注册的uri一致的话
			if path == req.RequestURI {
				//对当前的访问进行对应过滤
				err := handler(rw, req)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		webHandle(rw, req)
	}
}
