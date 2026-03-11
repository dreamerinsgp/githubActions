import java.lang.reflect.Method;
import java.lang.reflect.Modifier;

/**
 * 反射演示：展示反射如何「在运行时」发现类的方法、参数、返回值，并动态调用
 *
 * 反射 = 程序在运行时，可以查看和操作自己的结构（类、方法、字段等）
 */
public class ReflectionDemo {
    public static void main(String[] args) throws Exception {
        System.out.println("========== 1. 获取 Class 对象 ==========");
        Class<?> clazz = Calculator.class;
        System.out.println("类名: " + clazz.getName());

        System.out.println("\n========== 2. 发现所有方法 (getDeclaredMethods) ==========");
        Method[] methods = clazz.getDeclaredMethods();
        for (Method m : methods) {
            String mod = Modifier.toString(m.getModifiers());
            String ret = m.getReturnType().getSimpleName();
            String name = m.getName();
            System.out.println("  方法: " + mod + " " + ret + " " + name + "(...)");
        }

        System.out.println("\n========== 3. 动态调用方法 (invoke) ==========");
        Calculator calc = (Calculator) clazz.getDeclaredConstructor().newInstance();

        // 获取 add 方法并调用
        Method addMethod = clazz.getMethod("add", int.class, int.class);
        Object result = addMethod.invoke(calc, 2, 3);
        System.out.println("add(2, 3) = " + result);

        // 获取 subtract 方法并调用
        Method subMethod = clazz.getMethod("subtract", int.class, int.class);
        Object subResult = subMethod.invoke(calc, 10, 4);
        System.out.println("subtract(10, 4) = " + subResult);

        System.out.println("\n========== 4. 反射调用私有方法 (setAccessible) ==========");
        Method multiMethod = clazz.getDeclaredMethod("multiply", int.class, int.class);
        multiMethod.setAccessible(true);  // 绕过 private 访问限制
        Object multiResult = multiMethod.invoke(calc, 5, 6);
        System.out.println("multiply(5, 6) = " + multiResult);
    }
}
