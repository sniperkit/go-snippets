package bblfsh.bash;

import org.junit.Test;

import java.io.ByteArrayInputStream;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;

import com.fasterxml.jackson.databind.ObjectMapper;
import static org.fest.assertions.Assertions.*;

public class RequestReaderTest {

    private static String content1 = "#!/usr/bin/env bash";
    private static String content2 = "#!/bin/bash";

    @Test
    public void twoValid() throws IOException {
        final String input = "{\"content\":\"" + content1 +
            "\"}\n{\"content\":\"" + content2 + "\"}\n";
        final InputStream in = new ByteArrayInputStream(input.getBytes());
        final ObjectMapper mapper = Registry.objectMapper();
        final RequestReader reader = new RequestReader(in, mapper);

        Request request = reader.read();
        Request expected = new Request();
        expected.content = content1;
        assertThat(request.content).isEqualTo(expected.content);
        assertThat(request).isEqualTo(expected);

        request = reader.read();
        expected = new Request();
        expected.content = content2;
        assertThat(request.content).isEqualTo(expected.content);
        assertThat(request).isEqualTo(expected);
    }

    @Test
    public void oneMalformedOnValid() throws IOException {
        final String input = "foo\n{\"content\":\"" + content1 + "\"}\n";
        final InputStream in = new ByteArrayInputStream(input.getBytes());
        final ObjectMapper mapper = Registry.objectMapper();
        final RequestReader reader = new RequestReader(in, mapper);

        boolean thrown = false;
        try {
            reader.read();
        } catch (IOException ex) {
            thrown = true;
        }
        assertThat(thrown).isTrue();

        Request request = reader.read();
        Request expected = new Request();
        expected.content = content1;
        assertThat(request.content).isEqualTo(expected.content);
        assertThat(request).isEqualTo(expected);
    }

    @Test
    public void throwOnClosed() throws IOException {
        final InputStream in = new ByteArrayInputStream(new byte[]{}) {
            @Override
            public int read(byte[] var1) throws IOException {
                throw new IOException("closed");
            }
        };
        final ObjectMapper mapper = Registry.objectMapper();
        final RequestReader reader = new RequestReader(in, mapper);

        boolean thrown = false;
        try {
            reader.read();
        } catch (IOException ex) {
            thrown = true;
        }
        assertThat(thrown).isTrue();
    }
}
