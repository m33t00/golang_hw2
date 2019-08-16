package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func sourceData(in, out chan interface{}) {
	source := []int{1, 2, 3, 4, 5, 6, 7}
	for _, sourceVal := range source {
		out <- sourceVal
	}
}

// CombineResults function
func CombineResults(in, out chan interface{}) {
	var stringBuffer []string

	for data := range in {
		stringBuffer = append(stringBuffer, data.(string))
	}

	sort.Strings(stringBuffer)

	out <- strings.Join(stringBuffer, "_")
}

// SingleHash function
func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	md5Mutex := &sync.Mutex{}
	defer wg.Wait()

	for source := range in {
		wg.Add(1)
		go generateSingleHash(strconv.Itoa(source.(int)), wg, md5Mutex, out)
	}
}

func generateSingleHash(source string, wg *sync.WaitGroup, mutex *sync.Mutex, out chan interface{}) {
	var components [2]string
	crc32Wg := &sync.WaitGroup{}
	defer wg.Done()

	mutex.Lock()
	md5 := DataSignerMd5(source)
	mutex.Unlock()

	for idx, data := range [2]string{source, md5} {
		crc32Wg.Add(1)
		go func(idx int, data string) {
			defer crc32Wg.Done()
			components[idx] = DataSignerCrc32(data)
		}(idx, data)
	}

	crc32Wg.Wait()

	out <- strings.Join(components[:], "~")
}

// MultiHash function
func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for data := range in {
		wg.Add(1)
		go generateMultihash(data.(string), wg, out)
	}
}

func generateMultihash(data string, wg *sync.WaitGroup, out chan interface{}) {
	var hashes [6]string
	crc32Wg := &sync.WaitGroup{}
	defer wg.Done()

	for th := range hashes {
		crc32Wg.Add(1)
		go func(th int) {
			defer crc32Wg.Done()
			hashes[th] = DataSignerCrc32(strconv.Itoa(th) + data)
		}(th)
	}

	crc32Wg.Wait()
	out <- strings.Join(hashes[:], "")
}

// ExecutePipeline functiont
func ExecutePipeline(jobs ...job) {
	channels := make([]chan interface{}, len(jobs)+1)
	wg := &sync.WaitGroup{}

	for idx, stageFn := range jobs {
		if idx == 0 {
			channels[idx] = make(chan interface{})
		}
		channels[idx+1] = make(chan interface{})

		wg.Add(1)
		go func(stageFn job, in, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)

			stageFn(in, out)
		}(stageFn, channels[idx], channels[idx+1], wg)
	}
	wg.Wait()

	fmt.Println("Pipeline done")
}
