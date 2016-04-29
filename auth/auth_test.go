package auth

import (
	"golang.org/x/oauth2"
	"testing"
	"time"
)

var timeProvided, _ = time.Parse(time.RFC3339, "2016-03-25T13:28:37.082536957Z") // 2016-03-25T13:28:37.082536957Z
var validToken = oauth2.Token{
	AccessToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdXBlcmFwcDIiLCJleHAiOjE0NTg3MjQ0MTMsImlhdCI6MTQ1ODcyMDgxMywiaXNzIjoiIiwiamlkIjoiZTE3NTNhMzYtMzU0NC00M2Q2LTk4MTUtMDNkNWMyMDliM2YyIiwibmJmIjoxNDU4NzIwODEzLCJzdWIiOiIyNDVhYjIyZi03ODYyLTRiNGQtODM2Mi0yMTRkZmE5ZGUyYzkifQ.GYYtd2e41_Ny8j-MdKgemifSPuzkzSdVEjScX5acQQBRI28gSrLg55wExYtvRz1295SGSW5mKU-wDvPrlP6csxjj70axMDtCt47rTHSQut5jcWDPDRpN_4wD9TUfLH15c_VQNp39yntZgSU_ygH61z7VuxSCu-EGbGxfREK2aBPpPp1FWMo7QiQ97oqcLcvCWvibHyRTr7gFeswo6yy6KO8kGBeJMPX4gPHFstO2Ghiod2VolWy1RyOoKhbT5iO9CgcdfNP5KB4ba9hES-mV0PXoxUCDCeH_EnWKXqXY2STDmEy6k1F0ye1UbD1dTnbo9hSwr7Zq8hWlkL6iTZM6-w",
	TokenType:    "Bearer",
	RefreshToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpZCI6IjMwMWI2NjZjLTk2NzYtNGUxZC04Y2ZkLTk1Yjc4ZTM1OWNhOCJ9.sNctk7YF5qw-B6zrvdYf-GQ00EVabgRrpWH4YYttzs5p1hhxbYiSJF7ESCRvWe8ttk7Ojq3OS1jOi7RqvlRdi5ibAuD_zYVJo-iSa84VeubpXzQYzg0ri19n_HDGRdUnGl-6lV9SK7UBJwm3685pm3o-MkqFagwRhbKc5HtMcC0ho683Lx5scpMJsUBtYB5uKv1MsQTeLxK8lBPpSvgQGe7SKt9fUhnri-jF8IgWdth3J-bRw61p4BQib6-KOcJrUsKa3J5legvgSDmT03Btr5l_dus6dIPTwyGYuiFKD71zzBQggTrzSdOeiqRz4MkPVBmB09LSnq6ijcgCaRjf3Q",
	Expiry:       timeProvided,
}
var validTokenInfoString = "{\"access_token\":\"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdXBlcmFwcDIiLCJleHAiOjE0NTg3MjQ0MTMsImlhdCI6MTQ1ODcyMDgxMywiaXNzIjoiIiwiamlkIjoiZTE3NTNhMzYtMzU0NC00M2Q2LTk4MTUtMDNkNWMyMDliM2YyIiwibmJmIjoxNDU4NzIwODEzLCJzdWIiOiIyNDVhYjIyZi03ODYyLTRiNGQtODM2Mi0yMTRkZmE5ZGUyYzkifQ.GYYtd2e41_Ny8j-MdKgemifSPuzkzSdVEjScX5acQQBRI28gSrLg55wExYtvRz1295SGSW5mKU-wDvPrlP6csxjj70axMDtCt47rTHSQut5jcWDPDRpN_4wD9TUfLH15c_VQNp39yntZgSU_ygH61z7VuxSCu-EGbGxfREK2aBPpPp1FWMo7QiQ97oqcLcvCWvibHyRTr7gFeswo6yy6KO8kGBeJMPX4gPHFstO2Ghiod2VolWy1RyOoKhbT5iO9CgcdfNP5KB4ba9hES-mV0PXoxUCDCeH_EnWKXqXY2STDmEy6k1F0ye1UbD1dTnbo9hSwr7Zq8hWlkL6iTZM6-w\",\"token_type\":\"Bearer\",\"refresh_token\":\"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpZCI6IjMwMWI2NjZjLTk2NzYtNGUxZC04Y2ZkLTk1Yjc4ZTM1OWNhOCJ9.sNctk7YF5qw-B6zrvdYf-GQ00EVabgRrpWH4YYttzs5p1hhxbYiSJF7ESCRvWe8ttk7Ojq3OS1jOi7RqvlRdi5ibAuD_zYVJo-iSa84VeubpXzQYzg0ri19n_HDGRdUnGl-6lV9SK7UBJwm3685pm3o-MkqFagwRhbKc5HtMcC0ho683Lx5scpMJsUBtYB5uKv1MsQTeLxK8lBPpSvgQGe7SKt9fUhnri-jF8IgWdth3J-bRw61p4BQib6-KOcJrUsKa3J5legvgSDmT03Btr5l_dus6dIPTwyGYuiFKD71zzBQggTrzSdOeiqRz4MkPVBmB09LSnq6ijcgCaRjf3Q\",\"expiry\":\"2016-03-25T13:28:37.082536957Z\"}"
var validTokenInfo = &TokenInfo{
	TokenInfo: string(validTokenInfoString),
}
var corruptedTokenInfoString = "{\"access_tokenmodified\":\"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdXBlcmFwcDIiLCJleHAiOjE0NTg3MjQ0MTMsImlhdCI6MTQ1ODcyMDgxMywiaXNzIjoiIiwiamlkIjoiZTE3NTNhMzYtMzU0NC00M2Q2LTk4MTUtMDNkNWMyMDliM2YyIiwibmJmIjoxNDU4NzIwODEzLCJzdWIiOiIyNDVhYjIyZi03ODYyLTRiNGQtODM2Mi0yMTRkZmE5ZGUyYzkifQ.GYYtd2e41_Ny8j-MdKgemifSPuzkzSdVEjScX5acQQBRI28gSrLg55wExYtvRz1295SGSW5mKU-wDvPrlP6csxjj70axMDtCt47rTHSQut5jcWDPDRpN_4wD9TUfLH15c_VQNp39yntZgSU_ygH61z7VuxSCu-EGbGxfREK2aBPpPp1FWMo7QiQ97oqcLcvCWvibHyRTr7gFeswo6yy6KO8kGBeJMPX4gPHFstO2Ghiod2VolWy1RyOoKhbT5iO9CgcdfNP5KB4ba9hES-mV0PXoxUCDCeH_EnWKXqXY2STDmEy6k1F0ye1UbD1dTnbo9hSwr7Zq8hWlkL6iTZM6-w\",\"token_type\":\"Bearer\",\"refresh_token\":\"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpZCI6IjMwMWI2NjZjLTk2NzYtNGUxZC04Y2ZkLTk1Yjc4ZTM1OWNhOCJ9.sNctk7YF5qw-B6zrvdYf-GQ00EVabgRrpWH4YYttzs5p1hhxbYiSJF7ESCRvWe8ttk7Ojq3OS1jOi7RqvlRdi5ibAuD_zYVJo-iSa84VeubpXzQYzg0ri19n_HDGRdUnGl-6lV9SK7UBJwm3685pm3o-MkqFagwRhbKc5HtMcC0ho683Lx5scpMJsUBtYB5uKv1MsQTeLxK8lBPpSvgQGe7SKt9fUhnri-jF8IgWdth3J-bRw61p4BQib6-KOcJrUsKa3J5legvgSDmT03Btr5l_dus6dIPTwyGYuiFKD71zzBQggTrzSdOeiqRz4MkPVBmB09LSnq6ijcgCaRjf3Q\",\"expiry\":\"2016-03-25T13:28:37.082536957Z\"}"
var corruptedTokenInfo = &TokenInfo{
	TokenInfo: string(corruptedTokenInfoString),
}

func TestEncodeToken(t *testing.T) {
	//test that converts to access_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdXBlcmFwcDIiLCJleHAiOjE0NTg5MTI1MTcsImlhdCI6MTQ1ODkwODkxNywiaXNzIjoiIiwiamlkIjoiMTJkNWVkNDctMWI2MC00ODViLWFkOGMtZjMwMGRmZWU5OGZlIiwibmJmIjoxNDU4OTA4OTE3LCJzdWIiOiIzYTRjYzE5Mi03Yzc0LTQ3MjItYjkxOS1jODc2NDVmMTc3N2IifQ.qfwser5GmDFMhoaPw-swlmecCuIcQPe4PLmvLDUC5zxxDwEKjT0TdCh3-an1YZrhNNyWQEFZjHM1NQ7NCUyaqA1G7gu5h_mjVGKagZMVt-q-ucU4ltWTi3rhc-MUySG5LsHOv8qVph2mW0RAA6RQQftnxB_ury_MZ4pKRrVK49C1O_YKwPfKzbSWcnGQvg5X2-By016SsW45St5krNfWabKGOYTBIbc1CDvRx8rrI01MYETpXBAY1JcEJCArjf6NSYGIfXF7P7B-OfOkkUvUJ8D3GucqkZ_zanqZTcj3XGPsTjIyoNCzTykrCKN_PFNFUPQ1dQDa4C-QJ5A8hy10YA,token_type=Bearer,refresh_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpZCI6IjAxMzhiYzc2LTRkODgtNGU2My1iY2U1LTcwZGNhNzljYjhhYiJ9.RmFWvGH0mFz9nQ7cnBOLqSked7MeM-ct2QZ0ATlAPok_95qlTCJfGp41sY2TRsvyQR9sAOeta30T4sGRPgL8Y3H3iW9cBKUdZuOhHohSwDv_jjoGSp8hNuLFTNzYdqKu2Tokv2L4HHPQMky1nXyz8pWSqu6phx5KN7D8Ti5oV80oBO503CjhDGCRwCtA4E3JFpCl49-D06dF_meCOldj1tLPvcaQ0DCF3wPToT9D9Peiwppj6-bfpDGtbCj49BM_wd7Ibgsy2aPdMCoiAOGmFi8Vmo-z7aGQRdvohG_TuRdmAUIjMzhA0445cTzfsDu6hbFPRmjf0tJNywkptYCVeQ,expiry=2016-03-25T13:28:37.082536957Z
	//expectedTokenInfo := "access_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdXBlcmFwcDIiLCJleHAiOjE0NTg5MTI1MTcsImlhdCI6MTQ1ODkwODkxNywiaXNzIjoiIiwiamlkIjoiMTJkNWVkNDctMWI2MC00ODViLWFkOGMtZjMwMGRmZWU5OGZlIiwibmJmIjoxNDU4OTA4OTE3LCJzdWIiOiIzYTRjYzE5Mi03Yzc0LTQ3MjItYjkxOS1jODc2NDVmMTc3N2IifQ.qfwser5GmDFMhoaPw-swlmecCuIcQPe4PLmvLDUC5zxxDwEKjT0TdCh3-an1YZrhNNyWQEFZjHM1NQ7NCUyaqA1G7gu5h_mjVGKagZMVt-q-ucU4ltWTi3rhc-MUySG5LsHOv8qVph2mW0RAA6RQQftnxB_ury_MZ4pKRrVK49C1O_YKwPfKzbSWcnGQvg5X2-By016SsW45St5krNfWabKGOYTBIbc1CDvRx8rrI01MYETpXBAY1JcEJCArjf6NSYGIfXF7P7B-OfOkkUvUJ8D3GucqkZ_zanqZTcj3XGPsTjIyoNCzTykrCKN_PFNFUPQ1dQDa4C-QJ5A8hy10YA,token_type=Bearer,refresh_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpZCI6IjAxMzhiYzc2LTRkODgtNGU2My1iY2U1LTcwZGNhNzljYjhhYiJ9.RmFWvGH0mFz9nQ7cnBOLqSked7MeM-ct2QZ0ATlAPok_95qlTCJfGp41sY2TRsvyQR9sAOeta30T4sGRPgL8Y3H3iW9cBKUdZuOhHohSwDv_jjoGSp8hNuLFTNzYdqKu2Tokv2L4HHPQMky1nXyz8pWSqu6phx5KN7D8Ti5oV80oBO503CjhDGCRwCtA4E3JFpCl49-D06dF_meCOldj1tLPvcaQ0DCF3wPToT9D9Peiwppj6-bfpDGtbCj49BM_wd7Ibgsy2aPdMCoiAOGmFi8Vmo-z7aGQRdvohG_TuRdmAUIjMzhA0445cTzfsDu6hbFPRmjf0tJNywkptYCVeQ,expiry=2016-03-25T13:28:37.082536957Z"
	actualTokenInfo, err := EncodeTokenInfo(&validToken)
	if err != nil {
		t.Error("got error when trying to encode correct token : %v", err)
	}
	if validTokenInfo.TokenInfo != actualTokenInfo.TokenInfo {
		t.Errorf("Encoded tokenInfo does not match, \nexpected:%v, \nactual:%v ", validTokenInfo, actualTokenInfo)
	}

	//test that nil token return nil
	nilTokenInfo, errNil := EncodeTokenInfo(nil)
	if errNil != nil {
		t.Error("got error when trying to encode nil token, should have got empty string : %v", err)
	}
	if nilTokenInfo != nil {
		t.Errorf("When trying to encode nil token, should have got nil tokenInfo")
	}

	//test that invalid token return err
	//invalidTokenWithoutAccessToken := oauth2.Token{
	//	TokenType:    "Bearer",
	//	RefreshToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpZCI6IjMwMWI2NjZjLTk2NzYtNGUxZC04Y2ZkLTk1Yjc4ZTM1OWNhOCJ9.sNctk7YF5qw-B6zrvdYf-GQ00EVabgRrpWH4YYttzs5p1hhxbYiSJF7ESCRvWe8ttk7Ojq3OS1jOi7RqvlRdi5ibAuD_zYVJo-iSa84VeubpXzQYzg0ri19n_HDGRdUnGl-6lV9SK7UBJwm3685pm3o-MkqFagwRhbKc5HtMcC0ho683Lx5scpMJsUBtYB5uKv1MsQTeLxK8lBPpSvgQGe7SKt9fUhnri-jF8IgWdth3J-bRw61p4BQib6-KOcJrUsKa3J5legvgSDmT03Btr5l_dus6dIPTwyGYuiFKD71zzBQggTrzSdOeiqRz4MkPVBmB09LSnq6ijcgCaRjf3Q",
	//	Expiry:       timeProvided,
	//}
	//invalidTokenWithoutTokenType := oauth2.Token{
	//	AccessToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdXBlcmFwcDIiLCJleHAiOjE0NTg3MjQ0MTMsImlhdCI6MTQ1ODcyMDgxMywiaXNzIjoiIiwiamlkIjoiZTE3NTNhMzYtMzU0NC00M2Q2LTk4MTUtMDNkNWMyMDliM2YyIiwibmJmIjoxNDU4NzIwODEzLCJzdWIiOiIyNDVhYjIyZi03ODYyLTRiNGQtODM2Mi0yMTRkZmE5ZGUyYzkifQ.GYYtd2e41_Ny8j-MdKgemifSPuzkzSdVEjScX5acQQBRI28gSrLg55wExYtvRz1295SGSW5mKU-wDvPrlP6csxjj70axMDtCt47rTHSQut5jcWDPDRpN_4wD9TUfLH15c_VQNp39yntZgSU_ygH61z7VuxSCu-EGbGxfREK2aBPpPp1FWMo7QiQ97oqcLcvCWvibHyRTr7gFeswo6yy6KO8kGBeJMPX4gPHFstO2Ghiod2VolWy1RyOoKhbT5iO9CgcdfNP5KB4ba9hES-mV0PXoxUCDCeH_EnWKXqXY2STDmEy6k1F0ye1UbD1dTnbo9hSwr7Zq8hWlkL6iTZM6-w",
	//	RefreshToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpZCI6IjMwMWI2NjZjLTk2NzYtNGUxZC04Y2ZkLTk1Yjc4ZTM1OWNhOCJ9.sNctk7YF5qw-B6zrvdYf-GQ00EVabgRrpWH4YYttzs5p1hhxbYiSJF7ESCRvWe8ttk7Ojq3OS1jOi7RqvlRdi5ibAuD_zYVJo-iSa84VeubpXzQYzg0ri19n_HDGRdUnGl-6lV9SK7UBJwm3685pm3o-MkqFagwRhbKc5HtMcC0ho683Lx5scpMJsUBtYB5uKv1MsQTeLxK8lBPpSvgQGe7SKt9fUhnri-jF8IgWdth3J-bRw61p4BQib6-KOcJrUsKa3J5legvgSDmT03Btr5l_dus6dIPTwyGYuiFKD71zzBQggTrzSdOeiqRz4MkPVBmB09LSnq6ijcgCaRjf3Q",
	//	Expiry:       timeProvided,
	//}
	//invalidTokenWithoutRefreshToken := oauth2.Token{
	//	AccessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdXBlcmFwcDIiLCJleHAiOjE0NTg3MjQ0MTMsImlhdCI6MTQ1ODcyMDgxMywiaXNzIjoiIiwiamlkIjoiZTE3NTNhMzYtMzU0NC00M2Q2LTk4MTUtMDNkNWMyMDliM2YyIiwibmJmIjoxNDU4NzIwODEzLCJzdWIiOiIyNDVhYjIyZi03ODYyLTRiNGQtODM2Mi0yMTRkZmE5ZGUyYzkifQ.GYYtd2e41_Ny8j-MdKgemifSPuzkzSdVEjScX5acQQBRI28gSrLg55wExYtvRz1295SGSW5mKU-wDvPrlP6csxjj70axMDtCt47rTHSQut5jcWDPDRpN_4wD9TUfLH15c_VQNp39yntZgSU_ygH61z7VuxSCu-EGbGxfREK2aBPpPp1FWMo7QiQ97oqcLcvCWvibHyRTr7gFeswo6yy6KO8kGBeJMPX4gPHFstO2Ghiod2VolWy1RyOoKhbT5iO9CgcdfNP5KB4ba9hES-mV0PXoxUCDCeH_EnWKXqXY2STDmEy6k1F0ye1UbD1dTnbo9hSwr7Zq8hWlkL6iTZM6-w",
	//	TokenType:   "Bearer",
	//	Expiry:      timeProvided,
	//}
	//invalidTokenWithoutExpiry := oauth2.Token{
	//	AccessToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzdXBlcmFwcDIiLCJleHAiOjE0NTg3MjQ0MTMsImlhdCI6MTQ1ODcyMDgxMywiaXNzIjoiIiwiamlkIjoiZTE3NTNhMzYtMzU0NC00M2Q2LTk4MTUtMDNkNWMyMDliM2YyIiwibmJmIjoxNDU4NzIwODEzLCJzdWIiOiIyNDVhYjIyZi03ODYyLTRiNGQtODM2Mi0yMTRkZmE5ZGUyYzkifQ.GYYtd2e41_Ny8j-MdKgemifSPuzkzSdVEjScX5acQQBRI28gSrLg55wExYtvRz1295SGSW5mKU-wDvPrlP6csxjj70axMDtCt47rTHSQut5jcWDPDRpN_4wD9TUfLH15c_VQNp39yntZgSU_ygH61z7VuxSCu-EGbGxfREK2aBPpPp1FWMo7QiQ97oqcLcvCWvibHyRTr7gFeswo6yy6KO8kGBeJMPX4gPHFstO2Ghiod2VolWy1RyOoKhbT5iO9CgcdfNP5KB4ba9hES-mV0PXoxUCDCeH_EnWKXqXY2STDmEy6k1F0ye1UbD1dTnbo9hSwr7Zq8hWlkL6iTZM6-w",
	//	TokenType:    "Bearer",
	//	RefreshToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpZCI6IjMwMWI2NjZjLTk2NzYtNGUxZC04Y2ZkLTk1Yjc4ZTM1OWNhOCJ9.sNctk7YF5qw-B6zrvdYf-GQ00EVabgRrpWH4YYttzs5p1hhxbYiSJF7ESCRvWe8ttk7Ojq3OS1jOi7RqvlRdi5ibAuD_zYVJo-iSa84VeubpXzQYzg0ri19n_HDGRdUnGl-6lV9SK7UBJwm3685pm3o-MkqFagwRhbKc5HtMcC0ho683Lx5scpMJsUBtYB5uKv1MsQTeLxK8lBPpSvgQGe7SKt9fUhnri-jF8IgWdth3J-bRw61p4BQib6-KOcJrUsKa3J5legvgSDmT03Btr5l_dus6dIPTwyGYuiFKD71zzBQggTrzSdOeiqRz4MkPVBmB09LSnq6ijcgCaRjf3Q",
	//}
	//actualInvalidTokenInfo, errInvalid := EncodeTokenInfo(&invalidTokenWithoutAccessToken)
	//if errInvalid == nil {
	//	t.Error("Should have got an error when trying to encode invalid token without access token")
	//}
	//if actualInvalidTokenInfo != "" {
	//	t.Errorf("Should have got an empty tokenInfo when trying to encode invalid token without access token\nactual:%v ", actualInvalidTokenInfo)
	//}
	//actualInvalidTokenInfo, errInvalid = EncodeTokenInfo(&invalidTokenWithoutTokenType)
	//if errInvalid == nil {
	//	t.Error("Should have got an error when trying to encode invalid token without token type")
	//}
	//if actualInvalidTokenInfo != "" {
	//	t.Errorf("Should have got an empty tokenInfo when trying to encode invalid token without token type\nactual:%v ", actualInvalidTokenInfo)
	//}
	//actualInvalidTokenInfo, errInvalid = EncodeTokenInfo(&invalidTokenWithoutRefreshToken)
	//if errInvalid == nil {
	//	t.Error("Should have got an error when trying to encode invalid token without refresh token")
	//}
	//if actualInvalidTokenInfo != "" {
	//	t.Errorf("Should have got an empty tokenInfo when trying to encode invalid token without refresh token\nactual:%v ", actualInvalidTokenInfo)
	//}
	//actualInvalidTokenInfo, errInvalid = EncodeTokenInfo(&invalidTokenWithoutExpiry)
	//if errInvalid == nil {
	//	t.Error("Should have got an error when trying to encode invalid token without expiry")
	//}
	//if actualInvalidTokenInfo != "" {
	//	t.Errorf("Should have got an empty tokenInfo when trying to encode invalid token without expiry\nactual:%v ", actualInvalidTokenInfo)
	//}

}

func TestDecodeToken(t *testing.T) {
	//test valid
	//test nil
	actualToken, err := DecodeTokenInfo(validTokenInfo)
	if err != nil {
		t.Error("Got an error when tryin to decode a valid TokenInfo", err)
	}
	if actualToken == nil {
		t.Fatal("Should have got a non nil Token when Decoding a valid tokeninfo")
	}
	if actualToken.AccessToken != validToken.AccessToken &&
		!actualToken.Expiry.Equal(validToken.Expiry) &&
		actualToken.RefreshToken != validToken.RefreshToken &&
		actualToken.TokenType != validToken.TokenType {
		t.Errorf("Got different result when trying to decode tokenInfo\nExpectedAccessToken: %v\nActualAccessToken: %v\nExpectedExpiry:%v\nActual Expiry:%v\nExpected refresh:%v\nActual refresh:%vExpected TokenType:%v\nActual TokenType:%v",
			validToken.AccessToken, actualToken.AccessToken,
			validToken.Expiry.Format(time.RFC3339), actualToken.Expiry.Format(time.RFC3339),
			validToken.RefreshToken, actualToken.RefreshToken,
			validToken.TokenType, actualToken.TokenType)
	}

	invalidToken, errNil := DecodeTokenInfo(nil)
	if errNil != nil {
		t.Error("got error when trying to decode nil tokenInfo, should have got nil token : ", err)
	}
	if invalidToken != nil {
		t.Error("When trying to decode empty tokenInfo, should have got nil token")
	}

	invalidToken, errNil = DecodeTokenInfo(&TokenInfo{})
	if errNil == nil {
		t.Error("should have got and error when trying to decode tokenInfo with empty tokenInfo, should have got nil token")
	}
	if invalidToken != nil {
		t.Error("When trying to decode empty tokenInfo, should have got nil token, instead got :", invalidToken.TokenType)
	}

	//test corrupted tokenInfo
	actualCorruptedToken, errCorrupted := DecodeTokenInfo(corruptedTokenInfo)
	if errCorrupted != nil {
		t.Error("Should have not got an error to decode an invalid token, instead got : ", errCorrupted)
	}
	if actualCorruptedToken.AccessToken != "" {
		t.Error("Should have got an empty access_token in resulting token when trying to decode an invalid tokeninfo, instead got access_token : ", actualCorruptedToken.AccessToken)
	}

}
