package plugins

import (
	"log"

	"github.com/lightyeario/kelp/api"
	"github.com/lightyeario/kelp/support/utils"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/support/config"
)

type strategyContainer struct {
	description string
	makeFn      func(sdex *SDEX, assetBase *horizon.Asset, assetQuote *horizon.Asset, stratConfigPath string) api.Strategy
}

// strategies is a map of all the strategies available
var strategies = map[string]strategyContainer{
	"buysell": strategyContainer{
		description: "creates buy and sell offers based on a reference price with a pre-specified liquidity depth",
		makeFn: func(sdex *SDEX, assetBase *horizon.Asset, assetQuote *horizon.Asset, stratConfigPath string) api.Strategy {
			var cfg buySellConfig
			err := config.Read(stratConfigPath, &cfg)
			utils.CheckConfigError(cfg, err)
			return makeBuySellStrategy(sdex, assetBase, assetQuote, &cfg)
		},
	},
	"mirror": strategyContainer{
		description: "mirrors an orderbook from another exchange by placing the same orders on Stellar",
		makeFn: func(sdex *SDEX, assetBase *horizon.Asset, assetQuote *horizon.Asset, stratConfigPath string) api.Strategy {
			var cfg mirrorConfig
			err := config.Read(stratConfigPath, &cfg)
			utils.CheckConfigError(cfg, err)
			return makeMirrorStrategy(sdex, assetBase, assetQuote, &cfg)
		},
	},
	"sell": strategyContainer{
		description: "creates sell offers based on a reference price with a pre-specified liquidity depth",
		makeFn: func(sdex *SDEX, assetBase *horizon.Asset, assetQuote *horizon.Asset, stratConfigPath string) api.Strategy {
			var cfg sellConfig
			err := config.Read(stratConfigPath, &cfg)
			utils.CheckConfigError(cfg, err)
			return makeSellStrategy(sdex, assetBase, assetQuote, &cfg)
		},
	},
	"balanced": strategyContainer{
		description: "dynamically prices two tokens based on their relative demand",
		makeFn: func(sdex *SDEX, assetBase *horizon.Asset, assetQuote *horizon.Asset, stratConfigPath string) api.Strategy {
			var cfg balancedConfig
			err := config.Read(stratConfigPath, &cfg)
			utils.CheckConfigError(cfg, err)
			return makeBalancedStrategy(sdex, assetBase, assetQuote, &cfg)
		},
	},
	"delete": strategyContainer{
		description: "deletes all orders for the configured orderbook",
		makeFn: func(sdex *SDEX, assetBase *horizon.Asset, assetQuote *horizon.Asset, stratConfigPath string) api.Strategy {
			return makeDeleteStrategy(sdex, assetBase, assetQuote)
		},
	},
}

// MakeStrategy makes a strategy
func MakeStrategy(
	sdex *SDEX,
	assetBase *horizon.Asset,
	assetQuote *horizon.Asset,
	strategy string,
	stratConfigPath string,
) api.Strategy {
	if strat, ok := strategies[strategy]; ok {
		return strat.makeFn(sdex, assetBase, assetQuote, stratConfigPath)
	}

	log.Fatalf("invalid strategy type: %s\n", strategy)
	return nil
}

// Strategies returns the list of strategies along with the description
func Strategies() map[string]string {
	m := make(map[string]string, len(strategies))
	for s := range strategies {
		m[s] = strategies[s].description
	}
	return m
}

type exchangeContainer struct {
	description string
	makeFn      func() api.Exchange
}

// exchanges is a map of all the exchange integrations available
var exchanges = map[string]exchangeContainer{
	"kraken": exchangeContainer{
		description: "Kraken is a popular centralized cryptocurrency exchange (https://www.kraken.com/)",
		makeFn:      makeKrakenExchange,
	},
}

// MakeExchange is a factory method to make an exchange based on a given type
func MakeExchange(exchangeType string) api.Exchange {
	if exchange, ok := exchanges[exchangeType]; ok {
		return exchange.makeFn()
	}

	log.Fatalf("invalid exchange type: %s\n", exchangeType)
	return nil
}

// Exchanges returns the list of exchanges along with the description
func Exchanges() map[string]string {
	m := make(map[string]string, len(exchanges))
	for name := range exchanges {
		m[name] = exchanges[name].description
	}
	return m
}
