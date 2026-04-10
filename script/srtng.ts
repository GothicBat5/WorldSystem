//what the fuck!
import * as readline from "readline";

const MAX_WORDS = 10000;

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

const words: string[] = [];

console.log("Enter some words [type END to finish]");

// Input loop
function ask(): void {
  rl.question("> ", (input) => {
    const word = input.trim();

    if (word.toUpperCase() === "END") 
    {
      finish();
      return;
    }

    if (!word) 
    {
      console.log("Empty input ignored.");
      return ask();
    }

    if (words.length >= MAX_WORDS) 
    {
      console.log("Max word limit reached.");
      return finish();
    }

    words.push(word);
    ask();
  });
}

// Finish + process
function finish(): void {
  rl.close();

  if (words.length === 0) 
  {
    console.log("\nNo words entered.");
    return;
  }

  // Sort (case-insensitive but preserves original)
  const sorted = [...words].sort((a, b) =>
    a.localeCompare(b, undefined, { sensitivity: "base" })
  );

  console.log("\nSorted words:");
  
  for (const word of sorted)
  {
    console.log(word);
  }
}

// Start
ask();
