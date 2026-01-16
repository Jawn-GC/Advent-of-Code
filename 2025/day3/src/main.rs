use std::fs::File;
use std::io::{BufReader, BufRead};

fn main() -> std::io::Result<()> {
    let path = "input.txt";
    let file = File::open(path).expect("Could not open file");
    let reader = BufReader::new(file);

    println!("Parsing battery banks...");
    let mut banks: Vec<String> = Vec::new();
    for result in reader.lines() {
        let line: String = result?;
        banks.push(line);
    }

    println!("Finding joltages...");
    let mut sum1 = 0;
    let mut sum2 = 0;
    for bank in banks {
        sum1 += get_max_joltage1(&bank);
        sum2 += get_max_joltage2(&bank, 12);
    }

    println!("[Part 1] The sum of max joltages of the battery banks is {sum1}");
    println!("[Part 2] The sum of max joltages of the battery banks is {sum2}");

    Ok(())
}

// The bank is a string of digits (1-9).
// The max joltage is the largest integer that can be created using any two digits (batteries).
// The order of the batteries may not be swapped within the bank.
fn get_max_joltage1(bank: &str) -> u32 {
    let bank_length = bank.len(); 

    let mut first: char = '0';
    let mut second: char = '0';
    for (i, battery1) in bank.chars().enumerate() {
        if i  == bank_length - 1 {break;} // The last battery cannot be the first digit

        if battery1 <= first {continue;}
        else if battery1 > first {
            first = battery1;
            second = '0';
        }
        
        for (j , battery2) in bank.chars().enumerate() {
            if j <= i {continue;}

            if battery2 > second {second = battery2;}
        }
    }

    let a = first.to_digit(10).expect("Not an integer");
    let b = second.to_digit(10).expect("Not an integer");

    a * 10 + b
}


// Same goal as the other function, but creates an n-digit number.
// The bank size is assumed to be large enough.
fn get_max_joltage2(bank: &str, max_batteries: u64) -> u64 {
    let bank_length = bank.len();
    let mut digit_positions: Vec<usize> = Vec::new();
    let mut digits: Vec<char> = Vec::new();

    for n in 0..max_batteries {
        let mut max_digit: char = '0';
        let mut cur_pos: usize = 0;

        let start_pos: usize;
        if n == 0 {start_pos = 0;}
        else {start_pos = digit_positions[n as usize - 1] + 1;}
        for (i, battery) in bank.chars().enumerate() {
            if i < start_pos {continue;}
            if i > bank_length - max_batteries as usize + n as usize {break;}

            if battery > max_digit {
                max_digit = battery;
                cur_pos = i;
            }
        }

        digit_positions.push(cur_pos);
        digits.push(max_digit);
    }

    let mut joltage: u64 = 0;
    for digit in digits {
        joltage = joltage * 10 + (digit as u8 - b'0') as u64;
    }

    joltage
}