package main

import (
 "fmt"
 "io"
 "net/http"
 "strconv"
 "strings"
 "time"
)

func main() {
 url := "http://srv.msk01.gigacorp.local/_stats"
 errorCount := 0
 maxErrors := 3

 for {
  resp, err := http.Get(url)
  if err != nil || resp == nil {
   errorCount++
   if errorCount >= maxErrors {
    fmt.Println("Unable to fetch server statistic.")
    break
   }
   time.Sleep(100 * time.Millisecond)
   continue
  }

  if resp.StatusCode != 200 {
   resp.Body.Close()
   errorCount++
   if errorCount >= maxErrors {
    fmt.Println("Unable to fetch server statistic.")
    break
   }
   time.Sleep(100 * time.Millisecond)
   continue
  }

  body, err := io.ReadAll(resp.Body)
  resp.Body.Close()

  if err != nil {
   errorCount++
   if errorCount >= maxErrors {
    fmt.Println("Unable to fetch server statistic.")
    break
   }
   time.Sleep(100 * time.Millisecond)
   continue
  }

  line := strings.TrimSpace(string(body))
  fields := strings.Split(line, ",")

  if len(fields) != 7 {
   errorCount++
   if errorCount >= maxErrors {
    fmt.Println("Unable to fetch server statistic.")
    break
   }
   time.Sleep(100 * time.Millisecond)
   continue
  }

  errorCount = 0

  loadAvg, _ := strconv.ParseFloat(fields[0], 64)
  memTotal, _ := strconv.ParseInt(fields[1], 10, 64)
  memUsed, _ := strconv.ParseInt(fields[2], 10, 64)
  diskTotal, _ := strconv.ParseInt(fields[3], 10, 64)
  diskUsed, _ := strconv.ParseInt(fields[4], 10, 64)
  netTotal, _ := strconv.ParseInt(fields[5], 10, 64)
  netUsed, _ := strconv.ParseInt(fields[6], 10, 64)

  memUsagePercent := (memUsed * 100) / memTotal

  freeDisk := diskTotal - diskUsed
  freeDiskMB := freeDisk / (1024 * 1024)

  netAvailable := netTotal - netUsed
  // перевод в мегабиты по десятичной системе (1_000_000 бит = 1 Мбит)
  netAvailableMbit := (netAvailable * 8) / 1_000_000

  if loadAvg > 30 {
   fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
  }
  if memUsagePercent > 80 {
   fmt.Printf("Memory usage too high: %d%%\n", memUsagePercent)
  }
  if (diskUsed * 100) > (diskTotal * 90) {
   fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMB)
  }
  if (netUsed * 100) > (netTotal * 90) {
   fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", netAvailableMbit)
  }

  time.Sleep(100 * time.Millisecond)
 }
}