package minijunit.api;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;

/**
 * HTTP 请求工具，用于 API 自动化测试
 */
public class HttpHelper {
    private final String baseUrl;
    private final HttpClient client;

    public HttpHelper(String baseUrl) {
        this.baseUrl = baseUrl.endsWith("/") ? baseUrl.substring(0, baseUrl.length() - 1) : baseUrl;
        this.client = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(5))
                .build();
    }

    public static HttpHelper localhost(int port) {
        return new HttpHelper("http://localhost:" + port);
    }

    public HttpResponseResult get(String path) throws Exception {
        HttpRequest req = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + path))
                .GET()
                .build();
        HttpResponse<String> resp = client.send(req, HttpResponse.BodyHandlers.ofString());
        return new HttpResponseResult(resp.statusCode(), resp.body());
    }

    public HttpResponseResult post(String path, String jsonBody) throws Exception {
        HttpRequest req = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + path))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(jsonBody))
                .build();
        HttpResponse<String> resp = client.send(req, HttpResponse.BodyHandlers.ofString());
        return new HttpResponseResult(resp.statusCode(), resp.body());
    }

    public HttpResponseResult put(String path, String jsonBody) throws Exception {
        HttpRequest req = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + path))
                .header("Content-Type", "application/json")
                .PUT(HttpRequest.BodyPublishers.ofString(jsonBody))
                .build();
        HttpResponse<String> resp = client.send(req, HttpResponse.BodyHandlers.ofString());
        return new HttpResponseResult(resp.statusCode(), resp.body());
    }

    public HttpResponseResult delete(String path) throws Exception {
        HttpRequest req = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + path))
                .DELETE()
                .build();
        HttpResponse<String> resp = client.send(req, HttpResponse.BodyHandlers.ofString());
        return new HttpResponseResult(resp.statusCode(), resp.body());
    }

    public static class HttpResponseResult {
        public final int statusCode;
        public final String body;

        public HttpResponseResult(int statusCode, String body) {
            this.statusCode = statusCode;
            this.body = body;
        }
    }
}
