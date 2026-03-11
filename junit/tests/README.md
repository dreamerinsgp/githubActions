# 反射演示

## 什么是反射？

**反射**：程序在运行时，可以查看和修改自己的结构（类、方法、字段、注解等）。

- 不用反射：编译时就知道要调用 `calc.add(1, 2)`
- 用反射：运行时才知道方法名 `"add"`，通过 `Method.invoke()` 动态调用

## 常用 API

| API | 作用 |
|-----|------|
| `Class.getDeclaredMethods()` | 获取类中所有方法 |
| `Method.getName()` | 获取方法名 |
| `Method.getParameterTypes()` | 获取参数类型 |
| `Method.invoke(obj, args)` | 动态调用方法 |
| `Method.setAccessible(true)` | 允许访问 private 方法 |

## 运行

```bash
cd C:\Users\wyq19\Test\junit\tests
javac Calculator.java ReflectionDemo.java
java ReflectionDemo
```
