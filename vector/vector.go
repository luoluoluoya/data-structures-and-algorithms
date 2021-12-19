// Package vector 向量
package vector

import (
	"data-structures-and-algorithms/types"
	"fmt"
	"math/rand"
	"strings"
)

const lowCapacity = 3

type Vector struct {
	elem []types.Sortable
}

// New 直接实例化
func New(cap int) *Vector {
	return &Vector{
		elem: make([]types.Sortable, 0, cap),
	}
}

// Copy 复制另一向量
func Copy(o *Vector) *Vector {
	vec := New(o.Capacity())
	vec.elem = vec.elem[:o.Size()]
	copy(vec.elem, o.elem)
	return vec
}

// CopySlice 复制切片以创建向量
func CopySlice(o ...types.Sortable) *Vector {
	vec := New(cap(o))
	vec.elem = vec.elem[:len(o)]
	copy(vec.elem, o)
	return vec
}

// Less 返回vec[i]与vec[j]的大小
func (vec *Vector) Less(i, j int) bool {
	return vec.elem[i].Less(vec.elem[j])
}

// Swap 交换i，j 元素位置
func (vec *Vector) Swap(i, j int) {
	vec.elem[i], vec.elem[j] = vec.elem[j], vec.elem[i]
}

// Size 返回向量的大小
func (vec *Vector) Size() int {
	return len(vec.elem)
}

// Empty 判定向量是否为空
func (vec *Vector) Empty() bool {
	return vec.Size() <= 0
}

// Capacity 返回向量的容量
func (vec *Vector) Capacity() int {
	return cap(vec.elem)
}

// Front 返回向量中第一个元素。
func (vec *Vector) Front() (types.Sortable, error) {
	if vec.Empty() {
		return nil, fmt.Errorf("Vector is empty")
	}
	return vec.elem[0], nil
}

// Back 返回向量中最后一个元素。
func (vec *Vector) Back() (types.Sortable, error) {
	if vec.Empty() {
		return nil, fmt.Errorf("Vector is empty")
	}
	return vec.elem[vec.Size()-1], nil
}

// At 返回向量中位置 n 处元素。
func (vec *Vector) At(r int) (types.Sortable, error) {
	if !vec.validRank(r) {
		return nil, fmt.Errorf("access out of bounds. len = %v, idx = %v", vec.Size(), r)
	}
	return vec.elem[r], nil
}

// Assign 为向量分配新内容，替换其当前内容。 当 n 非法时返回false。
func (vec *Vector) Assign(r int, v types.Sortable) bool {
	if !vec.validRank(r) {
		return false
	}
	vec.elem[r] = v
	return true
}

// Disordered 返回向量的逆序对数
func (vec *Vector) Disordered() int {
	var num = 0
	for i := 1; i < vec.Size(); i++ {
		if vec.elem[i].Less(vec.elem[i-1]) {
			num++
		}
	}
	return num
}

// String 字符串形式
func (vec *Vector) String() string {
	if vec.Empty() {
		return "{}"
	}
	items := make([]string, 0, vec.Size())
	for _, item := range vec.elem {
		items = append(items, fmt.Sprintf("%v", item))
	}
	return "{" + strings.Join(items, ", ") + "}"
}

// insert 在 [0,size] 的指定位置处进行插入。 当 n 非法时返回false。
func (vec *Vector) insert(r int, v types.Sortable) bool {
	size := vec.Size()
	if r != size && !vec.validRank(r) {
		return false
	}
	vec.expand()
	vec.elem = vec.elem[:size+1]
	for i := size; i > r; i-- {
		vec.elem[i] = vec.elem[i-1]
	}
	vec.elem[r] = v
	return true
}

// Remove 移除向量中秩为 r 的元素
func (vec *Vector) Remove(r int) (types.Sortable, error) {
	if !vec.validRank(r) {
		return nil, fmt.Errorf("out of bounds. len = %v, idx = %v", vec.Size(), r)
	}
	e := vec.elem[r]
	vec.RemoveRange(r, r+1)
	return e, nil
}

// RemoveRange 移除秩在区间 [lo,hi) 中的元素
func (vec *Vector) RemoveRange(lo, hi int) {
	if lo >= hi || !vec.validRank(lo) {
		return
	}
	vec.elem = append(vec.elem[:lo], vec.elem[hi:]...)
	vec.shrink()
}

// Clear 清空向量，不收缩所占空间
func (vec *Vector) Clear() {
	vec.elem = vec.elem[:0]
}

// Push 尾部进行插入
func (vec *Vector) Push(v types.Sortable) {
	vec.insert(vec.Size(), v)
}

// Pop 尾部进行删除
func (vec *Vector) Pop() (types.Sortable, error) {
	return vec.Remove(vec.Size() - 1)
}

// Scrambling 向量整体置乱
func (vec *Vector) Scrambling() {
	vec.ScramblingRange(0, vec.Size())
}

// ScramblingRange 向量区间[lo, hi)置乱
func (vec *Vector) ScramblingRange(lo, hi int) {
	for ; lo < hi; hi-- {
		vec.Swap(rand.Intn(hi), hi-1)
	}
}

// Traverse 遍历向量元素
func (vec *Vector) Traverse(visit func(*types.Sortable)) {
	for i := 0; i < vec.Size(); i++ {
		visit(&vec.elem[i])
	}
}

// Find 无序向量查找：多个元素时返回秩最大者，失败时返回-1
func (vec *Vector) Find(v types.Sortable) int {
	return vec.FindRange(v, 0, vec.Size())
}

// FindRange 无序向量区间 [lo, hi) 查找：失败时返回-1
func (vec *Vector) FindRange(v types.Sortable, lo, hi int) int {
	for ; lo < hi && vec.elem[hi-1] != v; hi-- {
	}
	if lo == hi {
		return -1
	}
	return hi - 1
}

// Deduplicate 无序向量去重
func (vec *Vector) Deduplicate() int {
	oldSize := vec.Size()
	set := make(map[types.Sortable]struct{})
	for i := 0; i < vec.Size(); i++ {
		if _, ok := set[vec.elem[i]]; ok {
			_, _ = vec.Remove(i)
		} else {
			set[vec.elem[i]] = struct{}{}
		}
	}
	return oldSize - vec.Size()
}

// Search 有序向量整体查找, 返回不大于v的元素的最大秩
func (vec *Vector) Search(v types.Sortable) int {
	return vec.SearchRange(v, 0, vec.Size())
}

// SearchRange 有序向量区间 [lo, hi) 查找, 返回不大于v的元素的最大秩
func (vec *Vector) SearchRange(v types.Sortable, lo, hi int) int {
	return vec.binarySearchV2(v, lo, hi)
}

// Uniquify 有序向量去重
func (vec *Vector) Uniquify() int {
	var i, j = 0, 1
	for ; j < vec.Size(); j++ {
		if vec.elem[i] != vec.elem[j] {
			i++
			vec.elem[i] = vec.elem[j]
		}
	}
	vec.elem = vec.elem[:i+1]
	vec.shrink()
	return j - vec.Size()
}

// 朴素二分查找
func (vec *Vector) binarySearchV1(v types.Sortable, lo, hi int) int {
	for lo < hi {
		mid := (lo + hi) >> 1
		if v.Less(vec.elem[mid]) {
			hi = mid
		} else if vec.elem[mid].Less(v) {
			lo = mid + 1
		} else {
			return mid
		}
	}
	return lo
}

// 优化二分查找
func (vec *Vector) binarySearchV2(v types.Sortable, lo, hi int) int {
	for lo < hi {
		mid := (lo + hi) >> 1
		if v.Less(vec.elem[mid]) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo - 1
}

// 扩容：空间不足时对容量执行翻倍
func (vec *Vector) expand() {
	if vec.Size() < vec.Capacity() {
		return
	}
	cap := vec.Capacity()
	if cap < lowCapacity {
		cap = lowCapacity
	} else {
		cap <<= 1
	}
	newElem := make([]types.Sortable, vec.Size(), cap)
	copy(newElem, vec.elem)
	vec.elem = newElem
}

// 缩容：维持空间利用率在 50% 之上
func (vec *Vector) shrink() {
	if vec.Size()<<2 > vec.Capacity() || vec.Capacity() < lowCapacity<<1 {
		return
	}
	cap := vec.Capacity() >> 1
	newElem := make([]types.Sortable, vec.Size(), cap)
	copy(newElem, vec.elem)
	vec.elem = newElem
}

// 验证秩是否合法
func (vec *Vector) validRank(r int) bool {
	return 0 <= r && r < vec.Size()
}
