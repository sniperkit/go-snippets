package tech.sourced.babelfish;

import java.io.*;
import java.util.Arrays;
import java.util.Scanner;

public class Main {

    enum ProcessCycle {CONTINUE, STOP}

    public static void main(String args[]) {
        // TODO: change this for a parsing of args
        boolean doPrint = false;

        if (doPrint) {
            System.err.println("Program args: " + Arrays.toString(args));
            final EclipseCPPParser parser = new EclipseCPPParser();
            try {
                // TODO: remove this testing code
                String testfile = "src/test/resources/test.cpp";
                String code = new Scanner(new File(testfile)).useDelimiter("\\Z").next();
                parser.printAST(code);
            } catch (Exception e) {
                System.err.println("Parsing error: " + e.toString());
                System.exit(1);
            }
        } else {
            BufferedReader in = new BufferedReader(new InputStreamReader(System.in));
            BufferedOutputStream out = new BufferedOutputStream(System.out);
            boolean prettyPrint = true;

            while (true) {
                if (process(in, out, prettyPrint) == ProcessCycle.STOP)
                    return; // stdin closed or unwrittable
            }
        }
    }

    //
    //Try to send and error trough the response driver, print to stdout and stop and print on stderr if that fails
    //
    private static ProcessCycle trySendError(BufferedOutputStream out, String msg, Exception e) {
        // TODO: debug, remove
        e.printStackTrace();
        try {
            final TranslationUnitJSONMapper responseJSONMapper =
                    new TranslationUnitJSONMapper(false, new ByteArrayOutputStream());
            DriverResponse response = new DriverResponse(responseJSONMapper);
            response.sendError(out, e, msg);
            return ProcessCycle.CONTINUE;
        } catch (Exception j) {
            System.err.println(e.getMessage());
            System.err.println("BAILING OUT, CANT WRITE ERRORS");
            System.err.println("ADITTIONAL ERROR WHILE SENDING ERROR BELOW!");
            System.err.println(j.getMessage());
            return ProcessCycle.STOP;
        }
    }

    static private ProcessCycle process(BufferedReader bufferInputStream, BufferedOutputStream bufferOutputStream)
            throws IOException {
        return process(bufferInputStream, bufferOutputStream, false);
    }

    static private ProcessCycle process(BufferedReader bufferInputStream, BufferedOutputStream bufferOutputStream,
                                        boolean prettyPrint) {
        DriverResponse response = null;
            try {
                final EclipseCPPParser parser = new EclipseCPPParser();
                final TranslationUnitJSONMapper responseJSONMapper =
                        new TranslationUnitJSONMapper(prettyPrint, new ByteArrayOutputStream());
                response = new DriverResponse(responseJSONMapper);
                String requestContent;

                // TODO: check that this detects EOF on stdin correctly
                // TODO: check if this can be instantiated outside of the loop

                final String inStr = bufferInputStream.readLine();
                if (inStr == null) {
                    // stdin closed
                    return ProcessCycle.STOP;
                }
                // FIXME: this creates a new ObjectMapper, check if this can be optimized
                requestContent = DriverRequest.load(inStr).content;

                response.parseCode(parser, requestContent);
                response.send(bufferOutputStream);
                return ProcessCycle.CONTINUE;

            } catch (DriverRequest.RequestLoadException e) {
                return trySendError(bufferOutputStream, "Error reading the petition: ", e);
            } catch (DriverResponse.ResponseSendException e) {
                return trySendError(bufferOutputStream, "Error serializing the AST to JSON: ", e);
            } catch (IOException e) {
                return trySendError(bufferOutputStream, "A problem occurred while processing the petition: ", e);
            }
    }
}
