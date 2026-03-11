package minijunit;

import java.util.Objects;

public final class Assert {

    private Assert() {
    }

    public static void assertEquals(Object expected, Object actual) {
        if (!Objects.equals(expected, actual)) {
            throw new AssertionError("expected: " + expected + ", but was: " + actual);
        }
    }

    public static void assertTrue(boolean condition) {
        if (!condition) {
            throw new AssertionError("expected: true, but was: false");
        }
    }

    public static void assertTrue(boolean condition, String message) {
        if (!condition) {
            throw new AssertionError(message);
        }
    }

    public static void assertFalse(boolean condition) {
        if (condition) {
            throw new AssertionError("expected: false, but was: true");
        }
    }

    public static void assertNull(Object actual) {
        if (actual != null) {
            throw new AssertionError("expected: null, but was: " + actual);
        }
    }

    public static void assertNotNull(Object actual) {
        if (actual == null) {
            throw new AssertionError("expected: not null, but was: null");
        }
    }

    public static <T extends Throwable> void assertThrows(Class<T> expectedType, Runnable runnable) {
        try {
            runnable.run();
            throw new AssertionError("Expected " + expectedType.getSimpleName() + " but no exception was thrown");
        } catch (Throwable t) {
            if (!expectedType.isInstance(t)) {
                AssertionError e = new AssertionError("Expected " + expectedType.getSimpleName() + " but got " + t.getClass().getSimpleName() + ": " + t.getMessage());
                e.initCause(t);
                throw e;
            }
        }
    }
}
