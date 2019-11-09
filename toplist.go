package main

import (
    "sync"
)

type TopList struct {
    list []ListElem
    lock *sync.Mutex
}

func NewTopList() *TopList {
    return &TopList {
        make([]ListElem, 0),
        &sync.Mutex{},
    }
}

type ListElem struct {
    Elem interface{}
    Score int
    Perms []int
}

func (tl *TopList) Add(le ListElem) {
    tl.lock.Lock()
    defer tl.lock.Unlock()

    idx := 0
    for i := 0; i < len(tl.list); i++ {
        if i == len(tl.list) - 1 || tl.list[i].Score == 0 {
            idx = i
            break
        }
        if le.Score > tl.list[i].Score {
            idx = i
            break
        }
    }
    
    newList := make([]ListElem, 0)
    if len(tl.list) > 0 {
        newList = append(newList, tl.list[0:idx]...)
    }
    newList = append(newList, le)
    if len(tl.list) > 0 {
        newList = append(newList, tl.list[idx:len(tl.list)]...)
    }

    tl.list = takeTop(newList, 10)
}

func takeTop(list []ListElem, count int) []ListElem {
    if count > len(list) {
        count = len(list)
    }
    return list[0:count]
}
