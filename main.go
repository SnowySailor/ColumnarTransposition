package main

import (
    "strings"
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
    "this",
    "test",
    "normies",
    "friendo",
    "dnd",
    "blaise",
    "character",
    "dungeon",
    "master",
}

var permsChan = make(chan []int, 100)
var top = NewTopList()
var wg sync.WaitGroup

func main() {
    fileString := strings.ToLower("LMYAAKIRNPOEWTVTTTXRTUVFSTSSUERRHIOOEEAIUMTAKOAEWMOEATRESETGEATAHTTREHDGUOSLIYHEMBGRSALBUIAETNXNOAAIYAETYDECLBDHNXHESDLENTALUMMAEHOUX")
    fileBytes := []byte(fileString)
    input := string(fileBytes)

    log.Printf("Trying to decode data: %s\n", fileString)

    go generateSequential(100)
    wg.Add(1)
    go generatePerms(9)

    threads := 6
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

        decoded := columnDecrypt(input, perm)
        fitness := determineFitness(decoded)
        if fitness >= 3 {
            top.Add(ListElem{decoded,fitness,perm})
        }
    }
    wg.Done()
}

func columnDecrypt(s1 string, perm []int) string {
    input := strings.ToUpper(s1)
    input_length := len(input)
    
    cols := len(perm)
    rows := int(math.Ceil(float64(input_length)/float64(cols)))

    for {
        if input_length % cols == 0 { break; }
        input_length = input_length + 1
    }
    input_length = input_length - len(input) + 1
    
    mainArray := make([][]rune, rows)
    for i := 0; i < rows; i++ {
        mainArray[i] = make([]rune, cols + 1)
        for s := 0; s < cols + 1; s++ {
            mainArray[i][s] = rune(0)
        }
    }
    
    for j := 0; j < input_length; j++ {
        mainArray[rows - 1][cols - j] = rune(6969)
    }

    kc := 0;
    q := 0;
    for {
        if kc >= len(perm) { break; }
        for i := 0; i < len(perm); i++ {
            if perm[i] == kc {
                for j := 0; j < rows; j++ {
                    if mainArray[j][i] == rune(0) {
                        mainArray[j][i] = rune(input[q])
                        q = q + 1
                    }
                }
                kc = kc + 1
            }
        }
    }

    for j := 0; j < cols + 1; j++ {
        if mainArray[rows-1][j] == rune(6969) {
            mainArray[rows-1][j] = rune(0)
        }
    }
    
    ptext := ""
    for _, row := range mainArray {
        for _, char := range row {
            ptext += string(char)
        }
    }

    return strings.ToLower(ptext)
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
        permsChan <- elements
    }
    log.Printf("Trying sequential keys up to length %d\n", maxLength)
}

func generatePerms(maxLength int) {
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
            permsChan <- tmp
        }
    }

    close(permsChan)

    wg.Done()
}
