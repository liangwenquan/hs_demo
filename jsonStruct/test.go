package jsonStruct

type Themes struct {
	Code int `json:"code"`
	Message string `json:"message"`
	ThemeEntireInfoList struct {
		ThemeEntireInfoList []struct {
			BasicInfo struct {
				ThemeID int `json:"themeId"`
				ThemeName string `json:"themeName"`
				ThemeDesc string `json:"themeDesc"`
				Event string `json:"event"`
				CurHotIndex float64 `json:"curHotIndex"`
				CurMktDelta float64 `json:"curMktDelta"`
			} `json:"basicInfo"`
			ThemeStockList []struct {
				ThemeID int `json:"themeId"`
				TickerSymbol string `json:"tickerSymbol"`
				SecShortName string `json:"secShortName"`
				CurPrice float64 `json:"curPrice"`
				ChangePct float64 `json:"changePct"`
				MarketValue int64 `json:"marketValue"`
				Suspension int `json:"suspension"`
			} `json:"themeStockList"`
		} `json:"themeEntireInfoList"`
		Count int `json:"count"`
	} `json:"themeEntireInfoList"`
}

