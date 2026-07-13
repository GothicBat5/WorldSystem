use std::io::{self, Write};

fn main() 
{
    let height = read_height();
    let symbol = read_symbol();

    draw_pyramid(height, symbol);
}

fn read_height() -> usize {
    loop {
        let mut input = String::new();
        print!("Height: ");
        io::stdout().flush().unwrap();

        io::stdin()
            .read_line(&mut input)
            .expect("Failed to read input.");

        match input.trim().parse::<usize>() 
        {
            Ok(num) if num > 0 => return num,
            _ => println!("Please enter a positive integer."),
        }
    }
}

fn read_symbol() -> char {
    loop {
        let mut input = String::new();
        print!("Symbol (single character): ");
        io::stdout().flush().unwrap();

        io::stdin()
            .read_line(&mut input)
            .expect("Failed to read input.");

        let trimmed = input.trim();
        if trimmed.chars().count() == 1 {
            return trimmed.chars().next().unwrap();
        } 
        else {
            println!("Please enter exactly one character.");
        }
    }
}

fn draw_pyramid(height: usize, symbol: char) 
{
    for row in 1..=height {
        let spaces = height - row;
        let symbols = row * 2 - 1;
        println!("{:width$}{}", "", symbol.to_string().repeat(symbols), width = spaces);
    }
}
