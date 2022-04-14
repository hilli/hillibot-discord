package finance

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/equity"
	wc "github.com/wcharczuk/go-chart/v2"
)

// Chops up a string of symbols into a slice of strings
func SymbolsToList(symbols string) []string {
	symbolList := strings.Split(symbols, ",")
	for i, symbol := range symbolList {
		symbolList[i] = strings.ToUpper(strings.TrimSpace(symbol))
	}
	return symbolList
}

// Returns a rendered graph (PNG) in a bytes.Buffer pointer
func RenderedGraphForSymbols(symbols []string, buffer *bytes.Buffer) {
	graph := graphForSymbols(symbols)
	graph.Elements = []wc.Renderable{
		wc.Legend(&graph),
	}
	graph.Render(wc.PNG, buffer)
	// Debug write to file as well
	f, err := os.OpenFile("finance.png", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
			panic(err)
	}
	defer f.Close()
	graph.Render(wc.PNG, f)
}

func getSymbolInfo(symbolList []string) string {
	var result string
	if len(symbolList) == 0 || symbolList[0] == "" {
		result = "Please name SYMBOLS to your preferred stock ticker names, comma seperated. Ie AAPL,TSLA"
		return result
	}
	iter := equity.List(symbolList)
	// Iterate over results. Will exit upon any error.
	for iter.Next() {
		q := iter.Equity()
		result = result + fmt.Sprintf("%s (%s): Bid: %.2f Ask: %.2f Price: %.2f High: %.2f Low: %.2f Close: %.2f Post: %.2f Currency: %s Market State: %s\n",
			q.Symbol,
			q.ShortName,
			q.Bid,
			q.Ask,
			q.RegularMarketPrice,
			q.RegularMarketDayHigh,
			q.RegularMarketDayLow,
			q.RegularMarketPreviousClose,
			q.RegularMarketPrice+q.PreMarketChange,
			q.CurrencyID,
			q.MarketState)
	}
	log.Info().Msg(result)
	return result
}

func graphForSymbols(symbols []string) wc.Chart {
	var chartSeries []wc.Series
	for _, symbol := range symbols {
		chartSeries = append(chartSeries, getTimeSeriesForSymbol(symbol))
	}
	graph := wc.Chart{
		Background: wc.Style{
			Padding: wc.Box{
				Top:  20,
				Left: 20,
			},
		},
		Series: chartSeries,
	}

	graph.Elements = []wc.Renderable{
		wc.Legend(&graph),
	}
	return graph
}

func getTimeSeriesForSymbol(symbol string) wc.Series {
	startTime := datetime.FromUnix(int((time.Now().AddDate(0, 0, -1)).Unix()))
	endTime := datetime.FromUnix(int(time.Now().Unix()))
	params := &chart.Params{
		Symbol:   symbol,
		Start:    startTime,
		End:      endTime,
		Interval: datetime.FifteenMins,
	}
	iter := chart.Get(params)

	var XValues []time.Time
	var YValues []float64
	for iter.Next() {
		intTime := int64(iter.Bar().Timestamp)
		XValues = append(XValues, time.Unix(intTime, 0))
		YValues = append(YValues, iter.Bar().Close.InexactFloat64())
	}
	if err := iter.Err(); err != nil {
		fmt.Println(err)
	}
	series := wc.TimeSeries{
		Name:    symbol,
		XValues: XValues,
		YValues: YValues,
	}
	return series
}
