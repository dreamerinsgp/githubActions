package minijunit.api;

import minijunit.*;

/**
 * Go API 自动化测试：对 router 中的接口进行回归测试
 * 运行前需先启动 Go 服务：go run main.go
 */
public class GoApiTest {
    private static final int DEFAULT_PORT = 8080;
    private HttpHelper http;

    @Before
    public void setUp() {
        String baseUrl = System.getenv("API_URL");
        if (baseUrl == null || baseUrl.isEmpty()) {
            baseUrl = "http://localhost:" + DEFAULT_PORT;
        }
        http = new HttpHelper(baseUrl);
    }

    @Test
    public void healthShouldReturnOk() throws Exception {
        HttpHelper.HttpResponseResult r = http.get("/api/health");
        Assert.assertEquals(200, r.statusCode);
        Assert.assertTrue(r.body.contains("ok"), "body should contain 'ok', got: " + r.body);
    }

    @Test
    public void slowShouldReturn200() throws Exception {
        HttpHelper.HttpResponseResult r = http.get("/api/items/slow?ms=10");
        Assert.assertEquals(200, r.statusCode);
    }

    @Test
    public void listItemsShouldReturn200() throws Exception {
        HttpHelper.HttpResponseResult r = http.get("/api/items");
        Assert.assertEquals(200, r.statusCode);
        Assert.assertTrue(r.body.contains("items"), "body should contain 'items'");
    }

    @Test
    public void createItemShouldReturn201() throws Exception {
        HttpHelper.HttpResponseResult r = http.post("/api/items", "{\"name\":\"api-test\",\"description\":\"minijunit\"}");
        Assert.assertEquals(201, r.statusCode);
        Assert.assertTrue(r.body.contains("id"), "body should contain 'id'");
        Assert.assertTrue(r.body.contains("api-test"), "body should contain created name");
    }

    @Test
    public void getItemShouldReturn200Or404() throws Exception {
        HttpHelper.HttpResponseResult createResp = http.post("/api/items", "{\"name\":\"get-test\",\"description\":\"\"}");
        Assert.assertEquals(201, createResp.statusCode);
        int id = extractId(createResp.body);
        HttpHelper.HttpResponseResult r = http.get("/api/items/" + id);
        Assert.assertEquals(200, r.statusCode);
    }

    @Test
    public void updateItemShouldReturn200() throws Exception {
        HttpHelper.HttpResponseResult createResp = http.post("/api/items", "{\"name\":\"orig\",\"description\":\"\"}");
        int id = extractId(createResp.body);
        HttpHelper.HttpResponseResult r = http.put("/api/items/" + id, "{\"name\":\"updated\"}");
        Assert.assertEquals(200, r.statusCode);
    }

    @Test
    public void deleteItemShouldReturn200() throws Exception {
        HttpHelper.HttpResponseResult createResp = http.post("/api/items", "{\"name\":\"to-delete\",\"description\":\"\"}");
        int id = extractId(createResp.body);
        HttpHelper.HttpResponseResult r = http.delete("/api/items/" + id);
        Assert.assertEquals(200, r.statusCode);
    }

    private int extractId(String json) {
        java.util.regex.Matcher m = java.util.regex.Pattern.compile("\"id\"\\s*:\\s*(\\d+)").matcher(json);
        if (m.find()) {
            return Integer.parseInt(m.group(1));
        }
        throw new AssertionError("cannot extract id from: " + json);
    }

    public static void main(String[] args) {
        TestResult r = MiniJUnitRunner.run(GoApiTest.class);
        System.out.println(r.summary());
        System.exit(r.failCount() > 0 ? 1 : 0);
    }
}
