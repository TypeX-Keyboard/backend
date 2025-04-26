package bot

import (
	"fmt"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"testing"
)

func TestSBot_Active(t *testing.T) {
	ctx := gctx.New()
	redis := g.Redis()
	var cursor int64 = 0
	keys := make([]string, 0)
	for {
		// SCAN 命令进行遍历
		result, err := redis.Do(ctx, "SCAN", cursor, "MATCH", "ACTIVE:*", "COUNT", 100)
		if err != nil {
			fmt.Println("the query failed:", err)
			return
		}
		// 解析结果
		data := result.Array()
		cursor = gconv.Int64(data[0])
		keys = append(keys, gconv.Strings(data[1])...)
		g.Dump(cursor)
		if cursor == 0 {
			break
		}
	}
	fmt.Println("matched Keys:", keys)
}

func TestSBot_ActiveFriendCounts(t *testing.T) {
	ctx := gctx.New()
	s := New()
	_, err := s.ActiveFriendCounts(ctx)
	if err != nil {
		t.Error(err)
	}
}

func TestSBot_ActiveTask(t *testing.T) {
	ctx := gctx.New()
	s := New()
	err := s.ActiveTask(ctx)
	if err != nil {
		t.Error(err)
	}
}

func TestEveOf1M(t *testing.T) {
	s := New()
	m := s.EveOf1M()
	g.Dump(m)
}

func TestSBot_MiningTask(t *testing.T) {
	ctx := gctx.New()
	s := New()

	err := s.MiningTask(ctx)
	if err != nil {
		t.Error(err)
	}
}

func TestRank(t *testing.T) {
	ctx := gctx.New()
	s := New()
	rank, selfRank, amount, err := s.Rank(ctx, "5XaNdHKpr5dWtPn1e6NrD7Y4tyqvVBvSUtj2hW5rqgzi")
	if err != nil {
		t.Error(err)
	}
	g.Dump(rank)
	g.Dump(selfRank, amount)
}
