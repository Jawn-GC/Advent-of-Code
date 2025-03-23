public class Point {
	private int row;
	private int col;
	
	public Point(int row, int col) {
		this.row = row;
		this.col = col;
	}
	
	public int getRow() {
		return this.row;
    }

	public void setRow(int row) {
		this.row = row;
    }

	public int getColumn() {
		return this.col;
	}

	public void setColumn(int col) {
		this.col = col;
	}

	public Point add(Point other) {
		int new_row = this.getRow() + other.getRow();
		int new_col = this.getColumn() + other.getColumn();
		return new Point(new_row, new_col);
	}
	
	// String representation of a Point object
	@Override
	public String toString() {
		return "Point(Row: " + this.row + ", Column: " + this.col + ")";
	}
}