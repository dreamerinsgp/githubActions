package minijunit;

import java.util.ArrayList;
import java.util.List;

public class TestResult {
    private int runCount;
    private int failCount;
    private final List<String> failures = new ArrayList<>();

    public void recordPass() {
        runCount++;
    }

    public void recordFail(String testName, Throwable cause) {
        failCount++;
        String msg = cause != null ? cause.getMessage() : "unknown";
        failures.add(testName + ": " + msg);
    }

    public int runCount() {
        return runCount;
    }

    public int failCount() {
        return failCount;
    }

    public String summary() {
        StringBuilder sb = new StringBuilder();
        sb.append("\n");
        sb.append("--- MiniJUnit Summary ---\n");
        sb.append("Run: ").append(runCount).append(", Pass: ").append(runCount - failCount).append(", Fail: ").append(failCount).append("\n");
        if (!failures.isEmpty()) {
            sb.append("Failures:\n");
            for (String f : failures) {
                sb.append("  - ").append(f).append("\n");
            }
        }
        sb.append("--------------------------\n");
        return sb.toString();
    }
}
