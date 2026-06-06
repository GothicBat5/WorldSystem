use std::net::TcpStream;
use std::io::{Read, Write};

fn main() {

    let mut stream = TcpStream::connect("127.0.0.1:7878").expect("Could not connect");
    let message = "Hello server";

    stream.write_all(message.as_bytes()).expect("Failed to send message");

    let mut buffer = [0; 512];
    let size = stream.read(&mut buffer).expect("Failed to read response");
    let response = String::from_utf8_lossy(&buffer[..size]);

    println!("Server says: {}", response);
}
