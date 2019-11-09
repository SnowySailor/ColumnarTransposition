package main

import (
    "strings"
    //"io/ioutil"
    prmt "github.com/gitchander/permutation"
    "sync"
    "math"
    "log"
)

var words = []string{
    "the",
    "and",
    "bostock",
    "lenindale",
    "an",
    "roast",
    "town",
    "evil",
    "log",
    "then",
    "them",
    "because",
    "be",
    "road",
}

var permsChan = make(chan []int, 100)
var top = NewTopList()
var wg sync.WaitGroup

func main() {
    fileBytes := []byte(strings.ToLower("TSESBCATEORHTSTOKROTNOIHTASDOFOTASETNTIATWHDIBEDODSHNE"))
    input := string(fileBytes)

    wg.Add(1)
    go generatePerms(12)

    threads := 20
    for i := 0; i < threads; i++ {
        wg.Add(1)
    	go consumePermsChan(input)
    }
    wg.Wait()
    log.Println("Done")
    log.Printf("Top 100: %+v\n", top)
}

func consumePermsChan(input string) {
	for {
        perm, open := <- permsChan
        if (!open) {
            break
        }

        //log.Println(perm)
		decoded := decode(input, perm)
		fitness := determineFitness(decoded)
		if fitness > 5 {
			top.Add(ListElem{decoded,fitness})
			//log.Printf("Got a potential match with %v: %s\n ", perm, decoded)
		}
	}
	wg.Done()
}

func decode(input string, perm []int) string {
	cols := len(perm)
	rows := int(math.Ceil(float64(len(input))/float64(cols)))

	matrix := make([][]rune, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]rune, cols)
	}

	for pos, char := range input {
        row := (pos % rows)
        col := (pos / rows)
		matrix[row][col] = char
	}

    // for _, a := range matrix {
    //     acc := ""
    //     for _, b := range a {
    //         acc = acc + string(b)
    //     }
    //     log.Println(acc)
    // }

    acc := ""
    for i := 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            acc += string(matrix[i][j])
        }
    }

    return acc
}

func reorderColumns(matrix [][]int, perms []int) [][]int {
    moves := make([]int, len(perms))
    for i := 0; i < len(perms); i++ {
        idx := indexOf(i, perms)
        moves[perms[idx]] = idx
    }
}

func indexOf(target int, list []int) int {
    for idx, i := range list {
        if i == target {
            return idx
        }
    }
    return -1
}

func determineFitness(str string) int {
    fitness := 0
    for _, search := range words {
        if strings.Contains(str, search) {
            fitness += 1
        }
    }
    return fitness
}

func generatePerms(maxLength int) {
    for length := 1; length < maxLength; length++ {
        log.Printf("Trying difficulty %d\n", length)
        elements := make([]int, length)
        for i := 0; i < length; i++ {
            elements[i] = i
        }

        p := prmt.New(prmt.IntSlice(elements))
        for p.Next() {
            //log.Println(elements)
    		permsChan <- elements
    	}
    }
    close(permsChan)

    wg.Done()
}
