import java.io.*;
import java.util.Scanner;
import java.util.List;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

public class Main {

    static Map<String, Point> deltas = Map.ofEntries(
        Map.entry("^", new Point(-1, 0)),
        Map.entry(">", new Point(0, 1)),
        Map.entry("v", new Point(1, 0)),
        Map.entry("<", new Point(0, -1))
    );

    record Point(int row, int col) {
        Point add(Point other) {
            return new Point(
                row + other.row,
                col + other.col
            );
        }
    }

    public static void main(String[] args) {
        File file;
        String filename = "input.txt";
        Scanner scanner;
        try {
            file = new File(filename);
            scanner = new Scanner(file);
        } catch (FileNotFoundException e) {
            System.out.println("Could not open " + filename + ": " + e.getMessage());
            return;
        }

        // The input is expected to be one line.
        System.out.println("Reading instructions...");
        String[] instructions = null;
        while (scanner.hasNextLine()) {
            instructions = scanner.nextLine().split("");
        }
        scanner.close();

        System.out.println("Following instructions...");
        Map<Point, Integer> houses_part1 = new HashMap<>();
        followInstructions(Arrays.stream(instructions).toList(), houses_part1);

        List<String> santa_instructions = new ArrayList<>();
        List<String> robot_instructions = new ArrayList<>();
        for (int i = 0; i < instructions.length; i++) {
            if (i % 2 == 0) {
                santa_instructions.add(instructions[i]);
            } else {
                robot_instructions.add(instructions[i]);
            }
        }
        Map<Point, Integer> houses_part2 = new HashMap<>();
        followInstructions(santa_instructions, houses_part2);
        followInstructions(robot_instructions, houses_part2);

        System.out.println("[Part 1] The number of houses that received presents is: " + houses_part1.entrySet().size());
        System.out.println("[Part 2] The number of houses that received presents is: " + houses_part2.entrySet().size());
    }

    public static void followInstructions(List<String> instructions, Map<Point, Integer> houses) {
        Point current_point = new Point(0, 0);
        houses.put(current_point, 1); // Initial state
        for (String dir : instructions) {
            current_point = current_point.add(deltas.get(dir));

            if (houses.containsKey(current_point)) {
                int val = houses.get(current_point);
                houses.put(current_point, val + 1);
            } else {
                houses.put(current_point, 1);
            }
        }
    }
}