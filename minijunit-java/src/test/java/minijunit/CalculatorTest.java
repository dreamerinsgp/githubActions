package minijunit;

public class CalculatorTest {
    Calculator calc;

    @Before
    public void setUp() {
        calc = new Calculator();
    }

    @After
    public void tearDown() {
        // optional cleanup
    }

    @Test
    public void addShouldReturnSum() {
        Assert.assertEquals(3, calc.add(1, 2));
    }

    @Test
    public void subtractShouldReturnDifference() {
        Assert.assertEquals(2, calc.subtract(5, 3));
    }

    @Test
    public void divideByZero() {
        Assert.assertThrows(ArithmeticException.class, () -> calc.div(1, 0));
    }

    public static void main(String[] args) {
        TestResult r = MiniJUnitRunner.run(CalculatorTest.class);
        System.out.println(r.summary());
        System.exit(r.failCount() > 0 ? 1 : 0);
    }
}
