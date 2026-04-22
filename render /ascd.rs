use std::f32;

const WIDTH: usize = 80;
const HEIGHT: usize = 40;

#[derive(Clone, Copy)]
struct Vec3 {
    x: f32,
    y: f32,
    z: f32,
}

fn project(v: Vec3) -> Vec3 {
    
    Vec3 {
        x: v.x / v.z,
        y: v.y / v.z,
        z: v.z,
    }
}

fn to_screen(v: Vec3) -> (i32, i32) {
    let x = ((v.x + 1.0) * 0.5 * WIDTH as f32) as i32;
    let y = ((1.0 - (v.y + 1.0) * 0.5) * HEIGHT as f32) as i32;
    (x, y)
}

fn draw_point(buffer: &mut Vec<Vec<char>>, x: i32, y: i32) 
  {
    if x >= 0 && x < WIDTH as i32 && y >= 0 && y < HEIGHT as i32 {
        buffer[y as usize][x as usize] = '#';
    }
}

fn draw_line(buffer: &mut Vec<Vec<char>>, mut x0: i32, mut y0: i32, x1: i32, y1: i32)
  {
    let dx = (x1 - x0).abs();
    let dy = -(y1 - y0).abs();
    let sx = if x0 < x1 { 1 } else { -1 };
    let sy = if y0 < y1 { 1 } else { -1 };
    let mut err = dx + dy;

    loop {
        draw_point(buffer, x0, y0);

        if x0 == x1 && y0 == y1 {
            break;
        }

        let e2 = 2 * err;
        if e2 >= dy {
            err += dy;
            x0 += sx;
        }
        if e2 <= dx {
            err += dx;
            y0 += sy;
        }
    }
}

fn main() {
    let mut buffer = vec![vec!['.'; WIDTH]; HEIGHT];

    let tri = [
        Vec3 { x: -0.5, y: -0.5, z: 1.0 },
        Vec3 { x:  0.5, y: -0.5, z: 1.0 },
        Vec3 { x:  0.0, y:  0.5, z: 1.0 },
    ];

    let mut points = vec![];
    for v in tri {
        let p = project(v);
        points.push(to_screen(p));
    }

    // Draw triangle edges
    draw_line(&mut buffer, points[0].0, points[0].1, points[1].0, points[1].1);
    draw_line(&mut buffer, points[1].0, points[1].1, points[2].0, points[2].1);
    draw_line(&mut buffer, points[2].0, points[2].1, points[0].0, points[0].1);

   
    for row in buffer {
        let line: String = row.into_iter().collect();
        println!("{}", line);
    }
}
