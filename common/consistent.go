package common

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

//这里是一致性hash算法的实现类，需要实现一个hash的圆环，用于将不同节点的validate服务
//注册到这个hash上，在前端的请求通过slb均匀的发送的不同的validate上时，可以通过这个
//圆环算法，找到对应的数据存储在圆环中哪个validate的机器上，并对其进行访问
//如果数据存储在自己上，则直接访问，否则重定向请求到目标机器上获取数据

//声明新的切片类型  假定圆环的长度为32位无符号整形，设定以下三个方法用于排序
type units []uint32

//返回切片长度
func (x units) Len() int {
	return len(x)
}

//比对两个数大小
func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

//切片中两个值的交换
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var errEmpty = errors.New("Hash环没有数据！")

//创建结构体用于保存一致性hash信息
type Consistent struct {
	//hash 环，key为hash值，string为节点的信息，可以为节点的ip
	circle map[uint32]string
	//已经排序好的节点hash切片
	sortedHashes units
	//虚拟节点个数  用于增加hash的平衡性
	//这里解释一下，如果实体节点数量太少，可能出现数据大量存储在某个节点上，另外的节点只有少量的数据
	//这个虚拟节点即为解决次问题
	VirtualNode int
	//map的读写锁
	sync.RWMutex
}

//构造方法
func NewConsistent() *Consistent {
	return &Consistent{
		//初始map
		circle: make(map[uint32]string),
		//设置虚拟节点的数量  先默认给个20
		VirtualNode: 20,
	}
}

//自动生成key值
//element代表实体节点信息，index代表生成虚拟节点后面需要的拼接值
func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

//获取hash的位置
func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		var srcatch [64]byte
		copy(srcatch[:], key)
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	//使用IEEE多项式返回数据的CRC-32校验和
	return crc32.ChecksumIEEE([]byte(key))
}

//更新排序，方便查找
func (c *Consistent) updateSortedHashes() {
	//初始化排序数组
	hashes := c.sortedHashes[:0]
	//判断切片容量是否过大，如果过大则充值
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}
	//添加hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	//对所有的节点进行hash值的排序
	sort.Sort(hashes)
	c.sortedHashes = hashes
}

//添加节点
func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

//添加节点
func (c *Consistent) add(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		c.circle[c.hashKey(c.generateKey(element, i))] = element
	}
	c.updateSortedHashes()
}

//删除节点
func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

//删除节点
func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

//顺时针查找最近的节点
func (c *Consistent) search(key uint32) int {
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//使用二分查找算法来搜索指定切片满足条件的最小值
	i := sort.Search(len(c.sortedHashes), f)
	//如果超出范围则设置i=0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

//根据数据标识获取最近的服务器节点信息
//这里传入的name为服务器的节点名称，可以是ip
func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	//根据服务器的ip信息获取hash值
	key := c.hashKey(name)
	//根据这个uint32的key获取其在circle中的最近一个值
	i := c.search(key)
	//得到这个值的键，并返回
	return c.circle[c.sortedHashes[i]], nil
}
