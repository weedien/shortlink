package main

import (
	"encoding/json"
	"fmt"
	"shortlink/internal/interfaces/rest/dto/resp"
)

func main() {

	var GroupCountList []resp.GroupCountDTO
	g1 := resp.GroupCountDTO{
		Gid:            "public",
		ShortLinkCount: 99,
	}
	g2 := resp.GroupCountDTO{
		Gid:            "test",
		ShortLinkCount: 11,
	}
	GroupCountList = append(GroupCountList, g1, g2)

	r := resp.ShortLinkGroupCountQueryResp(GroupCountList)

	jsonData, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Error converting struct to JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}
