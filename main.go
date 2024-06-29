package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os/exec"
    "runtime"
    "time"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
    TotalMemory uint64  `json:"totalMemory"`
    FreeMemory  uint64  `json:"freeMemory"`
    UsedMemory  uint64  `json:"usedMemory"`
    MemoryUsage float64 `json:"memoryUsage"`
    CPUUsage    float64 `json:"cpuUsage"`
}

func getSystemInfo() SystemInfo {
    v, _ := mem.VirtualMemory()
    percent, _ := cpu.Percent(time.Second, false)

    return SystemInfo{
        TotalMemory: v.Total / 1024 / 1024,
        FreeMemory:  v.Free / 1024 / 1024,
        UsedMemory:  v.Used / 1024 / 1024,
        MemoryUsage: v.UsedPercent,
        CPUUsage:    percent[0],
    }
}

func systemInfoHandler(w http.ResponseWriter, r *http.Request) {
    info := getSystemInfo()
    json.NewEncoder(w).Encode(info)
}

func openBrowser(url string) {
    var err error

    switch runtime.GOOS {
    case "linux":
        err = exec.Command("xdg-open", url).Start()
    case "windows":
        err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    default:
        err = fmt.Errorf("unsupported platform")
    }

    if err != nil {
        log.Fatal(err)
    }
}

func printSystemInfo() {
    for {
        info := getSystemInfo()
        fmt.Printf("CPU Usage: %.2f%%, Memory Usage: %.2f%% (Used: %d MB, Total: %d MB)\n",
            info.CPUUsage, info.MemoryUsage, info.UsedMemory, info.TotalMemory)
        time.Sleep(1 * time.Minute)
    }
}

func main() {
    tmpl := template.Must(template.ParseFiles("templates/index.html"))

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    http.HandleFunc("/info", systemInfoHandler)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if err := tmpl.Execute(w, nil); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })

    go printSystemInfo()

    fmt.Println("Starting server at :8080")
    go openBrowser("http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}