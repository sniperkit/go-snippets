package bblfsh.bash;

import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;

public class RequestReader {

    private final BufferedReader reader;
    private final ObjectMapper mapper;

    public RequestReader(final InputStream in, final ObjectMapper mapper) {
        this.reader = new BufferedReader(new InputStreamReader(in));
        this.mapper = mapper;
    }

    public Request read() throws IOException {
        final String line = this.reader.readLine();
        if (line == null) {
            throw new CloseException("end of the stream reached");
        }
        return mapper.readValue(line, Request.class);
    }

}
