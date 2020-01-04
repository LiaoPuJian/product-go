package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"product-go/common"
	"product-go/encrypt"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

//设置好validate的集群
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

var hashConsistent *common.Consistent

//这个结构体用来存放控制信息
type AccessControl struct {
	sourcesArray map[int]interface{}
	*sync.RWMutex
}

var accessControl = &AccessControl{sourcesArray: make(map[int]interface{})}

//根据uid获取指定的数据
func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.sourcesArray[uid]
}

//设置uid对应的数据
func (m *AccessControl) SetNewRecord(uid int, data interface{}) {
	m.Lock()
	defer m.Unlock()
	m.sourcesArray[uid] = data
}

func (m *AccessControl) GetDistributedRight(r *http.Request) bool {
	uid, err := r.Cookie("uid")
	if err != nil {
		return false
	}

	//采用一致性hash算法，获取uid的hash值，并对应到对应的机器节点
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	if hostRequest == localHost {
		return m.GetDataFromMap(uid.Value)
	} else {
		return GetDataFromOtherMap(hostRequest, r)
	}
}

//从本机的数据map中获取数据
func (m *AccessControl) GetDataFromMap(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)
	if data != nil {
		return true
	}
	return false
}

//从其他节点的数据map中回去数据
func GetDataFromOtherMap(host string, r *http.Request) bool {
	//获取uid
	uid, err := r.Cookie("uid")
	if err != nil {
		return false
	}
	//获取sign
	sign, err := r.Cookie("sign")
	if err != nil {
		return false
	}
	//模拟接口访问
	client := &http.Client{}
	url := fmt.Sprintf("http://%s:%s/check", host, port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	//指定请求的cookie（将当前的cookie写入到这次请求中）
	cookieUid := &http.Cookie{Name: "uid", Value: uid.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: sign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	//执行请求
	response, err := client.Do((req))
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("执行/check正常的业务逻辑")
}

func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("开始执行/check下的验证")
	//添加基于cookie的权限验证
	return CheckUserInfo(r)
}

func CheckUserInfo(r *http.Request) error {
	//获取cookie中存储的uid
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return err
	}
	//获取用户的sign
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return err
	}
	//对sign执行解密，判断其和uid是否相等
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return err
	}
	fmt.Println("结果对比：")
	fmt.Println("用户ID:", uidCookie.Value)
	fmt.Println("解密后的用户ID:", string(signByte))

	if string(signByte) == uidCookie.Value {
		return nil
	} else {
		return errors.New("身份校验失败，不匹配")
	}
}

func main() {

	//分布式权限验证的负载均衡（采用一致性hash算法）
	hashConsistent = common.NewConsistent()

	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	//1、获取一个过滤器
	filter := common.NewFilter()
	//注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	//设置过滤器的路由
	http.HandleFunc("/check", filter.Handle(Check))
	//启动服务
	http.ListenAndServe(":8083", nil)
}
