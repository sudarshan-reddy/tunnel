package main

import (
	"container/list"
	"fmt"
	"time"
)

type tunnelData chan []byte

type Tunnel struct {
	*list.List
	closed chan struct{}
}

func New() *Tunnel {
	t := &Tunnel{List: list.New(), closed: make(chan struct{})}
	t.Init()
	return t
}

func (t *Tunnel) Put(v chan []byte) {
	t.PushBack(v)
}

func (t *Tunnel) Get() chan []byte {
	dCh := make(chan []byte)
	go func() {
		defer close(dCh)
		for {
			select {
			case <-t.closed:
				return
			default:
				for e := t.Front(); e != nil; e = e.Next() {
					dataCh := e.Value.(chan []byte)
					data, ok := <-dataCh
					if !ok {
						continue
					}
					dCh <- data
				}
			}
		}
	}()
	return dCh
}

func (t *Tunnel) Close() {
	close(t.closed)
}

func main() {

	list := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	//var semaphore = make(chan struct{}, 5)
	t := New()

	go func() {
		for _, each := range list {
			//	semaphore <- struct{}{}
			var iCh = make(chan []byte)
			t.Put(iCh)

			go func(i int, byteCh chan []byte) {
				defer close(byteCh)
				//defer func() {
				//	<-semaphore
				//}()
				fmt.Println("--->", i)
				byteCh <- []byte(string(i))
			}(each, iCh)
		}
	}()

	go func() {
		for stuff := range t.Get() {
			fmt.Println(stuff)
		}
	}()

	time.Sleep(3 * time.Second)
	t.Close()

}
