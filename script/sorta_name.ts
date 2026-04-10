/*
1 import * as fs from "fs";
                      ~~~~

main.ts:2:27 - error TS2792: Cannot find module 'readline-sync'. Did you mean to set the 'moduleResolution' option to 'nodenext', or to add aliases to the 'paths' option?

2 import * as readline from "readline-sync";
                            ~~~~~~~~~~~~~~~
*/

import * as fs from "fs";
import * as readline from "readline-sync";

type Roles = Record<string, string[]>;
const roles: Roles = {};

// Utility functions
function trim(s: string): string {
  return s.trim();
}

function capitalize(str: string): string {
  if (!str) return "";
  return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
}

function input(prompt: string): string {
  return trim(readline.question(prompt));
}

// Core features
function addPerson(): void {
  const name = input("Enter name: ");
  
  if (!name) 
  {
    console.log("Name cannot be empty.\n");
    return;
  }

  const role = capitalize(input("Enter role: "));
  if (!role) 
  {
    console.log("Role cannot be empty.\n");
    return;
  }

  if (!roles[role]) 
  {
    roles[role] = [];
  }

  roles[role].push(name);
  console.log(`Added ${name} as ${role}\n`);
}

function showGroups(): void {
  console.log("\n------ People Groups ------");

  const sortedRoles = Object.keys(roles).sort();

  if (sortedRoles.length === 0) 
  {
    console.log("No group found.\n"); //NIXED BUG, nice
    return;
  }

  for (const role of sortedRoles) 
  {
    const people = roles[role].sort();
    console.log(`${role}s: ${people.join(", ")}`);
  }

  console.log();
}

function deletePerson(): void {
  const name = input("Enter name to delete: ");

  for (const role in roles) 
  {
    const people = roles[role];
    const index = people.findIndex(
      p => p.toLowerCase() === name.toLowerCase()
    );

    if (index !== -1) 
    {
      const removed = people.splice(index, 1)[0];
      console.log(`Removed ${removed} from ${role}\n`);
      return;
    }
  }

  console.log("Person not found.\n");
}

function searchPerson(): void {
  const name = input("Search name: ");
  let found = false;

  for (const role in roles) 
  {
      
    for (const person of roles[role]) 
    {
      if (person.toLowerCase().includes(name.toLowerCase())) 
      {
        console.log(`${person} is a ${role}`);
        found = true;
      }
    }
  }

  if (!found) 
  {
    console.log("No person found.\n"); //FIXED BUG
  } 
  else
  {
    console.log();
  }
}

//Save/Load
function saveToFile(): void {
  try {
    let data = "";

    for (const role in roles) 
    {
      for (const person of roles[role]) 
      {
        data += `${role}|${person}\n`;
      }
    }

    fs.writeFileSync("roles.txt", data);
    console.log("Data saved!\n");
  } 
  catch {
    console.log("Error saving file.\n");
  }
}

function loadFromFile(): void {
  if (!fs.existsSync("roles.txt")) return;

  const lines = fs.readFileSync("roles.txt", "utf-8").split("\n");

  for (const line of lines) 
  {
    const [role, name] = line.split("|");
    
    if (role && name)
    {
      if (!roles[role]) roles[role] = [];
      roles[role].push(name);
    }
  }
}

//Menu
function showMenu(): void {
  console.log("===== MENU =====");
  console.log("1. Add Person");
  console.log("2. Show Groups");
  console.log("3. Delete Person");
  console.log("4. Search Person");
  console.log("5. Save");
  console.log("6. Quit");
}

// Main loop
function main(): void {
  loadFromFile();

  while (true) 
  {
    showMenu();
    const choice = input("Choose: ");

    switch (choice) {
      case "1":
        addPerson();
        break;
      case "2":
        showGroups();
        break;
      case "3":
        deletePerson();
        break;
      case "4":
        searchPerson();
        break;
      case "5":
        saveToFile();
        break;
      case "6":
        console.log("\nProgram ended here.");
        return;
      default:
        console.log("Invalid choice.\n");
    }
  }
}

main();
