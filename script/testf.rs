//Error, how ironic
use std::thread;
use std::time::Instant;
use std::sync::{Arc, Mutex};

fn main()
{
    println!("Demo safety...");
    
    let mayValue: Option<i32> = Some(42);
    
    //no null crashed!
    match mayValue 
    {
        Some(v) => println!("SAFE: {}", v),
        None => println!("NO VALUE FOUND"),
    }
    
    let date = String::from("The Ownership");
    
    take_ownership(data);
    
    let mut sum: u64 = 0;
    
    for i in 0..50_000_000 {
        
        sum += 1;
    }
    
    let duration = start.elapsed();
    
    println!("Sum: {}", sum);
    println!("Time: {}", duration);
    
    let counter = Arc::new(Mutex::new(0));
    
    let mut handles = vec![];
    
    for _ in 0..10 {
        
        let counter = thread::spawn(move || {
            
            for _ in 0..1000 {
                let mut num = counter.lock().unwrap();
                *num += 1;
            }
        });
        
        handles.push(handle);
    }
    
    for handle in handles {
        
        handle.join().unwrap();
    }
    
    println!("Final: {}", *counter.lock().unwrap());
}

fn take_ownership(s: String)
{
    println!("Taking ownership of: {}", s);
    
    
}
