use std::net::{TcpListener, TcpStream};
use std::io::{Read, Write};

fn handle_client(mut stream: TcpStream) 
{

    let mut buffer = [0; 512];

    match stream.read(&mut buffer) {
      
        Ok(size) => {
            let received = String::from_utf8_lossy(&buffer[..size]);

            println!("Received: {}", received);

            let response = "Hello from Rust server!";
            stream.write_all(response.as_bytes()).unwrap();
        }
        Err(e) => {
            println!("Failed to read: {}", e);
        }
    }
}

fn main() 
{

    let listener = TcpListener::bind("127.0.0.1:7878").expect("Could not bind port");

    println!("Server running on port 7878...");

    for stream in listener.incoming() {

        match stream {
            Ok(stream) => {
                handle_client(stream);
            }
            Err(e) => {
                println!("Connection failed: {}", e);
            }
        }
    }
}
