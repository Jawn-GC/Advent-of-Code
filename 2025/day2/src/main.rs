use std::fs::File;
use std::io::{BufReader, BufRead};

#[derive(Debug)]
struct Range {
    start: u64,
    end: u64,
}

fn main() -> std::io::Result<()> {
    let path = "input.txt";
    let file = File::open(path).expect("Could not open file");
    let reader = BufReader::new(file);

    println!("Parsing ID ranges...");
    let mut ranges: Vec<Range> = Vec::new();
    for result in reader.lines() {
        let line: String = result?;
        let ranges_str: Vec<&str> = line.split(",").collect();
        for r in ranges_str {
            let points: Vec<&str> = r.split("-").collect();
            let start: u64 = points[0].trim().parse::<u64>().expect("Invalid integer");
            let end: u64 = points[1].trim().parse::<u64>().expect("Invalid integer");
            ranges.push(Range{start, end});
        } 
    }

    println!("Finding invalid IDs...");
    let mut invalid_ids: Vec<u64> = Vec::new();
    for range in ranges {
        let mut ids = get_invalid_ids1(range);
        invalid_ids.append(&mut ids); // append moves values from one vector to another, so they both must be mutable
    }

    let mut sum: u64 = 0;
    for id in invalid_ids {
        sum += id;
    }
    println!("[Part 1]: The sum of invalid ids is {sum}");

    Ok(())
}

// Invalid IDs are those that look like a sequence of digits repeated twice.
// So, invalid IDs must have an even number of digits.
// Let pattern be the sequence of digits to be repeated.
// If range.start is odd, start checking for IDs with a pattern of 10^n such that n causes the current_id to have 1 more digit than range.start.
// For example, start checking ID 1010 if range.start = 758.
// If range.start is even, start checking for IDs with a pattern equal the digits in the first half of range.start.
// For example, start checking ID 4545 if range.start = 4513. 
// Continue checking IDs until current_id is out of range.
fn get_invalid_ids1(range: Range) -> Vec<u64> {
    let mut ids = Vec::new();
    let mut num_digits: u64 = get_num_digits(range.start);
    let mut pattern: u64;
    if num_digits % 2 == 1 {
        num_digits += 1;
        pattern = multiply_by_ten(1, num_digits/2-1);
    } else {
        pattern = divide_by_ten(range.start, num_digits/2);
    }

    loop {
        let current_id: u64 = multiply_by_ten(pattern, get_num_digits(pattern)) + pattern;
        if current_id >= range.start && current_id <= range.end {ids.push(current_id);}
        else if current_id > range.end {break;}
        pattern += 1;
    }

    ids
}

fn get_num_digits(mut num: u64) -> u64 {
    let mut counter: u64 = 0;
    while num > 0 {
        num = num / 10;
        counter += 1;
    }
    counter
}

// Only accurate for nonnegative n
fn multiply_by_ten(x: u64, n: u64) -> u64 {
    let mut result: u64 = x; 
    for _ in 1..=n {
        result *= 10; 
    }
    result
}

// Only accurate for nonnegative n
fn divide_by_ten(x: u64, n: u64) -> u64 {
    let mut result: u64 = x; 
    for _ in 1..=n {
        result /= 10; 
    }
    result
}