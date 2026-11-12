package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var md5Mutex = &sync.Mutex{}

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	wg := new(sync.WaitGroup)
	for _, j := range jobs {
		out := make(chan interface{})
		wg.Add(1)

		go func(j job, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)

			j(in, out)
		}(j, in, out)

		in = out
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	var outerWG sync.WaitGroup

	for data := range in {
		outerWG.Add(1)

		go func(data string) {
			defer outerWG.Done()

			var left, right string
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()

				left = DataSignerCrc32(data)
			}()

			go func() {
				defer wg.Done()

				md5Mutex.Lock()
				md5 := DataSignerMd5(data)
				md5Mutex.Unlock()

				right = DataSignerCrc32(md5)
			}()

			wg.Wait()

			out <- left + "~" + right
		}(fmt.Sprint(data))
	}

	outerWG.Wait()
}

func MultiHash(in, out chan interface{}) {
	var outerWG sync.WaitGroup

	for data := range in {
		outerWG.Add(1)

		go func(data string) {
			defer outerWG.Done()

			hashes := make([]string, 6)

			var wg sync.WaitGroup
			wg.Add(6)

			for i := 0; i < 6; i++ {
				go func(n int) {
					defer wg.Done()

					hashes[n] = DataSignerCrc32(
						strconv.Itoa(n) + data,
					)
				}(i)
			}

			wg.Wait()

			out <- strings.Join(hashes, "")
		}(fmt.Sprint(data))
	}

	outerWG.Wait()
}

func CombineResults(in, out chan interface{}) {
	hashes := make([]string, 0)
	for data := range in {
		hashes = append(hashes, data.(string))
	}

	sort.Strings(hashes)

	out <- strings.Join(hashes, "_")
}
