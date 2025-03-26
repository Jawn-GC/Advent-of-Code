import java.io.*;
import java.util.Scanner;
import java.util.List;
import java.util.ArrayList;

record Box(int length, int width, int height) {

    // Total surface area of the box
    int getSurfaceArea() {
        return 2*length*width + 2*length*height + 2*width*height;
    }

    // Volume of the box
    int getVolume() {
        return length*width*height;
    }

    // Area of the smallest side
    int getSmallestArea() {
        int smallest_area = length*width;

        if (smallest_area > length*height) {
            smallest_area = length*height;
        }

        if (smallest_area > width*height) {
            smallest_area = width*height;
        }

        return smallest_area;
    }

    // Shortest perimeter around the box
    int getSmallestPerimeter() {
        int p1 = 2*length + 2*width;
        int p2 = 2*length + 2*height;
        int p3 = 2*width + 2*height;
        int smallest_perimeter = p1;

        if (smallest_perimeter > p2) {
            smallest_perimeter = p2;
        }

        if (smallest_perimeter > p3) {
            smallest_perimeter = p3;
        }

        return smallest_perimeter;
    }
}

public class Main {
    public static void main(String[] args) {
        File file;
        Scanner scanner;
        try {
            file = new File("input.txt");
            scanner = new Scanner(file);
        } catch (FileNotFoundException e) {
            System.out.println("File not found: " + e.getMessage());
            return;
        }

        System.out.println("Parsing gift dimensions...");
        List<Box> boxes = new ArrayList<>();
        while(scanner.hasNextLine()) {
            String line = scanner.nextLine();
            String[] temp = line.split("x"); // {l, w, h}
            int l = Integer.parseInt(temp[0]);
            int w = Integer.parseInt(temp[1]);
            int h = Integer.parseInt(temp[2]);
            boxes.add(new Box(l,w,h));
        }
        scanner.close();

        System.out.println("Calculating wrapping paper surface areas...");
        int total_surface_area = 0;
        for (Box box : boxes) {
            total_surface_area += box.getSurfaceArea();
            total_surface_area += box.getSmallestArea();
        }

        System.out.println("Calculating ribbon length...");
        int total_length = 0;
        for (Box box : boxes) {
            total_length += box.getSmallestPerimeter();
            total_length += box.getVolume();
        }

        System.out.println("[Part 1] Total surface area of wrapping paper required: " + total_surface_area);
        System.out.println("[Part 2] Total length of ribbon required: " + total_length);
    }
}