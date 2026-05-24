open System
open System.IO
open System.Collections.Generic


let roles = Dictionary<string, List<string>>()

let trim (s: string) = s.Trim()

let capitalize (str: string) = if String.IsNullOrWhiteSpace(str) then ""
    else
        let first = Char.ToUpper(str.[0])
        let rest = str.Substring(1).ToLower()
        string first + rest

let input (prompt: string) = Console.Write(prompt)
    trim (Console.ReadLine())


let addPerson () = let name = input "Enter name: "
    if String.IsNullOrWhiteSpace(name) then
        printfn "Name cannot be empty.\n"
    else
        let role = capitalize (input "Enter role: ")
        if String.IsNullOrWhiteSpace(role) then
            printfn "Role cannot be empty.\n"
        else
            if not (roles.ContainsKey(role)) then
                roles.[role] <- List<string>()
            roles.[role].Add(name)
            printfn "Added %s as %s\n" name role

let showGroups () = printfn "\n------ People Groups ------"
    let sortedRoles = roles.Keys |> Seq.sort |> Seq.toList
    if sortedRoles.IsEmpty then
        printfn "No group found.\n"
    else
        for role in sortedRoles do
            let people = roles.[role] |> Seq.sort |> String.concat ", "
            printfn "%ss: %s" role people
        printfn ""

let deletePerson () = let name = input "Enter name to delete: "
    let mutable found = false
    for role in roles.Keys do
        let people = roles.[role]
        let idx = people.FindIndex(fun p -> p.ToLower() = name.ToLower())
        if idx <> -1 then
            let removed = people.[idx]
            people.RemoveAt(idx)
            printfn "Removed %s from %s\n" removed role
            found <- true
    if not found then
        printfn "Person not found.\n"

let searchPerson () = let name = input "Search name: "
    let mutable found = false
    for role in roles.Keys do
        for person in roles.[role] do
            if person.ToLower().Contains(name.ToLower()) then
                printfn "%s is a %s" person role
                found <- true
    if not found then
        printfn "No person found.\n"
    else
        printfn ""

let saveToFile () =  try
        use writer = new StreamWriter("roles.txt")
        for role in roles.Keys do
            for person in roles.[role] do
                writer.WriteLine($"{role}|{person}")
        printfn "Data saved!\n"
    with _ ->
        printfn "Error saving file.\n"

let loadFromFile () =  if File.Exists("roles.txt") then
        let lines = File.ReadAllLines("roles.txt")
        for line in lines do
            let parts = line.Split('|')
            if parts.Length = 2 then
                let role, name = parts.[0], parts.[1]
                if not (roles.ContainsKey(role)) then
                    roles.[role] <- List<string>()
                roles.[role].Add(name)

let showMenu () =
    printfn "===== MENU =====" 
    
    printfn "1. Add Person"
    
    printfn "2. Show Groups"
    
    printfn "3. Delete Person"
    
    printfn "4. Search Person"
    
    printfn "5. Save"
    
    printfn "6. Quit"

let rec mainLoop () =
    showMenu()
    let choice = input "Choose: "
    match choice with
    | "1" -> addPerson(); mainLoop()
    
    | "2" -> showGroups(); mainLoop()
    
    | "3" -> deletePerson(); mainLoop()
    
    | "4" -> searchPerson(); mainLoop()
    
    | "5" -> saveToFile(); mainLoop()
    
    | "6" -> printfn "\nProgram ended here."
    | _   -> printfn "Invalid choice.\n"; mainLoop()

[<EntryPoint>]
let main _ =
    loadFromFile()
    mainLoop()
    0
