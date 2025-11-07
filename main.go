package main

import (
    "bufio"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "strings"
)

func main() {
    url := "http://srv.msk01.gigacorp.local/_stats"
    maxErrors := 3
    errorCount := 0

    for {
        resp, err := http.Get(url)
        if err != nil || resp.StatusCode != 200 {
            errorCount++
            if errorCount >= maxErrors {
                fmt.Println("Unable to fetch server statistic.")
                os.Exit(1)
            }
            continue
        }
        defer resp.Body.Close()

        scanner := bufio.NewScanner(resp.Body)
        var line string
        if scanner.Scan() {
            line = scanner.Text()
        } else {
            errorCount++
            if errorCount >= maxErrors {
                fmt.Println("Unable to fetch server statistic.")
                os.Exit(1)
            }
            continue
        }

        fields := strings.Split(line, ",")
        if len(fields) != 7 {
            errorCount++
            if errorCount >= maxErrors {
                fmt.Println("Unable to fetch server statistic.")
                os.Exit(1)
            }
            continue
        }

        // Парсим значения
        loadAvg, _ := strconv.ParseFloat(fields[0], 64)
        memTotal, _ := strconv.ParseInt(fields[1], 10, 64)
        memUsed, _ := strconv.ParseInt(fields[2], 10, 64)
        diskTotal, _ := strconv.ParseInt(fields[3], 10, 64)
        diskUsed, _ := strconv.ParseInt(fields[4], 10, 64)
        netTotal, _ := strconv.ParseInt(fields[5], 10, 64)
        netUsed, _ := strconv.ParseInt(fields[6], 10, 64)

        errorCount = 0

        // Проверяем пороги и выводим сообщения
        if loadAvg > 30 {
            fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
        }
        memUsagePercent := float64(memUsed) / float64(memTotal) * 100
        if memUsagePercent > 80 {
            fmt.Printf("Memory usage too high: %.0f%%\n", memUsagePercent)
        }
        freeDisk := diskTotal - diskUsed
        freeDiskMB := freeDisk / (1024 * 1024)
        if float64(diskUsed) > 0.9*float64(diskTotal) {
            fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMB)
        }
        netAvailable := netTotal - netUsed
        netAvailableMbit := netAvailable * 8 / (1024 * 1024) // биты в мегабиты
        if float64(netUsed) > 0.9*float64(netTotal) {
            fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", netAvailableMbit)
        }

        break
    }
}
