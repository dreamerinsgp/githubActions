# MiniJUnit - Java 类 JUnit 测试框架

精简的测试框架，实现：注解发现、断言、Before-After 生命周期、运行器、结果报告。

## 运行

### 使用 javac

```powershell
cd minijunit-java

# 编译
javac -d target\classes src\main\java\minijunit\*.java
javac -d target\test-classes -cp target\classes src\test\java\minijunit\*.java src\test\java\minijunit\api\*.java

# 单元测试
java -cp "target\classes;target\test-classes" minijunit.CalculatorTest

# API 自动化测试（需先启动 Go 服务：go run main.go）
java -cp "target\classes;target\test-classes" minijunit.api.GoApiTest
```

### 使用 Maven

```powershell
mvn compile test-compile exec:java
```

## API 自动化测试（测试 Go 项目接口）

1. 在 Test 项目根目录启动 Go 服务：`go run main.go`
2. 运行：`java -cp "target\classes;target\test-classes" minijunit.api.GoApiTest`

覆盖接口：`/api/health`、`/api/items/slow`、`/api/items`（GET/POST/PUT/DELETE）。

## 项目结构

```
src/main/java/minijunit/
  Test.java, Before.java, After.java  # 注解
  Assert.java                         # 断言
  TestResult.java                     # 结果收集
  MiniJUnitRunner.java                # 运行器

src/test/java/minijunit/
  Calculator.java, CalculatorTest.java  # 单元测试示例

src/test/java/minijunit/api/
  HttpHelper.java, GoApiTest.java       # Go API 自动化测试
```
