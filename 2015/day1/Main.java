import java.io.*;
import java.util.Scanner;
// import java.util.List;
// import java.util.ArrayList;

public class Main {
    public static void main(String[] args) {

        System.out.println("Reading elevator instructions...");
        File file;
        Scanner scanner;
        try {
            file = new File("input.txt");
            scanner = new Scanner(file);
        } catch (FileNotFoundException e) {
            System.out.println("File not found: " + e.getMessage());
            return;
        }

        // The instructions are expected to be only 1 line.
        System.out.println("Reading elevator instructions...");
        String instructions = "";
        while (scanner.hasNextLine()) {
            instructions = scanner.nextLine();
        }
        scanner.close();

        System.out.println("Following elevator instructions...");
        int starting_floor = 0;
        int ending_floor = followInstructions(starting_floor, instructions);
        System.out.println("[Part 1] Santa arrives at floor " + ending_floor + ".");
        int basement_step = findBasement(starting_floor, instructions);
        System.out.println("[Part 2] Santa first reaches the basement on Step " + basement_step + ".");
    }

    public static int followInstructions(int starting_floor, String instructions) {
        int current_floor = starting_floor;

        for (int i = 0; i < instructions.length(); i++){
            if (instructions.charAt(i) == '(') {
                current_floor++;
            } else if (instructions.charAt(i) == ')') {
                current_floor--;
            }
        }

        return current_floor;
    }

    public static int findBasement(int starting_floor, String instructions) {
        int current_floor = starting_floor;

        for (int i = 0; i < instructions.length(); i++){
            if (instructions.charAt(i) == '(') {
                current_floor++;
            } else if (instructions.charAt(i) == ')') {
                current_floor--;
            }

            if (current_floor == -1) {
                return i+1;
            }
        }

        return -1; // No instruction causes the elevator to reach the basement
    }
}