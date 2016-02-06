package main

import "fmt"
import "time"

func collatz() {

    const numOfOdds = 100
    const maxNum = 2*numOfOdds - 1

    var solvedChanArr [numOfOdds]chan int64

    var nextOddArr [numOfOdds]int64

    for i:=0; i < numOfOdds; i++ {

        solvedChanArr[i] = make(chan int64)
    }

    for i:=0; i < numOfOdds; i++ {

        go func(odd int64) {

            num := (odd*3) + 1

            for num %2 == 0 {

                num = num / 2;
            }

            if num > maxNum {

                solvedChanArr[odd/2] <- 2
                
            } else {

                nextOddArr[odd/2] = num

                if num == 1 {

                    solvedChanArr[odd/2] <- 1

                } else {

                    var solved = <- solvedChanArr[num/2]

                    // Broadcast your value
                    go func(s int64) {
                        solvedChanArr[odd/2] <- s
                    } (solved)

                    // Re-broadcast the value you just read
                    go func(s int64) {
                        solvedChanArr[num/2] <- s
                    } (solved)
                }
            }
        }((int64)(2*i + 1))
    }

    var positiveCount = 0
    var unknownCount = 0

    for i:=0; i < numOfOdds; i++ {

        // There is a potential for deadlock here that I'm
        // still trying to fix.  On my machine, it almost never
        // occurs
        solved := <- solvedChanArr[i]
        nextOdd := nextOddArr[i]

        if solved == 1 {

            fmt.Println("Positive:\t", (2*i + 1)," -> ", nextOdd)
            positiveCount++

        } else if solved == 2 {

            fmt.Println("Unknown:\t", (2*i + 1))
            unknownCount++

        } else {

            fmt.Println("What?\t", (2*i + 1))
        }
    }

    // Cleanup
    for i:=0; i < numOfOdds; i++ {

        close(solvedChanArr[i])
    }

    fmt.Println("-----------------------------------------------")
    fmt.Printf("Positive:\t%%%2.2f\tUnknown:\t%%%2.2f\n",
        ((float64)(positiveCount*100)/(float64)(numOfOdds)),
        ((float64)(unknownCount*100)/(float64)(numOfOdds)))

}

func main() {

    var start = time.Now()
    collatz()
    var elapsed = time.Since(start)
    fmt.Println("-----------------------------------------------")
    fmt.Println("Time: ", elapsed)
}
