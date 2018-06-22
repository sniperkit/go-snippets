package tech.sourced.babelfish;

import org.apache.commons.io.IOUtils;
import org.junit.Test;
import org.junit.runner.Request;

import java.io.*;


public class CppParserTest {

    //Only to see the JSON output
    @Test
    public void responsePack() throws IOException {
//        File file = new File("src/test/resources/test.cpp");
//        final BufferedReader reader = new BufferedReader(new FileReader(file));
//        String source = IOUtils.toString(reader);
//        EclipseCPPParser parser = new EclipseCPPParser();
//        final ByteArrayOutputStream baos = new ByteArrayOutputStream();
//        final TranslationUnitJSONMapper responseMapper = new TranslationUnitJSONMapper(true, baos);
//
//        DriverResponse response = new DriverResponse("1.0.0", "CPP", "14");
//        response.setMapper(responseMapper);
//        response.parseCode(parser,source);
//        response.send();
//        System.out.println();
    }

    // FIXME: rewrite this this packing the request outside the DriverRequest object
    @Test
    public void requestPackUnpack() throws IOException {
//        File file = new File("src/test/resources/test.cpp");
//        final BufferedReader reader = new BufferedReader(new FileReader(file));
//        String source = IOUtils.toString(reader);
//        final ByteArrayOutputStream baos = new ByteArrayOutputStream();
//        final RequestMapper requestMapper = new RequestMapper(true, baos);
//
//        DriverRequest request = new DriverRequest("parse", "CPP", "14", source);
//        request.setMapper(requestMapper);
//        request.send();
//
//        DriverRequest request2 = DriverRequest.load(new String(baos.toByteArray()));
//        Boolean equals = request.equals(request2);
//        assert(equals);
    }
}