package main

import (
    "fmt"
    "time"
)

func main() {
    var alarmHour, alarmMinute int

   
    fmt.Print("Enter alarm hour : ")
    _, err := fmt.Scan(&alarmHour)

    if err != nil || alarmHour < 0 || alarmHour > 23 {
        fmt.Println("Invalid hour. Please enter a number between 0 and 23.")
        return
    }

 
    fmt.Print("Enter alarm minute : ")
    _, err = fmt.Scan(&alarmMinute)

    if err != nil || alarmMinute < 0 || alarmMinute > 59 {
        fmt.Println("Invalid minute. Please enter a number between 0 and 59.")
        return
    }

    fmt.Printf("Alarm set for %02d:%02d\n", alarmHour, alarmMinute)

    for {
        now := time.Now()
        currentHour := now.Hour()
        currentMinute := now.Minute()

        // Print current time with seconds updating
        fmt.Printf("\rCurrent time: %02d:%02d:%02d", currentHour, currentMinute, now.Second())

        // Check alarm condition
        if currentHour == alarmHour && currentMinute == alarmMinute {
            fmt.Println("\n\nAlarm ringing! Wake up!")
            break
        }

        time.Sleep(1 * time.Second)
    }
}
