import java.io.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Arrays;
import java.util.Scanner;

public class Main {	
	
    private static final List<Point> DELTAS = Arrays.asList(
        new Point(-1, 0), // UP 
        new Point(-1, 1), // UP RIGHT
        new Point(0, 1),  // RIGHT
        new Point(1, 1),  // DOWN RIGHT
        new Point(1, 0),  // DOWN
        new Point(1, -1), // DOWN LEFT
        new Point(0, -1), // LEFT
        new Point(-1, -1) // UP LEFT
    );

	public static void main(String[] args) throws Exception {
        
        // Reading input...
        List<List<String>> grid = new ArrayList<>(); // Declare an empty 2D List object
        try {
            File file = new File("input.txt");
            Scanner scanner = new Scanner(file);

            while (scanner.hasNextLine()) {
                String line = scanner.nextLine();
                List<String> grid_row = new ArrayList<>(); // New empty row for grid
				String[] splitLine = line.split(""); // Split every symbol into its own string

                for (String s : splitLine) { // Add symbols to the new grid row (enhanced for-loop)
                    grid_row.add(s);
                }

                grid.add(grid_row);
            }
            scanner.close();

        } catch (FileNotFoundException e) {
            System.out.println("File not found: " + e.getMessage());
            return;
        }

        // Parsing numbers...
        List<Number> partNumbers = getPartNumbers(grid);
        int partNumberSum = 0;
        for (Number num : partNumbers) {
            partNumberSum += num.toInt();
        }

        // Locating gears...
        List<Point> gears = getGears(grid);
        int gearRatioSum = 0;
        for (Point gear : gears) {
            gearRatioSum += getGearRatio(gear, partNumbers);
        }

        System.out.println("[Part 1] Part number sum: " + partNumberSum);
        System.out.println("[Part 2] Gear ratio sum: " + gearRatioSum);
	}  

    // Returns all Numbers found in a grid. A Number is composed of consecutive Digits (left to right)
    public static List<Number> getNumbers(List<List<String>> grid) {
        List<Number> numbers = new ArrayList<>();
        int height = grid.size(); // Number of rows
        int width = grid.get(0).size(); // Number of columns
        
        for (int i = 0; i < height; i++) {
            Number newNumber = new Number(); 
            boolean readingNumber = false; // Will be true if the previous symbol in the row was an integer
            for (int j = 0; j < width; j++) {
                String symbol = grid.get(i).get(j);

                if (symbol.matches("\\d{1}")) { // If symbol is a single-digit integer
                    readingNumber = true;
                    newNumber.addDigit(new Digit(new Point(i, j), symbol));

                    // If we are on a digit at the end of a row, add the newNumber to the numbers list
                    if (j == width - 1) {
                        numbers.add(newNumber);
                    }

                } else if (readingNumber) { // Add newNumber to the numbers list if we passed its last digit
                    readingNumber = false;
                    numbers.add(newNumber);
                    newNumber = new Number();
                }
            }
        }

        // Set the parentNumber field for each digit
        for (Number num : numbers) {
            for (Digit digit : num.getDigits()) {
                digit.setParentNumber(num);
            }
        }

        return numbers;
    }

    // Returns a list of Numbers which have at least one Digit adjacent to a non-numeric symbol (except ".")
    public static List<Number> getPartNumbers(List<List<String>> grid) {
        List<Number> partNumbers = new ArrayList<>();
        List<Number> numbers = getNumbers(grid);

        for (Number num : numbers) {
            digitLoop: for (Digit d : num.getDigits()) {
                for (Point delta : DELTAS) {
                    Point adj_point = d.getPosition().add(delta);
                    if (isOOB(grid, adj_point)) {
                        continue;
                    }

                    String symbol = grid.get(adj_point.getRow()).get(adj_point.getColumn());
                    if (!symbol.matches("[\\d.]")) {
                        partNumbers.add(num);
                        break digitLoop;
                    }
                }
            }
        }

        return partNumbers;
    }

    // Checks if a Point is Out of Bounds of the grid
    public static boolean isOOB(List<List<String>> grid, Point p) {
        int height = grid.size();
        int width = grid.get(0).size();

        if (p.getRow() < 0 || p.getRow() >= height || p.getColumn() < 0 || p.getColumn() >= width) {
            return true;
        }
        
        return false;
    }

    public static List<Point> getGears(List<List<String>> grid) {
        List<Point> gears = new ArrayList<>();
        int height = grid.size();
        int width = grid.get(0).size();

        for (int i = 0; i < height; i++) {
            for (int j = 0; j < width; j++) {
                String symbol = grid.get(i).get(j);
                if (symbol.equals("*")) {
                    gears.add(new Point(i,j));
                }
            }
        }

        return gears;
    }

    public static int getGearRatio(Point gear, List<Number> partNumbers) {
        List<Number> adjNumbers = new ArrayList<>();
        List<Digit> digits = new ArrayList<>();

        for (Number num : partNumbers) {
            for (Digit d : num.getDigits()) {
                digits.add(d);
            }
        }

        for (Point delta : DELTAS) {
            Point adjPoint = gear.add(delta);
            for (Digit d : digits) {
                if (adjPoint.getRow() == d.getPosition().getRow() && adjPoint.getColumn() == d.getPosition().getColumn()) {
                    if (!adjNumbers.contains(d.getNumber())) {
                        adjNumbers.add(d.getNumber());
                    }
                }

                if (adjNumbers.size() == 2) {
                    return adjNumbers.get(0).toInt() * adjNumbers.get(1).toInt();
                } else if (adjNumbers.size() > 2) {
                    return 0;
                }
            }
        }
        
        return 0;
    }
}