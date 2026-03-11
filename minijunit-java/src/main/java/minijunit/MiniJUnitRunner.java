package minijunit;

import java.lang.reflect.Constructor;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.List;

public class MiniJUnitRunner {

    public static TestResult run(Class<?> testClass) {
        TestResult result = new TestResult();

        List<Method> testMethods = new ArrayList<>();
        Method beforeMethod = null;
        Method afterMethod = null;

        for (Method m : testClass.getDeclaredMethods()) {
            if (m.isAnnotationPresent(Test.class)) {
                testMethods.add(m);
            }
            if (m.isAnnotationPresent(Before.class)) {
                beforeMethod = m;
            }
            if (m.isAnnotationPresent(After.class)) {
                afterMethod = m;
            }
        }

        for (Method testMethod : testMethods) {
            Object instance;
            try {
                Constructor<?> ctor = testClass.getDeclaredConstructor();
                instance = ctor.newInstance();
            } catch (Exception e) {
                result.recordFail(testMethod.getName(), e);
                continue;
            }

            try {
                if (beforeMethod != null) {
                    beforeMethod.invoke(instance);
                }
                testMethod.invoke(instance);
                result.recordPass();
            } catch (InvocationTargetException e) {
                Throwable cause = e.getCause();
                result.recordFail(testMethod.getName(), cause != null ? cause : e);
            } catch (Exception e) {
                result.recordFail(testMethod.getName(), e);
            } finally {
                if (afterMethod != null) {
                    try {
                        afterMethod.invoke(instance);
                    } catch (Exception e) {
                        result.recordFail(testMethod.getName() + " (in @After)", e);
                    }
                }
            }
        }

        return result;
    }
}
