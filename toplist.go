package main

type TopList struct {
    list []ListElem
}

func NewTopList() *TopList {
    return &TopList {
        make([]ListElem, 0),
    }
}

type ListElem struct {
    Elem interface{}
    Score int
}

func (tl *TopList) Add(le ListElem) {
    idx := 0
    for i := 0; i < len(tl.list); i++ {
        if le.Score < tl.list[i].Score {
            idx = i
            break
        }
    }
    
    newList := make([]ListElem, 0)
    newList = append(newList, tl.list[0:idx]...)
    newList = append(newList, le)
    if len(tl.list) > 1 {
        newList = append(newList, tl.list[idx:len(tl.list)-1]...)
    }

    tl.list = takeTop(newList, 100)
}

func takeTop(list []ListElem, count int) []ListElem {
    if count > len(list) {
        count = len(list)
    }
    return list[0:count]
}
