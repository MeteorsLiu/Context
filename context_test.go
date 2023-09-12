package context

import (
	cntx "context"
	"sync"
	"testing"
	"time"
)

func TestCancel(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	ctx, cancel := WithCancel(cntx.Background())
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		select {
		case <-ticker.C:
			cancel()
		case <-ctx.Done():
			cancel()
			t.Log("5 cancel")
		}

	}()
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(3 * time.Second)
		select {
		case <-ticker.C:
			cancel()
			t.Log("3 stop")
		case <-ctx.Done():
			cancel()
			t.Log("3 cancel")
		}
	}()
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		select {
		case <-ticker.C:
			cancel()
			t.Log("1 stop")
		case <-ctx.Done():
			cancel()
			t.Log("1 cancel")
		}
	}()
	wg.Wait()

	wg.Add(2)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		<-ticker.C
		t.Log("1 start")
		select {
		case <-ticker.C:
			cancel()
			t.Log("1 b stop")
		case <-ctx.Done():
			t.Log("1 b cancel")
		}
	}()
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		<-ticker.C
		t.Log("5 start")
		ticker.Reset(5 * time.Second)
		select {
		case <-ticker.C:
			cancel()
			t.Log("5 b stop")
		case <-ctx.Done():
			t.Log("5 b cancel")
		}
	}()
	wg.Wait()
}

func TestTimeout(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	ctx, cancel := WithTimeout(cntx.Background(), 2*time.Second)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		select {
		case <-ticker.C:
			cancel()
		case <-ctx.Done():
			cancel()
			t.Log("5 cancel")
		}

	}()
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(3 * time.Second)
		select {
		case <-ticker.C:
			cancel()
			t.Log("3 stop")
		case <-ctx.Done():
			cancel()
			t.Log("3 cancel")
		}
	}()
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second)
		select {
		case <-ticker.C:
			cancel()
			t.Log("1 stop")
		case <-ctx.Done():
			cancel()
			t.Log("1 cancel")
		}
	}()
	wg.Wait()
	wg.Add(2)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		<-ticker.C
		t.Log("1 start")
		select {
		case <-ticker.C:
			cancel()
			t.Log("1 b stop")
		case <-ctx.Done():
			t.Log("1 b cancel")
		}
	}()
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		<-ticker.C
		t.Log("5 start")
		ticker.Reset(5 * time.Second)
		select {
		case <-ticker.C:
			cancel()
			t.Log("5 b stop")
		case <-ctx.Done():
			t.Log("5 b cancel")
		}
	}()
	wg.Wait()
}
