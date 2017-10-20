package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

//https://oauthplay.azurewebsites.net/でトークンを取る。
//まずはそこから始めよう。
func main() {
	urls := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Add("code", "OAQABAAIAAABHh4kmS_aKT5XrjzxRAtHz-bbO_fNFrtTx4cRSLZkQlo0BaJ9K46mGD9AR7ZbycPoChN5txboNG3rf6tHmyRGt5ie6Mac2DPiWG-17NfZBXhjA_39jWGUjXi57AN0gbOMdET21bH2gyldU7mVtJuho7emXKvmtHEmmU1LNOsdKXtwHj2MFzj2MgMUtFQTOGPkkV12nH4j0hg_qsS58g-BUHWRjloZtfUW_dkdfYT6AICmHBRzhZpQCfCEa1n4aLYWrt1S-ly7gZUrpH3_CrWCqPCrPqh_sfVg0i2ZA4-M8V0OffqUJOOG2J1BsoD0AuvNZV2Wg4o3RHnO_DmYWSHVg33FRe8viKzxXpL99x_SiKOxnb6xr8C9wE3M6WsTwOlKsx_KIzky7kacmskvHIrttFVgHaC3JUt0OW6fluo1AyNhNXiuL0bu3gyei6ndxgRUY2PyWBKBoDzzPYrQaLUiGByzja35IbS0ZyBArLpE_4HNP7vEfMyiPboHNq0lyEbc2wFeEAZK15J2dJAILLvXDB5q-Zj4Z1Te5iy8WaUehV4YqDPf89Jd8L674eaGR_hHYDx4Y8M8CZqBu6HzuR-seTMx7CcO6gkyVU4FwM_ouNxTxD8aXQqWJHOJiCXDYh9dL--LM17IIwUYbcefw8HWLQWiEOgkXX2bmkS55HQ2n_yVghpddjcmq7a9ujVMLtaREkDRqWdzDUpocq58EqtVavzOzHv-uyo9dQHRaWyxMNaovkHH5GqNQgivHAREl_ZIgAA")
	values.Add("scope", "user.read mail.read contacts.read calendars.read")
	values.Add("client_id", "0a77c1da-b4d0-48d3-a71e-07005dd5c429")
	values.Add("redirect_uri", "https://service.intra.tsunagunet.com:8443")
	values.Add("client_secret", "YcbWtns5OvfxVe26Z4KN9nn")

	req, _ := http.NewRequest("POST", urls, strings.NewReader(values.Encode()))
	log.Printf("req1 : %s\n", req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dump, _ := httputil.DumpRequestOut(req, true)
	log.Printf("req set header : %s\n", dump)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	dumpResp, _ := httputil.DumpResponse(resp, true)
	log.Printf("resp : %s\n", dumpResp)

	byteArray, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(byteArray))
}
