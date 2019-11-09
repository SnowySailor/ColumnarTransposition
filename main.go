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
    "have",
    "in",
    "good",
    "time",
    "person",
    "people",
    "first",
    "new",
    "world",
    "man",
    "thing",
    "day",
    "way",
    "make",
    "know",
    "about",
    "into",
    "over",
    "after",
    "work",
    "look",
    "want",
    "give",
    "take",
    "steal",
}

var permChan = make(chan []int, 100)
var top = NewTopList()
var wg sync.WaitGroup

func main() {
    fileBytes := []byte(strings.ToLower("IA EDR CI HEAN RXIE EGRSLAD BEOASX TFDNGO DNN STOAXS O IOBKNAEMU ROXTSSTC ATEL TE WDTH THOPMONETHCAN X"))
    input := string(fileBytes)

    go generateSequential(100)
    wg.Add(1)
    go generatePerm(10)

    threads := 8
    for i := 0; i < threads; i++ {
        wg.Add(1)
        go consumePermChan(input)
    }
    wg.Wait()
    log.Println("Done")
    log.Printf("Top 100: %+v\n", top)
}

func consumePermChan(input string) {
    for {
        perm, open := <- permChan
        if (!open) {
            break
        }

        decoded := decode(input, perm)
        fitness := determineFitness(decoded)
        if fitness > 3 {
            top.Add(ListElem{decoded,fitness,perm})
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

    matrix = reorderColumns(matrix, perm)

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

func reorderColumns(matrix [][]rune, perm []int) [][]rune {
    moves := make([]int, len(perm))
    for i := 0; i < len(perm); i++ {
        moves[perm[i]] = i
    }

    cols := len(perm)
    rows := len(matrix)
    newMatrix := make([][]rune, rows)
    for i := 0; i < rows; i++ {
        newMatrix[i] = make([]rune, cols)
    }

    for toColumn := 0; toColumn < len(perm); toColumn++ {
        fromColumn := moves[toColumn]
        for i := 0; i < rows; i++ {
            newMatrix[i][toColumn] = matrix[i][fromColumn]
        }
    }

    return newMatrix
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

func generateSequential(maxLength int) {
    for length := 1; length <= maxLength; length++ {
        elements := make([]int, length)
        for i := 0; i < length; i++ {
            elements[i] = i
        }
        permChan <- elements
    }
    log.Printf("Trying sequential keys up to length %d\n", maxLength)
}

func generatePerm(maxLength int) {
    for length := 1; length <= maxLength; length++ {
        log.Printf("Trying difficulty %d\n", length)

        elements := make([]int, length)
        for i := 0; i < length; i++ {
            elements[i] = i
        }

        p := prmt.New(prmt.IntSlice(elements))
        for p.Next() {
            tmp := make([]int, len(elements))
            copy(tmp, elements)
            permChan <- tmp
        }
    }

    close(permChan)

    wg.Done()
}
