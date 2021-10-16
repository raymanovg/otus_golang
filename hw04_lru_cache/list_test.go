package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		// deleting middle element of the list
		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		front := l.Front()
		require.Equal(t, 10, front.Value)
		require.Nil(t, front.Prev)
		require.NotNil(t, front.Next)
		require.Equal(t, 30, front.Next.Value)

		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)

		// deleting front and back element from list
		front = l.Front()
		l.Remove(front)
		require.Equal(t, 6, l.Len())
		front = l.Front()
		require.Nil(t, front.Prev)
		require.NotNil(t, front.Next)

		back := l.Back()
		l.Remove(back)
		require.Equal(t, 5, l.Len())
		back = l.Back()
		require.Nil(t, back.Next)
		require.NotNil(t, back.Prev)
	})
}
