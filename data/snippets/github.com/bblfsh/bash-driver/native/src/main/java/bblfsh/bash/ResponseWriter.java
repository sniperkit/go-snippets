package bblfsh.bash;

import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.OutputStream;

public class ResponseWriter {

    private final OutputStream out;
    private final ObjectMapper mapper;

    public ResponseWriter(final OutputStream out, final ObjectMapper mapper) {
        this.out = out;
        this.mapper = mapper;
    }

    public void write(final Response response) throws IOException {
        // mapper closes the output stream after write, that's why we use an
        // intermediate output stream
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        mapper.writeValue(out, response);
        out.write('\n');
        this.out.write(out.toByteArray());
    }

}
