use std::io;


fn show_error(message: &str)  
{
    eprintln!("\nError: {}", message);
    std::process::exit(1);
}

fn main() {

    print!("Enter your age: ");
 
    let _ = io::stdout().flush();

    let mut age_input = String::new();
    io::stdin().read_line(&mut age_input).expect("Failed to read line");

    let age_input = age_input.trim();

    if age_input.is_empty() {
        show_error("Age cannot be empty.");
    }


    let age = match age_input.parse::<u32>() {
        Ok(n) => n,
        Err(_) => show_error("Age must be a valid positive number."),
    };


    if age > 120 {
        show_error("That seems like an unrealistic age.");
    }


    print!("Enter your name: ");
    let _ = io::stdout().flush();

    let mut name = String::new();
    io::stdin().read_line(&mut name).expect("Failed to read line");

    let name = name.trim();

    if name.is_empty() {
        show_error("Name cannot be empty.");
    }

    //Determine Status
    let status = if age >= 18 {
        "an adult"
    } else {
        "a minor"
    };


    println!("\nHello, {}! You are {} years old, so you are {}.", name, age, status);
}
