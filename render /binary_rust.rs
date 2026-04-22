use rand::Rng;
use std::io;
use std::time::Instant;

fn read_number() -> usize {
    loop {
        let mut input = String::new();

        println!("<Lines> :");
        io::stdin().read_line(&mut input).expect("Failed to read");

        match input.trim().parse::<usize>() {
            Ok(n) => return n,
            Err(_) => println!("Invalid input."),
        }
    }
}

fn main() {
    let n = read_number();

    let start = Instant::now();

    let mut rng = rand::thread_rng();
    let mut buffer = String::with_capacity(16);

    for i in 1..=n {
    
        buffer.clear();

        for _ in 0..9 {
        
            buffer.push(if rng.gen_bool(0.5) { '1' } else { '0' });
        }

        println!("[ {}. ] = {}", i, buffer);
    }

    println!("Time elapsed: {:?}", start.elapsed());
}
