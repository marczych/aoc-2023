package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type HashMap struct {
	data [][]Node
}

type Node struct {
	key   string
	value int
}

func createHashMap() HashMap {
	hashMap := HashMap{}
	hashMap.data = make([][]Node, 256, 256)
	return hashMap
}

func (hashMap *HashMap) Get(key string) int {
	bucketIndex := hash(key)
	bucket := hashMap.data[bucketIndex]

	for _, node := range bucket {
		if key == node.key {
			return node.value
		}
	}

	return -1
}

func (hashMap *HashMap) Put(key string, value int) {
	bucketIndex := hash(key)
	bucket := hashMap.data[bucketIndex]

	for i := 0; i < len(bucket); i += 1 {
		if key == bucket[i].key {
			bucket[i].value = value
			return
		}
	}

	hashMap.data[bucketIndex] = append(bucket, Node{key, value})
}

func (hashMap *HashMap) Remove(key string) {
	bucketIndex := hash(key)
	bucket := hashMap.data[bucketIndex]
	nodeIndex := -1

	for index, node := range bucket {
		if key == node.key {
			nodeIndex = index
		}
	}

	if nodeIndex != -1 {
		hashMap.data[bucketIndex] = append(bucket[:nodeIndex], bucket[nodeIndex+1:]...)
	}
}

func (hashMap *HashMap) Total() int {
	total := 0
	for bucketIndex, bucket := range hashMap.data {
		for nodeIndex, node := range bucket {
			total += (bucketIndex + 1) * (nodeIndex + 1) * node.value
		}
	}
	return total
}

func hash(input string) int {
	currentValue := 0

	for _, char := range input {
		currentValue += int(char)
		currentValue *= 17
		currentValue %= 256
	}

	return currentValue
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	hashMap := createHashMap()
	for scanner.Scan() {
		line := scanner.Text()
		steps := strings.Split(line, ",")
		for _, step := range steps {
			if step[len(step)-1] == '-' {
				hashMap.Remove(step[:len(step)-1])
			} else if step[len(step)-2] == '=' {
				key := step[:len(step)-2]
				value := int(step[len(step)-1] - '0')
				hashMap.Put(key, value)
			} else {
				log.Fatal("Invalid step: ", step)
			}
		}
	}

	fmt.Println(hashMap.Total())
}
