open System

let maxWords = 10000
let mutable words : string list = []

printfn "Enter some words [type END to finish]"

let rec ask () =
    printf "> "
    let input = Console.ReadLine()
    if isNull input then
        
        finish ()
    else
        let word = input.Trim()
        if word.ToUpper() = "END" then
            finish ()
        elif word = "" then
            printfn "Empty input ignored."
            ask ()
        elif List.length words >= maxWords then
            printfn "Max word limit reached."
            finish ()
        else
            words <- words @ [word]
            ask ()

and finish () =
    if List.isEmpty words then
        printfn "\nNo words entered."
    else
        let sorted =
            words
            |> List.sortBy (fun w -> w.ToLowerInvariant())
        printfn "\nSorted words:"
        sorted |> List.iter (printfn "%s")

// Start
ask ()
