import java.util.ArrayList;

public class Number {
	private ArrayList<Digit> digits;
	
	public Number() {
		this.digits = new ArrayList<>(); 
	}
	
	public void addDigit(Digit d) {
		this.digits.add(d);
	}
	
	public ArrayList<Digit> getDigits() {
		return this.digits;
	}

	// String representation of a Number object
	@Override
	public String toString() {
		String str = "";

		for (Digit d : this.digits) {
			str += d.getSymbol();
		}

		return str; 
	}

	public int toInt() {
		if (this.digits.size() == 0) {
			return -1;
		} else {
			String numStr = this.toString();
			try {
				int num = Integer.parseInt(numStr);
				return num;
			} catch (NumberFormatException e) {
				return -1;
			}
		}
	}
}