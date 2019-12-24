package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

func main() {
	// in := gen(2, 3)

	// done := make(chan struct{})
	// defer close(done)

	// out := merge(done, sq(done, in), sq(done, in))
	// fmt.Println(<-out)

	// for n := range out {
	// 	fmt.Println(n)
	// }

	m, err := MD5All(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	var paths []string
	for path := range m {
		paths = append(paths, path)
	}

	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x\t%s\n", m[path], path)
	}

}

func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums))
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(done chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			select {
			case <-done:
				return
			case out <- n * n:
			}
		}
		close(out)
	}()
	return out
}

func merge(done chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int, 1)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func MD5All2(root string) (map[string][md5.Size]byte, error) {
	m := make(map[string][md5.Size]byte)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		m[path] = md5.Sum(data)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}

type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func sumFunc(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	c := make(chan result)
	errc := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}

			wg.Add(1)
			go func() {
				data, err := ioutil.ReadFile(path)
				select {
				case c <- result{path, md5.Sum(data), err}:
				case <-done:
				}
				wg.Done()
			}()

			select {
			case <-done:
				return errors.New("walk canceled")
			default:
				return nil
			}
		})

		go func() {
			wg.Wait()
			close(c)
		}()
		errc <- err
	}()

	return c, errc
}

func MD5All3(root string) (map[string][md5.Size]byte, error) {
	done := make(chan struct{})
	defer close(done)

	c, errc := sumFunc(done, root)

	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}

	return m, nil
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(paths)

		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}

			select {
			case <-done:
				return errors.New("walk canceled")
			case paths <- path:
			}

			return nil
		})

	}()

	return paths, errc
}

func digester(done <-chan struct{}, paths <-chan string, c chan<- result) {
	for path := range paths {
		data, err := ioutil.ReadFile(path)
		select {
		case c <- result{path, md5.Sum(data), err}:
		case <-done:
			return
		}
	}
}

func MD5All(root string) (map[string][md5.Size]byte, error) {

	done := make(chan struct{})
	defer close(done)

	var wg sync.WaitGroup
	const numDigesteres = 20
	wg.Add(numDigesteres)

	paths, errc := walkFiles(done, root)
	c := make(chan result)
	for i := 0; i < numDigesteres; i++ {
		go func() {
			digester(done, paths, c)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}

	return m, nil
}
