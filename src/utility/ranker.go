package utility

import (
	"fmt"
	"sort"
	"sync"
)

type TRankItem struct {
	RankID    int32
	RankValue int
}

type TRankList []*TRankItem

type TRanker struct {
	sync.Mutex
	List    TRankList
	ShowNum int //显示数量
	RankNum int //排行数量
}

func (list TRankList) Len() int {
	return len(list)
}

func (list TRankList) Less(i, j int) bool {
	if list[i].RankValue > list[j].RankValue {
		return true
	}

	return false
}

func (list *TRankList) Swap(i, j int) {
	var pItem *TRankItem = (*list)[i]
	(*list)[i] = (*list)[j]
	(*list)[j] = pItem
}

//show 显示的条数
//rank 总的排名条数
func (ranker *TRanker) InitRanker(show int, rank int) {
	ranker.ShowNum = show
	ranker.RankNum = rank
	ranker.List = make([]*TRankItem, rank)
	for i := 0; i < len(ranker.List); i++ {
		ranker.List[i] = &TRankItem{0, -1}
	}
}

//清空排行榜
func (ranker *TRanker) Clear() {
	for i := 0; i < len(ranker.List); i++ {
		ranker.List[i].RankID = 0
		ranker.List[i].RankValue = -1
	}
}

func (ranker *TRanker) SetRankItem(rankid int32, rankvalue int) int {
	ranker.Lock()
	defer ranker.Unlock()
	nCount := len(ranker.List)
	MinValue := ranker.List[nCount-1].RankValue
	if rankvalue <= MinValue {
		return -1
	}

	targetIndex := sort.Search(nCount, func(i int) bool {
		if ranker.List[i].RankValue <= rankvalue {
			return true
		}
		return false
	})

	myIndex := nCount - 1
	for i := targetIndex; i < nCount; i++ {
		if ranker.List[i].RankID == rankid || ranker.List[i].RankID == 0 {
			ranker.List[i].RankValue = rankvalue
			ranker.List[i].RankID = rankid
			myIndex = i
			break
		}
	}

	if myIndex == targetIndex {
		return targetIndex
	}

	for i := myIndex; i > targetIndex; i-- {
		ranker.List[i].RankID = ranker.List[i-1].RankID
		ranker.List[i].RankValue = ranker.List[i-1].RankValue
	}

	ranker.List[targetIndex].RankID = rankid
	ranker.List[targetIndex].RankValue = rankvalue
	return targetIndex
}

//战力排行榜专用, 其它的排行榜暂时不要用
func (ranker *TRanker) SetRankItemEx(rankid int32, orgvalue int, newvalue int) int {
	ranker.Lock()
	defer ranker.Unlock()
	nCount := len(ranker.List)
	MinValue := ranker.List[nCount-1].RankValue

	myIndex := -1
	if orgvalue >= MinValue {
		for i := 0; i < nCount; i++ {
			if rankid == ranker.List[i].RankID {
				myIndex = i
				break
			}
		}
		if myIndex > 0 {
			ranker.List[myIndex].RankValue = newvalue
			sort.Sort(&ranker.List)
		}

	} else {
		if newvalue > MinValue {
			ranker.List[nCount-1].RankValue = newvalue
			sort.Sort(&ranker.List)
		}
	}

	if newvalue > MinValue {
		return ranker.GetRankIndex(rankid, newvalue)
	}

	return -1
}

func (ranker *TRanker) GetRankIndex(rankid int32, rankvalue int) int {
	nCount := len(ranker.List)
	MinValue := ranker.List[nCount-1].RankValue
	if rankvalue <= MinValue {
		return -1
	}

	targetIndex := sort.Search(nCount,
		func(i int) bool {
			if ranker.List[i].RankValue <= rankvalue {
				return true
			}
			return false
		})

	if targetIndex == nCount {
		return -1
	}

	for i := targetIndex; i >= 0; i-- {
		if ranker.List[i].RankID == rankid {
			return i + 1
		}
	}

	return -1
}

func (ranker *TRanker) ForeachShow(handler func(int32, int)) {
	sum := 0
	for _, v := range ranker.List {
		if sum >= ranker.ShowNum || v.RankID == 0 {
			break
		}
		sum++
		handler(v.RankID, v.RankValue)
	}
}

func (ranker *TRanker) CopyFrom(src *TRanker) {
	if src == nil || src.List == nil {
		return
	}
	for i := 0; i < len(src.List); i++ {
		if i >= ranker.RankNum || src.List[i].RankID <= 0 {
			break
		}
		ranker.List[i].RankID = src.List[i].RankID
		ranker.List[i].RankValue = src.List[i].RankValue
	}

	return
}

func (ranker *TRanker) Print() {
	for i := 0; i < len(ranker.List); i++ {
		fmt.Println(ranker.List[i])
	}
}
