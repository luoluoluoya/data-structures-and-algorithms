// Package tree 二叉搜索树：在二叉树的基础上新增节点局部约束即可实现，任意节点x，x.lc.Data <= x.Data <= x.rc.Data
package tree

import (
	"data-structures-and-algorithms/types"
	"fmt"
	"strings"
)

// Bst 二叉查找树：不维护树的渐进平衡性
type Bst struct {
	hot  *BinNode // “命中”节点的父亲
	root *BinNode // 根节点
	size int
}

// NewBst 创建二叉搜索树
func NewBst() *Bst {
	return &Bst{}
}

// Size 树规模
func (t *Bst) Size() int {
	return t.size
}

// Empty 是否为空树
func (t *Bst) Empty() bool {
	return t.size <= 0
}

// Search 二叉树元素查找
func (t *Bst) Search(v types.Sortable) *BinNode {
	return *t.searchAt(&t.root, v)
}

// Insert 二叉树元素插入
func (t *Bst) Insert(v types.Sortable) *BinNode {
	x := t.searchAt(&t.root, v)
	if *x == nil {
		t.size++
		*x = &BinNode{Data: v, parent: t.hot}
	}
	return *x
}

// Remove 二叉树元素删除
func (t *Bst) Remove(v types.Sortable) bool {
	x := t.searchAt(&t.root, v)
	if *x == nil {
		return false
	}
	t.size--
	t.removeAt(x)
	return true
}

// searchAt 在以x为根节点的子树中查找元素v，设置hot指针, 并返回元素所在位置指针（指针的指针，便于上层直接赋值）
func (t *Bst) searchAt(x **BinNode, v types.Sortable) **BinNode {
	t.hot = nil
	if *x != nil {
		t.hot = (*x).parent
	}
	for !equal(*x, v) {
		t.hot = *x
		if v.Less((*x).Data) {
			x = &(*x).lc
		} else {
			x = &(*x).rc
		}
	}
	return x
}

// removeAt 从树t中摘除节点x：返回值指向实际被删除节点的接替者（中序下的直接后继）；hot指向实际被删除节点的父亲
// 算法描述：
// 1、 若当前节点无左子，则将当前节点替换为其右子。（无子结点也被囊括）
// 2、 若当前节点无右子，则将当前节点替换为其左子。
// 3、 双子俱全时：
// 	3.1、 在其右子树中定位其直接后继元素 p。
// 	3.2、 交换 x 与 p。删除 p。
func (t *Bst) removeAt(x **BinNode) *BinNode {
	w := *x           // 实际被删除节点
	var succ *BinNode // 后继节点
	if (*x).lc == nil {
		*x = (*x).rc
		succ = *x
	} else if (*x).rc == nil {
		*x = (*x).lc
		succ = *x
	} else {
		w = (*x).succ()
		w.Data, (*x).Data = (*x).Data, w.Data // todo 此处未作节点交换，只实现对数据项的交换 (外层节点引用的数据信息会出现异常)
		p := w.parent
		if p == (*x) {
			p.rc = w.rc
		} else {
			p.lc = w.rc
		}
		succ = w.rc
	}
	t.hot = w.parent
	if succ != nil { // 设置子节点指针
		succ.parent = t.hot
	}
	return succ
}

// String 中序遍历下输出树元素
func (t *Bst) String() string {
	items := make([]string, 0, t.Size())
	t.root.travelIn(func(v *types.Sortable) {
		items = append(items, fmt.Sprintf("%v", *v))
	})
	return "{" + strings.Join(items, ", ") + "}"
}

// equal 节点判等：外部节点假想为通配符哨兵
func equal(x *BinNode, v types.Sortable) bool {
	return x == nil || x.Data == v
}
