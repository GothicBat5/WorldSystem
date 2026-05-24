open System

let readAllInput () =
    Seq.initInfinite (fun _ -> Console.ReadLine())
    |> Seq.takeWhile (fun line -> not (isNull line))
    |> Seq.toList
    |> String.concat " "

let wordFrequency text =
    text.Split([|' '; '\t'; '\n'; '\r'; '.'; ','; ';'; ':'; '!'|], StringSplitOptions.RemoveEmptyEntries)
    |> Seq.map (fun w -> w.ToLowerInvariant())
    |> Seq.countBy id
    |> Seq.sortByDescending snd

[<EntryPoint>]
let main _ =
    printfn "Paste or type text (Ctrl+Z/Ctrl+D to end):"
    let input = readAllInput ()
    let freqs = wordFrequency input

    printfn "\nTop 10 words:"
    freqs
    |> Seq.truncate 10
    |> Seq.iter (fun (word, count) -> printfn "%s: %d" word count)

    0
