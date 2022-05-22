package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2022
 */

func TestWaitForSignal(_ *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		waitForSignal(ctx, cancel)
		wg.Done()
	}()
	time.Sleep(1 * time.Second)
	cancel()
	wg.Wait()
}
