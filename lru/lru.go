package lru

import (
	"container/list"
	"time"
)

//Cache 数据全存在这里
type Cache struct {
	//最大节点数量，0为不限制
	maxNode int
	//链表
	llist *list.List
	//hash map
	cache map[string]*list.Element
}

//node 单个数据节点
type node struct {
	key, value string
	expire     int64
}

//NewCache 新建一个Cache
func NewCache(maxNode int) *Cache {
	return &Cache{
		maxNode: maxNode,
		llist:   list.New(),
		cache:   make(map[string]*list.Element),
	}
}

//Add 新增数据
func (c *Cache) Add(key, value string, expire int64) bool {
	//记录过期时间戳
	if expire != 0 {
		expire += time.Now().Unix()
	}

	//数据已经存在，直接覆盖
	if e, ok := c.cache[key]; ok {
		c.llist.MoveToFront(e)
		e.Value.(*node).value = value
		e.Value.(*node).expire = expire
		return true
	}

	//数据不存在，新建node到list
	ne := c.llist.PushFront(&node{key, value, expire})
	c.cache[key] = ne

	//@todo 节点已满，移除老数据or过期数据
	if c.maxNode != 0 && c.llist.Len() > c.maxNode {
		be := c.llist.Back()
		if be != nil {
			c.removeElement(be)
		}
	}

	return true
}

//Get 查询数据
func (c *Cache) Get(key string) (value string, ok bool) {
	if e, ok := c.cache[key]; ok {
		n := e.Value.(*node)
		//惰性检查数据是否过期
		if n.expire != 0 && n.expire < time.Now().Unix() {
			c.removeElement(e)
			return "expire", false
		}
		c.llist.MoveToFront(e)
		return n.value, true
	}

	return "miss", false
}

//Delete 删除数据
func (c *Cache) Delete(key string) bool {
	if e, ok := c.cache[key]; ok {
		c.removeElement(e)
	}
	return true
}

//removeElement 移除指定数据节点
func (c *Cache) removeElement(e *list.Element) {
	c.llist.Remove(e)
	delete(c.cache, e.Value.(*node).key)
}
