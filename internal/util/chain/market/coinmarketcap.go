package market

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/net/gclient"
	"math/rand"
	"strings"
)

const (
	REST = "https://pro-api.coinmarketcap.com"
)

var ApiKeys = []string{"", ""}

type sMarket struct {
	cli *gclient.Client
}

type Imarket interface {
	Quotes(context.Context, []string, ...int) ([]Slug, error)
	TokenList(context.Context, int, int) ([]TokenMap, error)
}

func New() Imarket {
	cli := gclient.New().SetPrefix(REST)
	// 从数组中随机选择一个值
	randomIndex := rand.Intn(len(ApiKeys))
	randomValue := ApiKeys[randomIndex]
	cli.SetHeader("X-CMC_PRO_API_KEY", randomValue)
	return &sMarket{
		cli,
	}
}

func (s *sMarket) Quotes(ctx context.Context, ids []string, convertIds ...int) ([]Slug, error) {
	convertId := 825
	if len(convertIds) > 0 {
		convertId = convertIds[0]
	}
	r, err := s.cli.Get(ctx, fmt.Sprintf("/v2/cryptocurrency/quotes/latest?id=%s&convert_id=%d", strings.Join(ids, ","), convertId))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	body := r.ReadAll()
	var res *MarketRes
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	slugs := make([]Slug, 0)
	for _, slug := range res.Data {
		slugs = append(slugs, slug)
	}
	return slugs, nil
}

func (s *sMarket) TokenList(ctx context.Context, start, limit int) ([]TokenMap, error) {
	r, err := s.cli.Get(ctx, fmt.Sprintf("/v1/cryptocurrency/map?listing_status=%s&limit=%d&sort=%s&start=%d", "active", limit, "cmc_rank", start))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	body := r.ReadAll()
	var res *MapRes
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}
