package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func MakeSliceOfChan(length int) []chan interface{} {
	channels := make([]chan interface{}, length)
	for i := range channels {
		channels[i] = make(chan interface{}, 1)
	}
	return channels
}

func MakeSliceOfChanString(length int) []chan string {
	channels := make([]chan string, length)
	for i := range channels {
		channels[i] = make(chan string, 1)
	}
	return channels
}

func ExecutePipeline(jobs ...job) {
	nJobs := len(jobs)
	channels := MakeSliceOfChan(nJobs + 1)
	for i := range jobs {
		go func(function job, in, out chan interface{}) {
			function(in, out)
			close(out)
		}(jobs[i], channels[i], channels[i + 1])
	}
	<-channels[nJobs]
}

func CombineResults(in, out chan interface{}) {
	var results []string
	for result := range in {
		results = append(results, result.(string))
	}
	sort.Strings(results)
	out <- strings.Join(results, "_")
}

type Itoa func(interface{}) string
type Routine func(string, chan interface{}, *sync.WaitGroup)

func DoHash(in, out chan interface{}, itoa Itoa, routine Routine) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for input := range in {
		wg.Add(1)
		go routine(itoa(input), out, wg)
	}
}

func SingleHash(in, out chan interface{}) {
	itoa := Itoa(func(input interface{}) string {
		return strconv.Itoa(input.(int))
	})

	// so much love closures you know
	mutex := make(chan struct{}, 1)
	md5 := func(data string) string {
		mutex <- struct{}{}
		defer func() { <-mutex }()
		return DataSignerMd5(data)
	}

	crc32Routine := func(out chan string, data string) {
		out <- DataSignerCrc32(data)
	}

	routine := Routine(func(data string, out chan interface{}, wg *sync.WaitGroup) {
		defer wg.Done()

		hashes := MakeSliceOfChanString(2)

		go crc32Routine(hashes[0], data)
		go crc32Routine(hashes[1], md5(data))

		out <- <-hashes[0] + "~" + <-hashes[1]
	})

	DoHash(in, out, itoa, routine)
}

func MultiHash(in, out chan interface{}) {
	itoa := Itoa(func(input interface{}) string {
		return input.(string)
	})

	routine := Routine(func(data string, out chan interface{}, wg *sync.WaitGroup) {
		defer wg.Done()

		const N = 6
		hashes := make([]string, N)

		DoMultiHash := func() {
			wg := &sync.WaitGroup{}
			defer wg.Wait()

			for i := 0; i != N; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					hashes[i] = DataSignerCrc32(fmt.Sprintf("%d%s", i, data))
				}(i)  // по сути эту рутину можно было впихнуть в DoHash, а не в DoMultiHash, я нашел этому поистенне чудесное доказательство, но поля файла слишком узки для него...
			}
		}

		DoMultiHash()
		out <- strings.Join(hashes, "")
	})

	DoHash(in, out, itoa, routine)
}
