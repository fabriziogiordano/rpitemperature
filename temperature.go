package main

import (
    "os"
    "fmt"
    "os/exec"
    "strconv"
    "strings"
    "regexp"
    "time"
    "net/http"
    "net/url"
)

var counter int = 1
var success bool = false
var temp string = "0"
var humi string = "0"

func main() {
    temperature, humidity := temperatureAndHumimidy()
    cpu := cpu()
    gpu := gpu()

    post(temperature, humidity, cpu, gpu)
}

func post(temperature string, humidity string, cpu string, gpu string) int {
    urlGoogle := "https://script.google.com/macros/s/AKfycbwh1TJQjScV54PEfhpZnq--5g6flLVmpAUWnflgkOJA4xCNA88/exec"
    v := url.Values{}
    v.Set("CPU",  cpu)
    v.Set("GPU",  gpu)
    v.Set("Temperature", temperature)
    v.Set("Humidity",  humidity)
    urlGoogle = urlGoogle + "?" + v.Encode()

    response, err := http.Get(urlGoogle)
    defer response.Body.Close()

    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    }

    fmt.Println(temperature, humidity, cpu, gpu)
    return 1
}

func temperatureAndHumimidy() (string, string) {
    app := "sudo"
    arg0 := "/home/pi/Installs/Adafruit-Raspberry-Pi-Python-Code/Adafruit_DHT_Driver/Adafruit_DHT"
    arg1 := "2302"
    arg2 := "25"

    cmd := exec.Command(app, arg0, arg1, arg2)
    out, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return "0","0"
    }

    r, e := regexp.Compile("Temp='(.*)',Hum='(.*)'")
    if e != nil {
        println(err.Error())
        return "0","0"
    }
    res := r.FindStringSubmatch(strings.Trim(string(out), " \n"))

    if len(res) < 1 {
        counter = counter + 1
        if counter > 15 {
            return "0","0"
        }
        time.Sleep(2 * 1010 * time.Millisecond)
        temperatureAndHumimidy()
    }

    if !success {
        temp = res[1]
        humi = res[2]
        success = true
    } else {
        temp = temp
        humi = humi
    }

    return temp, humi
}

func cpu() string {
    app := "cat"
    arg0 := "/sys/class/thermal/thermal_zone0/temp"

    cmd := exec.Command(app, arg0)
    out, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return "0"
    }

    o, e := strconv.Atoi(strings.Trim(string(out), " \n"))

    if e != nil {
        fmt.Println(e.Error())
        return "0"
    }

    t1 := float64(o)/1000

    return strconv.FormatFloat(t1, 'f', 2, 32)
}

func gpu() string {
    app := "/opt/vc/bin/vcgencmd"
    arg0 := "measure_temp"

    cmd := exec.Command(app, arg0)
    out, err := cmd.Output()

    if err != nil {
        println(err.Error())
        return "0"
    }

    r, e := regexp.Compile("temp=(.*)'C")
    if e != nil {
        println(err.Error())
        return "0"
    }

    return string(r.FindStringSubmatch(string(out))[1])
}