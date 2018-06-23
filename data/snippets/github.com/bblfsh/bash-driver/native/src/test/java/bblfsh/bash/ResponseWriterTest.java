package bblfsh.bash;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.commons.io.IOUtils;
import org.junit.Test;
import static org.fest.assertions.Assertions.assertThat;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;

public class ResponseWriterTest {

    @Test
    public void error() throws IOException {
        final ByteArrayOutputStream out = new ByteArrayOutputStream();
        final ObjectMapper mapper = Registry.objectMapper();
        final ResponseWriter writer = new ResponseWriter(out, mapper);

        Response response;

        response = new Response();
        response.status = "fatal";
        response.errors = new ArrayList<>();
        response.errors.add("error");
        writer.write(response);

        final String result = new String(out.toByteArray());
        final String expected = "{\"status\":\"fatal\",\"errors\":[\"error\"]}\n";
        assertThat(result).isEqualTo(expected);
    }

    @Test
    public void valid() throws IOException {
        final String inputPath = "/helloWorld.bash";
        final String expectedPath = "/helloWorld.expected";

        final String source = IOUtils.toString(
                getClass().getResourceAsStream(inputPath),
                StandardCharsets.UTF_8);
        final String expected = IOUtils.toString(
                getClass().getResourceAsStream(expectedPath),
                StandardCharsets.UTF_8);

        final ByteArrayOutputStream out = new ByteArrayOutputStream();
        final ObjectMapper mapper = Registry.objectMapper();
        final ResponseWriter writer = new ResponseWriter(out, mapper);
        final Parser parser = new Parser();

        Response response = new Response();
        response.status = "ok";
        response.ast = parser.parse(source);

        writer.write(response);

        final String obtained = out.toString(StandardCharsets.UTF_8.name());
        assertThat(obtained).isEqualTo(expected);
    }

}
