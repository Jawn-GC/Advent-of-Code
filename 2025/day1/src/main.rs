use std::fs::File;
use std::io::{BufReader, BufRead};

#[derive(Debug)]
struct Rotation {
    dir: char,
    steps: i32,
}

fn main() -> std::io::Result<()> {
    let path = "input.txt";
    let file = File::open(path).expect("Could not open file");
    let reader = BufReader::new(file);

    println!("Reading instructions...");
    let mut instructions: Vec<Rotation> = Vec::new();
    for result in reader.lines() {
        let line: String = result?;
        let (dir_str, steps_str) = line.split_at(1); // 1 is the byte offset from the start of the String
        let dir: char = dir_str.chars().next().expect("Empty String");
        let steps: i32 = steps_str.trim().parse::<i32>().expect("Invalid integer");
        instructions.push(Rotation{dir,steps});
    }

    println!("Following instructions...");
    let num_ticks: i32 = 100;
    let start_pos: i32 = 50;
    println!("[Part 1]: The password is {}", calc_password1(num_ticks, start_pos, &instructions));
    println!("[Part 2]: The password is {}", calc_password2(num_ticks, start_pos, &instructions));

    Ok(())
}

// The password is the number of times that the dial lands on 0 after a complete rotation step.
fn calc_password1(num_ticks: i32, start_pos: i32, instructions: &Vec<Rotation>) -> i32 {
    let mut current_pos: i32 = start_pos;
    let mut zero_counter: i32 = 0;
    for r in instructions {
        current_pos = rotate(current_pos, r).rem_euclid(num_ticks);
        if current_pos == 0 {zero_counter += 1;}
    }
    zero_counter
}

// The password is the number of times that the dial ticks onto 0.
fn calc_password2(num_ticks: i32, start_pos: i32, instructions: &Vec<Rotation>) -> i32 {
    let mut current_pos: i32 = start_pos;
    let mut zero_counter: i32 = 0;
    for r in instructions {
        zero_counter += r.steps / 100;
        let temp = rotate(current_pos, &Rotation{dir: r.dir, steps: r.steps % 100}).rem_euclid(num_ticks);
        if (temp > current_pos && r.dir == 'L' && current_pos != 0) || (temp < current_pos && r.dir == 'R' && current_pos != 0) || temp == 0 {zero_counter += 1;}
        current_pos = temp;
    }
    zero_counter
}

fn rotate(start_pos: i32, rotation: &Rotation) -> i32 {
    let movement: i32 = {
        if rotation.dir == 'R' {rotation.steps}
        else if rotation.dir == 'L' {-1 * rotation.steps}
        else {0}
    };

    start_pos + movement
}