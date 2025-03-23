public class Digit {
	private Point position;       // Where the digit is located on a grid
	private String symbol;        // The value of the digit as a string
	private Number parentNumber;  // The Number object that this Digit object belongs to
	
	public Digit(Point p, String s) {
		this.position = p;
		this.symbol = s;
		this.parentNumber = null;
	}
	
	public Point getPosition() {
		return this.position;
	}
	
	public void setPosition(Point p) {
		this.position = p;
	}
	
	public String getSymbol() {
		return this.symbol;
	}
	
	public void setSymbol(String s) {
		this.symbol = s;
	}
	
	public Number getNumber() {
		return this.parentNumber;
	}
	
	public void setParentNumber(Number n) {
		this.parentNumber = n;
	}
	
	// String representation of a Digit object
	@Override
	public String toString() {
		return "Digit(Position: " + this.position.toString() + ", Symbol: " + this.symbol + ")";
	}
}