package consensus

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/cpacia/obxd/models"
	"sync"
	"testing"
	"time"

	"github.com/cpacia/obxd/net"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
)

func TestAvalancheEngine(t *testing.T) {
	mocknet := mocknet.New(context.Background())
	numNodes := 100
	numNoVotes := 0
	numAlwaysNoVotes := 0

	var (
		engines = make([]*AvalancheEngine, 0, numNodes)
	)

	for i := 0; i < numNodes; i++ {
		host, err := mocknet.GenPeer()
		if err != nil {
			t.Fatal(err)
		}
		network, err := net.NewNetwork(context.Background(), []net.Option{
			net.WithHost(host),
		}...)
		if err != nil {
			t.Fatal(err)
		}

		engine, err := NewAvalancheEngine(context.Background(), network)
		if err != nil {
			t.Fatal(err)
		}
		engines = append(engines, engine)
	}

	if err := mocknet.LinkAll(); err != nil {
		t.Fatal(err)
	}
	if err := mocknet.ConnectAllButSelf(); err != nil {
		t.Fatal(err)
	}

	for _, engine := range engines {
		engine.Start()
	}
	b := make([]byte, 32)
	rand.Read(b)
	chans := make([]chan Status, 0, numNodes)
	start := time.Now()
	for i, engine := range engines {
		chans = append(chans, make(chan Status))
		if i < numAlwaysNoVotes {
			engine.alwaysNo = true
			continue
		}

		engine.NewBlock(models.NewID(b), i >= numNoVotes, chans[i])
	}

	var wg sync.WaitGroup
	wg.Add(numNodes - numAlwaysNoVotes)
	for i := 0; i < numNodes; i++ {
		if i < numAlwaysNoVotes {
			continue
		}
		go func(x int) {
			status := <-chans[x]
			fmt.Printf("Node %d finished as %s\n", x, status)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(time.Since(start))
}
