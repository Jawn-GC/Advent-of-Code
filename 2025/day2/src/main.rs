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
    let mut invalid_ids1: Vec<u64> = Vec::new();
    let mut invalid_ids2: Vec<u64> = Vec::new();
    for range in ranges {
        let mut ids1 = get_invalid_ids1(&range);
        let mut ids2 = get_invalid_ids2(&range);
        invalid_ids1.append(&mut ids1); // append moves values from one vector to another, so they both must be mutable
        invalid_ids2.append(&mut ids2);
    }

    let mut sum1: u64 = 0;
    for id in invalid_ids1 {
        sum1 += id;
    }
    invalid_ids2.sort_unstable();
    invalid_ids2.dedup();
    let mut sum2: u64 = 0;
    for id in invalid_ids2 {
        sum2 += id;
    }
    println!("[Part 1]: The sum of invalid ids is {sum1}");
    println!("[Part 2]: The sum of invalid ids is {sum2}");

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
fn get_invalid_ids1(range: &Range) -> Vec<u64> {
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
        let current_id: u64 = generate_sequence(pattern, 2);
        if current_id >= range.start && current_id <= range.end {ids.push(current_id);}
        else if current_id > range.end {break;}
        pattern += 1;
    }

    ids
}

// Invalid IDs are those that look like a sequence of digits repeated 2 or more times (e.g. 123123123, 55)
// Start with a pattern of 1 and generate sequences with the pattern until the sequence has more digits the end of the range.
// Save the sequences that are within the range.
// Increment the pattern integer by one and repeat the sequence generation step.
// Exit the pattern increment loop once the pattern cannot be repeated even once without exceeding the number of digits of the max value.
fn get_invalid_ids2(range: &Range) -> Vec<u64> {
    let mut ids = Vec::new();
    let max_digits: u64 = get_num_digits(range.end);
    let min_digits: u64 = get_num_digits(range.start);

    let mut pattern: u64 = 1;
    loop {
        let pattern_length: u64 = get_num_digits(pattern);
        if pattern_length * 2 > max_digits {break;}
        for i in 2.. {
            if pattern_length * i > max_digits {
                break;
            }
            else if pattern_length * i < min_digits {
                continue;
            }

            let sequence: u64 = generate_sequence(pattern, i);
            
            if sequence >= range.start && sequence <= range.end {
                ids.push(sequence);
            }
            else if sequence > range.end {
                break;
            }
        }
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

// Does not handle cases that result in integer overflow
fn generate_sequence(pattern: u64, num_copies: u64) -> u64 {
    let mut sequence: u64 = pattern;
    let pattern_length: u64 = get_num_digits(pattern);
    for _ in 1..num_copies {
        sequence = multiply_by_ten(sequence, pattern_length) + pattern;
    }
    sequence
}