/**
 * 被测类：反射将动态发现并调用其方法
 */
public class Calculator {
    public int add(int a, int b) {
        return a + b;
    }

    public int subtract(int a, int b) {
        return a - b;
    }

    public String getName() {
        return "Calculator";
    }

    private int multiply(int a, int b) {
        return a * b;
    }
}
