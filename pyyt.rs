use std::io;

fn main() 
{
    let height = read_height();

    for row in 1..=height {
        let spaces = " ".repeat(height - row);
        let stars = "*".repeat(row * 2 - 1);

        println!("{}{}", spaces, stars);
    }
}

fn read_height() -> usize {
    loop {
        let mut input = String::new();

        println!("Height:");

        io::stdin()
            .read_line(&mut input)
            .expect("Failed to read input.");

        match input.trim().parse::<usize>() {
            Ok(num) if num > 0 => return num,

            _ => {
                println!("Invalid number.");
            }
        }
    }
}
