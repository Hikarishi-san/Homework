package main

import (
    "fmt"
    "io"
    "net/http"
    "strconv"
    "strings"
)

func main() {
    url := "http://srv.msk01.gigacorp.local/_stats"
    
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Unable to fetch server statistic.")
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        fmt.Println("Unable to fetch server statistic.")
        return
    }

    body, _ := io.ReadAll(resp.Body)
    line := strings.TrimSpace(string(body))
    
    fields := strings.Split(line, ",")
    if len(fields) != 7 {
        fmt.Println("Unable to fetch server statistic.")
        return
    }

    loadAvg, _ := strconv.ParseFloat(fields[0], 64)
    memTotal, _ := strconv.ParseInt(fields[1], 10, 64)
    memUsed, _ := strconv.ParseInt(fields[2], 10, 64)
    diskTotal, _ := strconv.ParseInt(fields[3], 10, 64)
    diskUsed, _ := strconv.ParseInt(fields[4], 10, 64)
    netTotal, _ := strconv.ParseInt(fields[5], 10, 64)
    netUsed, _ := strconv.ParseInt(fields[6], 10, 64)

    memUsagePercent := float64(memUsed) / float64(memTotal) * 100
    freeDisk := diskTotal - diskUsed
    freeDiskMB := freeDisk / (1024 * 1024)
    netAvailable := netTotal - netUsed
    netAvailableMbit := netAvailable * 8 / (1024 * 1024)

    if loadAvg > 30 {
        fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
    }
    if memUsagePercent > 80 {
        fmt.Printf("Memory usage too high: %.0f%%\n", memUsagePercent)
    }
    if float64(diskUsed) > 0.9*float64(diskTotal) {
        fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMB)
    }
    if float64(netUsed) > 0.9*float64(netTotal) {
        fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", netAvailableMbit)
    }
}
