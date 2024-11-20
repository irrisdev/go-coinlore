package main

import (
	"fmt"
	"log"

	"github.com/irrisdev/go-coinlore/coinlore"
)

func main() {
	
    client := coinlore.NewClient("https://api.coinlore.net")

    // Get global data
    global, err := client.GetGlobal()
    if err != nil {
        log.Fatalf("Error fetching global data: %v", err)
    }
    fmt.Printf("Global Coins Count: %d\n", global.CoinsCount)

    // Get coin data
    coin, err := client.GetCoin(905454)
    if err != nil {
        log.Fatalf("Error fetching coin data: %v", err)
    }
    fmt.Printf("Coin: %s, Price: %s USD\n", coin.Name, coin.PriceUSD)
}